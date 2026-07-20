package ws

import (
	"encoding/json"

	"github.com/martinzha28/bridge-platform/backend/internal/game"
)

// TableStateView is the lobby view of a table: who holds each seat and
// whether the game has started. Broadcast to everyone watching the table
// on any seat change, so the waiting room stays in sync before there is
// any game_state to send.
type TableStateView struct {
	TableID string            `json:"tableID"`
	Seats   map[string]string `json:"seats"` // direction -> "human" | "bot" | "" (open)
	Started bool              `json:"started"`
}

var seatOrder = [game.NumDirections]game.Direction{
	game.North, game.East, game.South, game.West,
}

// tableStateView builds the current lobby view. Must be called with t.mu held.
func (t *Table) tableStateView() TableStateView {
	seats := make(map[string]string, game.NumDirections)
	for _, d := range seatOrder {
		switch t.players[d].(type) {
		case *Client:
			seats[d.String()] = "human"
		case *botClient:
			seats[d.String()] = "bot"
		default:
			seats[d.String()] = ""
		}
	}
	return TableStateView{
		TableID: t.ID,
		Seats:   seats,
		Started: t.game != nil && t.game.Phase != game.PhaseDeal,
	}
}

// broadcastTableState sends the lobby view to every watching client.
// Must be called with t.mu held.
func (t *Table) broadcastTableState() {
	data, err := json.Marshal(ServerMessage{Type: MsgTableState, Payload: t.tableStateView()})
	if err != nil {
		return
	}
	for c := range t.observers {
		c.Send(data)
	}
}

// BroadcastTableState is broadcastTableState for callers that don't hold t.mu.
func (t *Table) BroadcastTableState() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.broadcastTableState()
}
