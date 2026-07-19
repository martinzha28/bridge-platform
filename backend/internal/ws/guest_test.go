package ws

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWSGuestConnection verifies that a connection with no authenticated
// user (the non-production case) is still upgraded and can drive a table.
func TestWSGuestConnection(t *testing.T) {
	hub := NewHub(nil)
	server := httptest.NewServer(http.HandlerFunc(hub.HandleUpgrade))
	t.Cleanup(server.Close)

	conn := dial(t, server)
	wsSend(t, conn, ClientMessage{Type: MsgCreateTable})
	resp := wsRecv(t, conn)

	if resp.Type != MsgTableCreated {
		t.Fatalf("type = %q, want %q", resp.Type, MsgTableCreated)
	}
}
