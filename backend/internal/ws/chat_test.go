package ws

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestChatBroadcastAndHistory(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")
	wsSend(t, host, ClientMessage{Type: MsgSit, Direction: "S"})
	if got := wsRecv(t, host); got.Type != MsgSeated {
		t.Fatalf("sit ack = %q, want seated", got.Type)
	}

	guest := dial(t, server)
	wsSend(t, guest, ClientMessage{Type: MsgJoinTable, TableID: id})
	wsRecvType(t, guest, MsgTableJoined)

	wsSend(t, host, ClientMessage{Type: MsgChat, Text: "  hello table  "})

	p := nextUserLine(t, guest)
	if p["sender"] != "Host" || p["seat"] != "South" || p["kind"] != "user" {
		t.Errorf("sender/seat/kind = %v/%v/%v, want Host/South/user", p["sender"], p["seat"], p["kind"])
	}
	if p["text"] != "hello table" {
		t.Errorf("text = %q, want trimmed %q", p["text"], "hello table")
	}

	// A client joining later is handed the recent history.
	late := dial(t, server)
	wsSend(t, late, ClientMessage{Type: MsgJoinTable, TableID: id})
	rows, ok := wsRecvType(t, late, MsgChatHistory).Payload.([]any)
	if !ok {
		t.Fatalf("history payload = %#v", rows)
	}
	var user []map[string]any
	for _, r := range rows {
		if m := r.(map[string]any); m["kind"] == "user" {
			user = append(user, m)
		}
	}
	if len(user) != 1 || user[0]["text"] != "hello table" {
		t.Errorf("user lines in history = %v, want just %q", user, "hello table")
	}
}

func TestChatFromObserver(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	id, _ := createAs(t, host, "Host", "g1")

	watcher := dial(t, server)
	joinAs(t, watcher, id, "Rail", "g2")

	wsSend(t, watcher, ClientMessage{Type: MsgChat, Text: "hi from the rail"})
	p := nextUserLine(t, host)
	if p["sender"] != "Rail" || p["seat"] != "" {
		t.Errorf("observer sender/seat = %v/%q, want Rail/\"\"", p["sender"], p["seat"])
	}
}

func TestChatDropsBlank(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	wsRecv(t, host)

	wsSend(t, host, ClientMessage{Type: MsgChat, Text: "   "})
	wsSend(t, host, ClientMessage{Type: MsgChat, Text: "real"})

	if got := wsRecvType(t, host, MsgChatMessage).Payload.(map[string]any)["text"]; got != "real" {
		t.Errorf("first chat_message text = %v, want blank dropped then %q", got, "real")
	}
}

// nextUserLine reads until a non-system chat line arrives.
func nextUserLine(t *testing.T, conn *websocket.Conn) map[string]any {
	t.Helper()
	for {
		p := wsRecvType(t, conn, MsgChatMessage).Payload.(map[string]any)
		if p["kind"] != "system" {
			return p
		}
	}
}
