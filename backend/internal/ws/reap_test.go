package ws

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func waitUntil(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition not met within 2s")
}

func TestTableReapedWhenLastHumanDisconnects(t *testing.T) {
	hub := NewHub(nil)
	hub.reapGrace = 20 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(hub.HandleUpgrade))
	t.Cleanup(server.Close)

	conn := dial(t, server)
	wsSend(t, conn, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, conn).Payload.(map[string]any)["tableID"].(string)

	// Seat the human and a bot so shutdown of bot goroutines is exercised.
	wsSend(t, conn, ClientMessage{Type: MsgSit, Direction: "S"})
	wsRecv(t, conn)
	wsSend(t, conn, ClientMessage{Type: MsgSitBot, Direction: "N", Difficulty: 1})
	wsRecv(t, conn)

	if _, ok := hub.GetTable(id); !ok {
		t.Fatal("table should exist while the client is connected")
	}

	conn.Close()

	waitUntil(t, func() bool {
		_, ok := hub.GetTable(id)
		return !ok
	})
}

func TestTableKeptWhileAnotherHumanConnected(t *testing.T) {
	hub := NewHub(nil)
	hub.reapGrace = 20 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(hub.HandleUpgrade))
	t.Cleanup(server.Close)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, host).Payload.(map[string]any)["tableID"].(string)
	wsSend(t, host, ClientMessage{Type: MsgSit, Direction: "N"})
	wsRecv(t, host)

	guest := dial(t, server)
	wsSend(t, guest, ClientMessage{Type: MsgJoinTable, TableID: id})
	wsRecv(t, guest)
	wsSend(t, guest, ClientMessage{Type: MsgSit, Direction: "S"})
	wsRecv(t, guest)

	host.Close()
	time.Sleep(200 * time.Millisecond) // let the server process the disconnect

	if _, ok := hub.GetTable(id); !ok {
		t.Fatal("table should survive while the guest is still connected")
	}
}
