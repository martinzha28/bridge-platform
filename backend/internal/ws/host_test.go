package ws

import (
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// waitHost reads table_state broadcasts until want(host) holds, where
// host is the name of the person flagged as host ("" if none).
func waitHost(t *testing.T, conn *websocket.Conn, want func(host string) bool) string {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		p := wsRecvType(t, conn, MsgTableState).Payload.(map[string]any)
		host := ""
		for _, raw := range p["people"].([]any) {
			if person := raw.(map[string]any); person["host"] == true {
				host = person["name"].(string)
			}
		}
		if want(host) {
			return host
		}
	}
	t.Fatal("host condition not met within 2s")
	return ""
}

// waitTableName reads table_state broadcasts until the table has the
// wanted name, failing on an error reply (e.g. a host-only refusal).
func waitTableName(t *testing.T, conn *websocket.Conn, want string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		msg := wsRecvRaw(t, conn)
		if msg.Type == MsgError {
			t.Fatalf("rename to %q refused: %s", want, msg.Error)
		}
		if msg.Type == MsgTableState && msg.Payload.(map[string]any)["name"] == want {
			return
		}
	}
	t.Fatalf("table never renamed to %q", want)
}

func sitAs(t *testing.T, conn *websocket.Conn, dir string) {
	t.Helper()
	wsSend(t, conn, ClientMessage{Type: MsgSit, Direction: dir})
	wsRecvType(t, conn, MsgSeated)
}

func wantHostError(t *testing.T, conn *websocket.Conn, msg ClientMessage) {
	t.Helper()
	wsSend(t, conn, msg)
	got := wsRecvType(t, conn, MsgError)
	if !strings.Contains(got.Error, "host") {
		t.Errorf("%s: error = %q, want a host-only refusal", msg.Type, got.Error)
	}
}

func TestCreatorIsHost(t *testing.T) {
	_, server := setupTestServer(t)
	host := dial(t, server)
	_, you := createAs(t, host, "Martin", "g1")
	if you["host"] != true {
		t.Errorf("creator you = %v, want host", you)
	}
	waitHost(t, host, func(h string) bool { return h == "Martin" })
}

func TestNonHostCannotManageTable(t *testing.T) {
	_, server := setupTestServer(t)
	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")
	sitAs(t, host, "S")
	wsSend(t, host, ClientMessage{Type: MsgSitBot, Direction: "N", Difficulty: 1})
	wsRecvType(t, host, MsgSeated)

	guest := dial(t, server)
	if you := joinAs(t, guest, id, "Bob", "g2"); you["host"] != false {
		t.Errorf("joiner you = %v, want not host", you)
	}

	wantHostError(t, guest, ClientMessage{Type: MsgRemoveBot, Direction: "N"})
	wantHostError(t, guest, ClientMessage{Type: MsgSitBot, Direction: "E", Difficulty: 1})
	wantHostError(t, guest, ClientMessage{Type: MsgSetName, Name: "mine now"})
	wantHostError(t, guest, ClientMessage{Type: MsgSetDescription, Description: "x"})
	sitAs(t, guest, "W") // sitting in an open seat is still allowed
	wantHostError(t, guest, ClientMessage{Type: MsgStart})
}

func TestHostPassesToLongestSeatedHuman(t *testing.T) {
	hub, server := setupTestServer(t)
	hub.leaveGrace = 20 * time.Millisecond

	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")
	sitAs(t, host, "S")

	ann := dial(t, server)
	joinAs(t, ann, id, "Ann", "g2")
	sitAs(t, ann, "N")
	time.Sleep(10 * time.Millisecond)

	bob := dial(t, server)
	joinAs(t, bob, id, "Bob", "g3")
	sitAs(t, bob, "E")

	joinAs(t, dial(t, server), id, "Watcher", "g4") // spectators never inherit

	host.Close()
	for {
		if line := nextSystemLine(t, bob); line == "Ann is now the host" {
			break
		}
	}
	waitHost(t, bob, func(h string) bool { return h == "Ann" })
	wsSend(t, ann, ClientMessage{Type: MsgSetName, Name: "Ann's table"})
	waitTableName(t, ann, "Ann's table")
}

func TestHostReconnectKeepsHost(t *testing.T) {
	hub, server := setupTestServer(t)
	hub.leaveGrace = 20 * time.Millisecond

	draft := dial(t, server)
	id, _ := createAs(t, draft, "Host", "g1")
	ann := dial(t, server)
	joinAs(t, ann, id, "Ann", "g2")
	sitAs(t, ann, "N")

	// The /create -> lobby handoff: a second socket for the same guest
	// joins, then the first closes.
	lobby := dial(t, server)
	if you := joinAs(t, lobby, id, "Host", "g1"); you["host"] != true {
		t.Errorf("host's second socket you = %v, want host", you)
	}
	draft.Close()
	time.Sleep(100 * time.Millisecond) // well past the grace

	wsSend(t, lobby, ClientMessage{Type: MsgSetName, Name: "still mine"})
	waitTableName(t, lobby, "still mine")
}

func TestNextSitterBecomesHostWhenNoneSeated(t *testing.T) {
	hub, server := setupTestServer(t)
	hub.leaveGrace = 20 * time.Millisecond

	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")
	watcher := dial(t, server)
	joinAs(t, watcher, id, "Watcher", "g2")

	host.Close()
	waitHost(t, watcher, func(h string) bool { return h == "" })

	sitAs(t, watcher, "S")
	waitHost(t, watcher, func(h string) bool { return h == "Watcher" })
}
