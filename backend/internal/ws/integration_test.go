package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func setupTestServer(t *testing.T) (*Hub, *httptest.Server) {
	t.Helper()
	hub := NewHub(nil)
	server := httptest.NewServer(http.HandlerFunc(hub.HandleUpgrade))
	t.Cleanup(server.Close)
	return hub, server
}

func dial(t *testing.T, server *httptest.Server) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func wsSend(t *testing.T, conn *websocket.Conn, msg ClientMessage) {
	t.Helper()
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func wsRecv(t *testing.T, conn *websocket.Conn) ServerMessage {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var msg ServerMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return msg
}

func TestWSCreateAndJoinTable(t *testing.T) {
	_, server := setupTestServer(t)

	// Client 1 creates a table
	conn1 := dial(t, server)
	wsSend(t, conn1, ClientMessage{Type: MsgCreateTable})
	resp := wsRecv(t, conn1)

	if resp.Type != MsgTableCreated {
		t.Fatalf("type = %q, want %q", resp.Type, MsgTableCreated)
	}

	payload := resp.Payload.(map[string]any)
	tableID := payload["tableID"].(string)
	if tableID == "" {
		t.Fatal("tableID should not be empty")
	}

	// Client 2 joins the table
	conn2 := dial(t, server)
	wsSend(t, conn2, ClientMessage{Type: MsgJoinTable, TableID: tableID})
	resp = wsRecv(t, conn2)

	if resp.Type != MsgTableJoined {
		t.Fatalf("type = %q, want %q", resp.Type, MsgTableJoined)
	}
}

func TestWSJoinNonexistent(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgJoinTable, TableID: "bogus"})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want %q", resp.Type, MsgError)
	}
}

func TestWSSitAndStart(t *testing.T) {
	_, server := setupTestServer(t)

	// Connect 4 clients, all create/join and sit
	conns := make([]*websocket.Conn, 4)
	dirs := []string{"N", "E", "S", "W"}

	// First client creates the table
	conns[0] = dial(t, server)
	wsSend(t, conns[0], ClientMessage{Type: MsgCreateTable})
	resp := wsRecv(t, conns[0])
	tableID := resp.Payload.(map[string]any)["tableID"].(string)

	// First client sits
	wsSend(t, conns[0], ClientMessage{Type: MsgSit, Direction: dirs[0]})
	resp = wsRecv(t, conns[0])
	if resp.Type != MsgSeated {
		t.Fatalf("client 0: type = %q, want %q", resp.Type, MsgSeated)
	}

	// Other 3 clients join and sit
	for i := 1; i < 4; i++ {
		conns[i] = dial(t, server)
		wsSend(t, conns[i], ClientMessage{Type: MsgJoinTable, TableID: tableID})
		wsRecv(t, conns[i]) // table_joined

		wsSend(t, conns[i], ClientMessage{Type: MsgSit, Direction: dirs[i]})
		resp = wsRecv(t, conns[i])
		if resp.Type != MsgSeated {
			t.Fatalf("client %d: type = %q, want %q", i, resp.Type, MsgSeated)
		}
	}

	// Start the game
	wsSend(t, conns[0], ClientMessage{Type: MsgStart, Seed: seedPtr(42)})

	// All 4 clients should receive game_state
	for i, conn := range conns {
		resp = wsRecv(t, conn)
		if resp.Type != MsgGameState {
			t.Errorf("client %d: type = %q, want %q", i, resp.Type, MsgGameState)
		}
	}
}

func TestWSBidBeforeSitting(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgBid, Call: "1NT"})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error", resp.Type)
	}
}

func TestWSUnknownMessageType(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: "nonsense"})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error", resp.Type)
	}
}

func TestWSSitWithoutJoining(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgSit, Direction: "N"})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error", resp.Type)
	}
}

func TestWSSitInvalidDirection(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgCreateTable})
	wsRecv(t, conn)

	wsSend(t, conn, ClientMessage{Type: MsgSit, Direction: "X"})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error", resp.Type)
	}
}

func TestWSSitEmptyDirection(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgCreateTable})
	wsRecv(t, conn)

	wsSend(t, conn, ClientMessage{Type: MsgSit, Direction: ""})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error", resp.Type)
	}
}

func TestWSSitOccupiedSeat(t *testing.T) {
	_, server := setupTestServer(t)

	conn1 := dial(t, server)
	wsSend(t, conn1, ClientMessage{Type: MsgCreateTable})
	resp := wsRecv(t, conn1)
	tableID := resp.Payload.(map[string]any)["tableID"].(string)

	wsSend(t, conn1, ClientMessage{Type: MsgSit, Direction: "N"})
	wsRecv(t, conn1) // seated

	conn2 := dial(t, server)
	wsSend(t, conn2, ClientMessage{Type: MsgJoinTable, TableID: tableID})
	wsRecv(t, conn2) // joined

	wsSend(t, conn2, ClientMessage{Type: MsgSit, Direction: "N"})
	resp = wsRecv(t, conn2)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error for occupied seat", resp.Type)
	}
}

func TestWSStartNotEnoughPlayers(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgCreateTable})
	wsRecv(t, conn)

	wsSend(t, conn, ClientMessage{Type: MsgSit, Direction: "N"})
	wsRecv(t, conn)

	wsSend(t, conn, ClientMessage{Type: MsgStart, Seed: seedPtr(42)})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error for not enough players", resp.Type)
	}
}

func TestWSInvalidBidString(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgCreateTable})
	wsRecv(t, conn)
	wsSend(t, conn, ClientMessage{Type: MsgSit, Direction: "N"})
	wsRecv(t, conn)

	wsSend(t, conn, ClientMessage{Type: MsgBid, Call: "ZZZZZ"})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error for invalid bid", resp.Type)
	}
}

func TestWSInvalidCardString(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgCreateTable})
	wsRecv(t, conn)
	wsSend(t, conn, ClientMessage{Type: MsgSit, Direction: "N"})
	wsRecv(t, conn)

	wsSend(t, conn, ClientMessage{Type: MsgPlayCard, Card: "ZZ"})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error for invalid card", resp.Type)
	}
}

func TestWSMalformedJSON(t *testing.T) {
	_, server := setupTestServer(t)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// Send garbage bytes
	conn.WriteMessage(websocket.TextMessage, []byte("{not json!!!"))

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var resp ServerMessage
	json.Unmarshal(data, &resp)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error for malformed JSON", resp.Type)
	}
}

func TestWSPlayCardBeforeSitting(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgPlayCard, Card: "SA"})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error", resp.Type)
	}
}

func TestWSStartBeforeSitting(t *testing.T) {
	_, server := setupTestServer(t)
	conn := dial(t, server)

	wsSend(t, conn, ClientMessage{Type: MsgStart, Seed: seedPtr(42)})
	resp := wsRecv(t, conn)

	if resp.Type != MsgError {
		t.Fatalf("type = %q, want error", resp.Type)
	}
}
