package ws

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	// chatHistoryLimit is how many recent lines a joiner is shown.
	chatHistoryLimit = 50
	// chatMaxLen caps a single message; longer text is truncated.
	chatMaxLen = 500
)

const (
	chatKindUser   = "user"
	chatKindSystem = "system" // "<name> joined" / "<name> left"
)

// ChatMessageView is one line of table chat. Ephemeral — the log lives
// only in memory and dies with the table.
type ChatMessageView struct {
	ID     int64  `json:"id"`     // per-table sequence, for ordering / React keys
	Kind   string `json:"kind"`   // chatKindUser | chatKindSystem
	Sender string `json:"sender"` // display name; "" for system lines
	Guest  bool   `json:"guest"`
	Seat   string `json:"seat"` // "" for observers, else "North".."West"
	Text   string `json:"text"`
	At     int64  `json:"at"` // unix millis
}

// appendChat stamps a message, records it in the table's ring buffer and
// returns it. Must be called with t.mu held.
func (t *Table) appendChat(msg ChatMessageView) ChatMessageView {
	t.chatSeq++
	msg.ID = t.chatSeq
	msg.At = time.Now().UnixMilli()
	t.chatLog = append(t.chatLog, msg)
	if len(t.chatLog) > chatHistoryLimit {
		t.chatLog = t.chatLog[len(t.chatLog)-chatHistoryLimit:]
	}
	return msg
}

// appendSystem records a system line. Must be called with t.mu held.
func (t *Table) appendSystem(text string) ChatMessageView {
	return t.appendChat(ChatMessageView{Kind: chatKindSystem, Text: text})
}

// broadcastChat sends one chat line to every watching client.
// Must be called with t.mu held.
func (t *Table) broadcastChat(msg ChatMessageView) {
	t.broadcastChatExcept(msg, nil)
}

// broadcastChatExcept is broadcastChat skipping one client.
// Must be called with t.mu held.
func (t *Table) broadcastChatExcept(msg ChatMessageView, skip *Client) {
	data, err := json.Marshal(ServerMessage{Type: MsgChatMessage, Payload: msg})
	if err != nil {
		return
	}
	for c := range t.observers {
		if c != skip {
			c.Send(data)
		}
	}
}

// Chat validates, records, and broadcasts a line from a client at this
// table. Blank messages are dropped; long ones are truncated.
func (t *Table) Chat(c *Client, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	if len(text) > chatMaxLen {
		text = strings.TrimSpace(text[:chatMaxLen])
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	m, ok := t.observers[c]
	if t.closed || !ok {
		return
	}
	t.broadcastChat(t.appendChat(ChatMessageView{
		Kind:   chatKindUser,
		Sender: m.name,
		Guest:  c.guest,
		Seat:   t.seatOf(c),
		Text:   text,
	}))
}

// SendChatHistory replays the recent chat log to a single client.
func (t *Table) SendChatHistory(c *Client) {
	t.mu.Lock()
	log := append([]ChatMessageView(nil), t.chatLog...)
	t.mu.Unlock()
	if len(log) == 0 {
		return
	}
	data, err := json.Marshal(ServerMessage{Type: MsgChatHistory, Payload: log})
	if err != nil {
		return
	}
	c.Send(data)
}
