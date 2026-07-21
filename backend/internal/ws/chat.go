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

// ChatMessageView is one line of table chat. Ephemeral — the log lives
// only in memory and dies with the table.
type ChatMessageView struct {
	ID     int64  `json:"id"`     // per-table sequence, for ordering / React keys
	Sender string `json:"sender"` // seat name ("West") or "Observer"
	Seat   string `json:"seat"`   // "" for observers, else "North".."West"
	Text   string `json:"text"`
	At     int64  `json:"at"` // unix millis
}

// appendChat records a message in the table's ring buffer and returns it.
// Must be called with t.mu held.
func (t *Table) appendChat(sender, seat, text string) ChatMessageView {
	t.chatSeq++
	msg := ChatMessageView{
		ID:     t.chatSeq,
		Sender: sender,
		Seat:   seat,
		Text:   text,
		At:     time.Now().UnixMilli(),
	}
	t.chatLog = append(t.chatLog, msg)
	if len(t.chatLog) > chatHistoryLimit {
		t.chatLog = t.chatLog[len(t.chatLog)-chatHistoryLimit:]
	}
	return msg
}

// broadcastChat sends one chat line to every watching client.
// Must be called with t.mu held.
func (t *Table) broadcastChat(msg ChatMessageView) {
	data, err := json.Marshal(ServerMessage{Type: MsgChatMessage, Payload: msg})
	if err != nil {
		return
	}
	for c := range t.observers {
		c.Send(data)
	}
}

// Chat validates, records, and broadcasts a line from a client. Blank
// messages are dropped; long ones are truncated.
func (t *Table) Chat(sender, seat, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	if len(text) > chatMaxLen {
		text = strings.TrimSpace(text[:chatMaxLen])
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return
	}
	t.broadcastChat(t.appendChat(sender, seat, text))
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
