package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/martinzha28/bridge-platform/backend/internal/game"
	"github.com/martinzha28/bridge-platform/backend/internal/repository"
)

type Table struct {
	ID        string
	mu        sync.Mutex
	game      *game.Game
	players   map[game.Direction]*Client
	gameRepo  *repository.GameRepository
	persisted bool
}

func NewTable(id string, gameRepo *repository.GameRepository) *Table {
	return &Table{
		ID:       id,
		players:  make(map[game.Direction]*Client),
		gameRepo: gameRepo,
	}
}

func (t *Table) Sit(c *Client, dir game.Direction) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, taken := t.players[dir]; taken {
		return fmt.Errorf("seat %v is taken", dir)
	}

	for d, player := range t.players {
		if player == c {
			return fmt.Errorf("already seated at %v", d)
		}
	}

	t.players[dir] = c
	return nil
}

func (t *Table) Start(seed int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.game != nil && t.game.Phase != game.PhaseComplete {
		return fmt.Errorf("game already in progress")
	}
	if len(t.players) != game.NumDirections {
		return fmt.Errorf("need 4 players to start, have %d", len(t.players))
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	t.game = game.NewGame(1)
	t.persisted = false
	if err := t.game.Deal(seed); err != nil {
		return err
	}

	t.broadcastState()
	return nil
}

func (t *Table) Bid(dir game.Direction, call game.Call) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.game == nil {
		return fmt.Errorf("game has not started")
	}
	if err := t.game.MakeCall(dir, call); err != nil {
		return err
	}

	t.broadcastState()
	t.persistIfComplete()
	return nil
}

func (t *Table) PlayCard(dir game.Direction, card game.Card) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.game == nil {
		return fmt.Errorf("game has not started")
	}
	if err := t.game.PlayCard(dir, card); err != nil {
		return err
	}

	t.broadcastState()
	t.persistIfComplete()
	return nil
}

func (t *Table) RemoveClient(c *Client) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for dir, player := range t.players {
		if player == c {
			delete(t.players, dir)
			break
		}
	}
}

// broadcastState sends each seated player their personalized game view.
// Must be called with t.mu held.
func (t *Table) broadcastState() {
	for dir, client := range t.players {
		view := t.game.ViewFor(dir)
		msg := ServerMessage{
			Type:    MsgGameState,
			Payload: view,
		}
		data, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		select {
		case client.send <- data:
		default:
		}
	}
}

// persistIfComplete saves the finished game to the database exactly once.
// Must be called with t.mu held.
func (t *Table) persistIfComplete() {
	if t.gameRepo == nil || t.persisted || t.game.Phase != game.PhaseComplete {
		return
	}

	t.persisted = true
	record := repository.GameFromSession(t.game)
	if err := t.gameRepo.SaveGame(context.Background(), &record); err != nil {
		t.persisted = false
		log.Printf("failed to persist game: %v", err)
	}
}
