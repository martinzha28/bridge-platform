package ws

import (
	"testing"
)

func seatMap(t *testing.T, msg ServerMessage) map[string]any {
	t.Helper()
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		t.Fatalf("table_state payload = %T", msg.Payload)
	}
	seats, ok := payload["seats"].(map[string]any)
	if !ok {
		t.Fatalf("seats = %T", payload["seats"])
	}
	return seats
}

func TestTableStateSeenByJoiner(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, host).Payload.(map[string]any)["tableID"].(string)

	wsSend(t, host, ClientMessage{Type: MsgSit, Direction: "S"})
	wsRecv(t, host)
	wsSend(t, host, ClientMessage{Type: MsgSitBot, Direction: "N", Difficulty: 1})
	wsRecv(t, host)

	// A guest joining sees the current seat map right away.
	guest := dial(t, server)
	wsSend(t, guest, ClientMessage{Type: MsgJoinTable, TableID: id})
	seats := seatMap(t, wsRecvType(t, guest, MsgTableState))
	if seats["South"] != "human" || seats["North"] != "bot" {
		t.Errorf("seats = %v, want S=human N=bot", seats)
	}
	if seats["East"] != "" || seats["West"] != "" {
		t.Errorf("E/W should be open, got %v / %v", seats["East"], seats["West"])
	}

	// Guest takes West; the next table_state reflects it.
	wsSend(t, guest, ClientMessage{Type: MsgSit, Direction: "W"})
	seats = seatMap(t, wsRecvType(t, guest, MsgTableState))
	if seats["West"] != "human" {
		t.Errorf("after sit W: West = %v, want human", seats["West"])
	}
}

func TestTableStateStartedFlag(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, host).Payload.(map[string]any)["tableID"].(string)
	for _, d := range []string{"N", "E", "S"} {
		wsSend(t, host, ClientMessage{Type: MsgSitBot, Direction: d, Difficulty: 1})
		wsRecv(t, host)
	}
	wsSend(t, host, ClientMessage{Type: MsgSit, Direction: "W"})
	wsRecv(t, host)

	// Fresh watcher: its first table_state is the pre-start snapshot.
	watcher := dial(t, server)
	wsSend(t, watcher, ClientMessage{Type: MsgJoinTable, TableID: id})
	if wsRecvType(t, watcher, MsgTableState).Payload.(map[string]any)["started"] != false {
		t.Fatal("started should be false before start")
	}

	wsSend(t, host, ClientMessage{Type: MsgStart, Seed: seedPtr(42)})
	if wsRecvType(t, watcher, MsgTableState).Payload.(map[string]any)["started"] != true {
		t.Error("started should be true after start")
	}
}

func TestStandVacatesSeat(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, host).Payload.(map[string]any)["tableID"].(string)

	wsSend(t, host, ClientMessage{Type: MsgSit, Direction: "S"})
	if got := wsRecv(t, host); got.Type != MsgSeated {
		t.Fatalf("sit ack = %q, want seated", got.Type)
	}

	wsSend(t, host, ClientMessage{Type: MsgStand})
	if seats := seatMap(t, wsRecvType(t, host, MsgTableState)); seats["South"] != "" {
		t.Errorf("after stand: South = %v, want empty", seats["South"])
	}
	if got := wsRecvType(t, host, MsgStood); got.Type != MsgStood {
		t.Fatalf("stand ack = %q, want stood", got.Type)
	}

	// The seat is free again: a joiner can take it.
	guest := dial(t, server)
	wsSend(t, guest, ClientMessage{Type: MsgJoinTable, TableID: id})
	wsRecvType(t, guest, MsgTableState)
	wsSend(t, guest, ClientMessage{Type: MsgSit, Direction: "S"})
	if seats := seatMap(t, wsRecvType(t, guest, MsgTableState)); seats["South"] != "human" {
		t.Errorf("guest sit S: South = %v, want human", seats["South"])
	}
}

func TestSetNameBroadcast(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, host).Payload.(map[string]any)["tableID"].(string)

	if name := wsRecvType(t, host, MsgTableState).Payload.(map[string]any)["name"]; name != "Practice table" {
		t.Errorf("default name = %v, want %q", name, "Practice table")
	}

	watcher := dial(t, server)
	wsSend(t, watcher, ClientMessage{Type: MsgJoinTable, TableID: id})
	wsRecvType(t, watcher, MsgTableState)

	wsSend(t, host, ClientMessage{Type: MsgSetName, Name: "  Kitchen table  "})
	if name := wsRecvType(t, watcher, MsgTableState).Payload.(map[string]any)["name"]; name != "Kitchen table" {
		t.Errorf("renamed to %v, want %q (trimmed)", name, "Kitchen table")
	}

	// A blank name resets to the default.
	wsSend(t, host, ClientMessage{Type: MsgSetName, Name: "   "})
	if name := wsRecvType(t, watcher, MsgTableState).Payload.(map[string]any)["name"]; name != "Practice table" {
		t.Errorf("blank rename = %v, want default", name)
	}
}

func TestSetDescriptionBroadcast(t *testing.T) {
	_, server := setupTestServer(t)

	host := dial(t, server)
	wsSend(t, host, ClientMessage{Type: MsgCreateTable})
	id := wsRecv(t, host).Payload.(map[string]any)["tableID"].(string)

	if d := wsRecvType(t, host, MsgTableState).Payload.(map[string]any)["description"]; d != "" {
		t.Errorf("default description = %q, want empty", d)
	}

	watcher := dial(t, server)
	wsSend(t, watcher, ClientMessage{Type: MsgJoinTable, TableID: id})
	wsRecvType(t, watcher, MsgTableState)

	wsSend(t, host, ClientMessage{Type: MsgSetDescription, Description: "  friendly practice, all welcome  "})
	if d := wsRecvType(t, watcher, MsgTableState).Payload.(map[string]any)["description"]; d != "friendly practice, all welcome" {
		t.Errorf("description = %q, want trimmed", d)
	}
}
