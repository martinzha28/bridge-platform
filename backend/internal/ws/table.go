package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/martinzha28/bridge-platform/backend/internal/bot"
	"github.com/martinzha28/bridge-platform/backend/internal/game"
	"github.com/martinzha28/bridge-platform/backend/internal/repository"
)

type Table struct {
	ID        string
	name      string
	description string
	mu        sync.Mutex
	game      *game.Game
	players   map[game.Direction]Player
	observers map[*Client]struct{} // clients watching the lobby / play
	reapTimer *time.Timer
	gameRepo  *repository.GameRepository
	persisted bool
	boardNum  int
	closed    bool
	chatLog   []ChatMessageView
	chatSeq   int64
}

func NewTable(id string, gameRepo *repository.GameRepository) *Table {
	return &Table{
		ID:        id,
		name:      "Practice table",
		players:   make(map[game.Direction]Player),
		observers: make(map[*Client]struct{}),
		gameRepo:  gameRepo,
		boardNum:  1,
	}
}

// AddObserver registers a client as watching this table and cancels any
// pending reap. Safe to call more than once for the same client.
func (t *Table) AddObserver(c *Client) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.reapTimer != nil {
		t.reapTimer.Stop()
		t.reapTimer = nil
	}
	t.observers[c] = struct{}{}
}

func (t *Table) RemoveObserver(c *Client) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.observers, c)
}

// HasObservers reports whether any client is still watching the table.
func (t *Table) HasObservers() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.observers) > 0
}

func (t *Table) Sit(p Player, dir game.Direction) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("table is closed")
	}
	if _, taken := t.players[dir]; taken {
		return fmt.Errorf("seat %v is taken", dir)
	}

	for d, player := range t.players {
		if player == p {
			return fmt.Errorf("already seated at %v", d)
		}
	}

	t.players[dir] = p
	t.broadcastTableState()
	return nil
}

func (t *Table) SitBot(dir game.Direction, difficulty bot.Difficulty) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("table is closed")
	}
	if _, taken := t.players[dir]; taken {
		return fmt.Errorf("seat %v is taken", dir)
	}

	bc := newBotClient(t, dir, bot.New(difficulty))
	t.players[dir] = bc
	go bc.run()
	t.broadcastTableState()
	return nil
}

func (t *Table) Start(seed *int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("table is closed")
	}
	if t.game != nil && t.game.Phase != game.PhaseComplete {
		return fmt.Errorf("game already in progress")
	}
	if len(t.players) != game.NumDirections {
		return fmt.Errorf("need 4 players to start, have %d", len(t.players))
	}

	s := time.Now().UnixNano()
	if seed != nil {
		s = *seed
	}

	t.game = game.NewGame(t.boardNum)
	t.persisted = false
	if err := t.game.Deal(s); err != nil {
		return err
	}

	t.broadcastTableState()
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
	t.onCompleteOnce()
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
	t.onCompleteOnce()
	return nil
}

// SetName updates the table's display name and notifies watchers.
// A blank name resets to the default; long names are truncated.
func (t *Table) SetName(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		name = "Practice table"
	}
	if len(name) > 60 {
		name = strings.TrimSpace(name[:60])
	}
	t.name = name
	t.broadcastTableState()
}

// SetDescription updates the table's blurb and notifies watchers.
func (t *Table) SetDescription(desc string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	desc = strings.TrimSpace(desc)
	if len(desc) > 280 {
		desc = strings.TrimSpace(desc[:280])
	}
	t.description = desc
	t.broadcastTableState()
}

// RemoveBot vacates a seat held by a bot and stops its goroutine.
// No-op if the seat is empty; errors if a human holds it or a game is live.
func (t *Table) RemoveBot(dir game.Direction) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("table is closed")
	}
	if t.game != nil && t.game.Phase != game.PhaseComplete {
		return fmt.Errorf("game already in progress")
	}
	p, ok := t.players[dir]
	if !ok {
		return nil
	}
	bc, ok := p.(*botClient)
	if !ok {
		return fmt.Errorf("seat %v is not a bot", dir)
	}
	delete(t.players, dir)
	close(bc.send)
	t.broadcastTableState()
	return nil
}

func (t *Table) RemovePlayer(p Player) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for dir, player := range t.players {
		if player == p {
			delete(t.players, dir)
			break
		}
	}
	t.broadcastTableState()
}

// Shutdown closes the table's bot goroutines and clears its seats.
// Idempotent; call after removing the table from the hub.
func (t *Table) Shutdown() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return
	}
	t.closed = true
	if t.reapTimer != nil {
		t.reapTimer.Stop()
		t.reapTimer = nil
	}
	for _, p := range t.players {
		if bc, ok := p.(*botClient); ok {
			close(bc.send)
		}
	}
	t.players = map[game.Direction]Player{}
	t.observers = map[*Client]struct{}{}
}

// broadcastState sends each seated player their personalized game view.
// Must be called with t.mu held.
func (t *Table) broadcastState() {
	for dir, player := range t.players {
		view := t.game.ViewFor(dir)
		msg := ServerMessage{
			Type:    MsgGameState,
			Payload: view,
		}
		data, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		player.Send(data)
	}
}

// onCompleteOnce advances the board number and persists the game exactly once.
// Must be called with t.mu held.
func (t *Table) onCompleteOnce() {
	if t.persisted || t.game.Phase != game.PhaseComplete {
		return
	}

	t.persisted = true
	t.boardNum++

	if t.gameRepo == nil {
		return
	}
	record := repository.GameFromSession(t.game)
	if err := t.gameRepo.SaveGame(context.Background(), &record); err != nil {
		log.Printf("failed to persist game: %v", err)
	}
}
