package ws

import "testing"

func TestChatBroadcastAndHistory(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, host).Payload.(map[string]any)["tableID"].(string)
	wsSend(t, host, ClientMessage{Type: MsgSit, Direction: "S"})
	if got := wsRecv(t, host); got.Type != MsgSeated {
		t.Fatalf("sit ack = %q, want seated", got.Type)
	}

	guest := dial(t, server)
	wsSend(t, guest, ClientMessage{Type: MsgJoinTable, TableID: id})
	wsRecvType(t, guest, MsgTableJoined)

	wsSend(t, host, ClientMessage{Type: MsgChat, Text: "  hello table  "})

	p := wsRecvType(t, guest, MsgChatMessage).Payload.(map[string]any)
	if p["sender"] != "South" || p["seat"] != "South" {
		t.Errorf("sender/seat = %v/%v, want South/South", p["sender"], p["seat"])
	}
	if p["text"] != "hello table" {
		t.Errorf("text = %q, want trimmed %q", p["text"], "hello table")
	}

	// A client joining later is handed the recent history.
	late := dial(t, server)
	wsSend(t, late, ClientMessage{Type: MsgJoinTable, TableID: id})
	rows, ok := wsRecvType(t, late, MsgChatHistory).Payload.([]any)
	if !ok || len(rows) != 1 {
		t.Fatalf("history payload = %#v", rows)
	}
	if rows[0].(map[string]any)["text"] != "hello table" {
		t.Errorf("history[0].text = %v", rows[0].(map[string]any)["text"])
	}
}

func TestChatFromObserver(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, host).Payload.(map[string]any)["tableID"].(string)

	watcher := dial(t, server)
	wsSend(t, watcher, ClientMessage{Type: MsgJoinTable, TableID: id})
	wsRecvType(t, watcher, MsgTableJoined)

	wsSend(t, watcher, ClientMessage{Type: MsgChat, Text: "hi from the rail"})
	p := wsRecvType(t, host, MsgChatMessage).Payload.(map[string]any)
	if p["sender"] != "Observer" || p["seat"] != "" {
		t.Errorf("observer sender/seat = %v/%q, want Observer/\"\"", p["sender"], p["seat"])
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
