package ws

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/martinzha28/bridge-platform/backend/internal/auth"
	"github.com/martinzha28/bridge-platform/backend/internal/models"
)

var guestNameRe = regexp.MustCompile(`^Guest \d{4}$`)

// createAs creates a table as a guest with the given nickname and
// returns the table ID and the server's "you" payload.
func createAs(t *testing.T, conn *websocket.Conn, name, guestID string) (string, map[string]any) {
	t.Helper()
	wsSend(t, conn, ClientMessage{Type: MsgCreateTable, DisplayName: name, GuestID: guestID})
	p := wsRecvType(t, conn, MsgTableCreated).Payload.(map[string]any)
	return p["tableID"].(string), youOf(t, p)
}

func youOf(t *testing.T, payload map[string]any) map[string]any {
	t.Helper()
	you, ok := payload["you"].(map[string]any)
	if !ok {
		t.Fatalf("payload has no \"you\": %v", payload)
	}
	return you
}

func joinAs(t *testing.T, conn *websocket.Conn, tableID, name, guestID string) map[string]any {
	t.Helper()
	wsSend(t, conn, ClientMessage{Type: MsgJoinTable, TableID: tableID, DisplayName: name, GuestID: guestID})
	return youOf(t, wsRecvType(t, conn, MsgTableJoined).Payload.(map[string]any))
}

// nextSystemLine reads until a system chat line arrives and returns its text.
func nextSystemLine(t *testing.T, conn *websocket.Conn) string {
	t.Helper()
	for {
		p := wsRecvType(t, conn, MsgChatMessage).Payload.(map[string]any)
		if p["kind"] == "system" {
			return p["text"].(string)
		}
	}
}

func TestSanitizeDisplayName(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"  Martin  ", "Martin", true},
		{"a  b", "a b", true},
		{"José Ñ", "José Ñ", true},
		{"under_score-dash", "under_score-dash", true},
		{strings.Repeat("x", 20), strings.Repeat("x", 20), true},
		{strings.Repeat("x", 21), "", false},
		{"a", "", false},
		{"   ", "", false},
		{"bad<name>", "", false},
		{"emoji 🃏", "", false},
	}
	for _, c := range cases {
		got, ok := sanitizeDisplayName(c.in)
		if ok != c.ok || got != c.want {
			t.Errorf("sanitizeDisplayName(%q) = %q, %v; want %q, %v", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestGuestNameFromCreate(t *testing.T) {
	_, server := setupTestServer(t)
	_, you := createAs(t, dial(t, server), "Martin", "g1")
	if you["name"] != "Martin" || you["guest"] != true {
		t.Errorf("you = %v, want guest Martin", you)
	}
	if _, ok := you["id"].(float64); !ok {
		t.Errorf("you.id = %#v, want a number", you["id"])
	}
}

func TestInvalidGuestNameFallsBack(t *testing.T) {
	_, server := setupTestServer(t)
	_, you := createAs(t, dial(t, server), "x", "g1")
	if !guestNameRe.MatchString(you["name"].(string)) {
		t.Errorf("name = %q, want Guest NNNN", you["name"])
	}
}

func TestDuplicateNameGetsSuffix(t *testing.T) {
	_, server := setupTestServer(t)
	id, _ := createAs(t, dial(t, server), "Martin", "g1")
	you := joinAs(t, dial(t, server), id, "martin", "g2")
	if you["name"] != "martin (2)" {
		t.Errorf("name = %q, want %q", you["name"], "martin (2)")
	}
}

func TestSameGuestSecondSocketKeepsName(t *testing.T) {
	_, server := setupTestServer(t)
	id, _ := createAs(t, dial(t, server), "Martin", "g1")
	you := joinAs(t, dial(t, server), id, "Martin", "g1")
	if you["name"] != "Martin" {
		t.Errorf("name = %q, want unsuffixed Martin for the same guest", you["name"])
	}
}

func TestPeopleInTableState(t *testing.T) {
	_, server := setupTestServer(t)
	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")
	wsSend(t, host, ClientMessage{Type: MsgSit, Direction: "S"})
	wsRecvType(t, host, MsgSeated)

	joinAs(t, dial(t, server), id, "Watcher", "g2")

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		p := wsRecvType(t, host, MsgTableState).Payload.(map[string]any)
		people, _ := p["people"].([]any)
		if len(people) != 2 {
			continue
		}
		a, b := people[0].(map[string]any), people[1].(map[string]any)
		if a["name"] != "Host" || a["seat"] != "South" || a["guest"] != true {
			t.Errorf("people[0] = %v, want seated guest Host at South", a)
		}
		if b["name"] != "Watcher" || b["seat"] != "" {
			t.Errorf("people[1] = %v, want unseated Watcher", b)
		}
		return
	}
	t.Fatal("never saw a table_state with two people")
}

func TestJoinPostsSystemLine(t *testing.T) {
	_, server := setupTestServer(t)
	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")

	joinAs(t, dial(t, server), id, "Bob", "g2")
	if got := nextSystemLine(t, host); got != "Bob joined" {
		t.Errorf("system line = %q, want %q", got, "Bob joined")
	}
}

func TestLeavePostsSystemLineAfterGrace(t *testing.T) {
	hub, server := setupTestServer(t)
	hub.leaveGrace = 20 * time.Millisecond
	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")

	bob := dial(t, server)
	joinAs(t, bob, id, "Bob", "g2")
	nextSystemLine(t, host) // Bob joined

	bob.Close()
	if got := nextSystemLine(t, host); got != "Bob left" {
		t.Errorf("system line = %q, want %q", got, "Bob left")
	}
}

func TestQuickRejoinSuppressesLeaveAndJoin(t *testing.T) {
	hub, server := setupTestServer(t)
	hub.leaveGrace = 300 * time.Millisecond
	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")

	bob := dial(t, server)
	joinAs(t, bob, id, "Bob", "g2")
	nextSystemLine(t, host) // Bob joined

	bob.Close()
	time.Sleep(50 * time.Millisecond)
	joinAs(t, dial(t, server), id, "Bob", "g2")
	time.Sleep(400 * time.Millisecond) // past the grace window

	wsSend(t, host, ClientMessage{Type: MsgChat, Text: "marker"})
	for {
		p := wsRecvType(t, host, MsgChatMessage).Payload.(map[string]any)
		if p["text"] == "marker" {
			return
		}
		if p["kind"] == "system" {
			t.Errorf("unexpected system line %q on a quick rejoin", p["text"])
		}
	}
}

type fakeUsers struct {
	names map[uuid.UUID]string
}

func (f fakeUsers) GetByID(_ context.Context, id uuid.UUID) (*models.User, error) {
	name, ok := f.names[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &models.User{ID: id, Username: name}, nil
}

// serverAsUser serves the hub with every connection authenticated as id.
func serverAsUser(t *testing.T, hub *Hub, id uuid.UUID) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.HandleUpgrade(w, r.WithContext(auth.InjectUserID(r.Context(), id)))
	}))
	t.Cleanup(server.Close)
	return server
}

func TestLoggedInUserGetsUsername(t *testing.T) {
	id := uuid.New()
	hub := NewHub(nil)
	hub.SetUserLookup(fakeUsers{names: map[uuid.UUID]string{id: "alice"}})
	server := serverAsUser(t, hub, id)

	_, you := createAs(t, dial(t, server), "Imposter", "g1")
	if you["name"] != "alice" || you["guest"] != false {
		t.Errorf("you = %v, want registered alice (client nickname ignored)", you)
	}
}

func TestUserLookupFailureFallsBackToGuest(t *testing.T) {
	hub := NewHub(nil)
	hub.SetUserLookup(fakeUsers{names: map[uuid.UUID]string{}})
	server := serverAsUser(t, hub, uuid.New())

	_, you := createAs(t, dial(t, server), "Martin", "g1")
	if you["name"] != "Martin" || you["guest"] != true {
		t.Errorf("you = %v, want guest Martin", you)
	}
}

func TestSameGuestRenamedSocketTakesNewName(t *testing.T) {
	_, server := setupTestServer(t)
	id, _ := createAs(t, dial(t, server), "Martin", "g1")
	you := joinAs(t, dial(t, server), id, "Marty", "g1")
	if you["name"] != "Marty" {
		t.Errorf("name = %q, want the newly chosen %q", you["name"], "Marty")
	}
}
