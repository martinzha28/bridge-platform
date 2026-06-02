package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/martinzha28/bridge-platform/backend/internal/bot"
	"github.com/martinzha28/bridge-platform/backend/internal/game"
	"github.com/martinzha28/bridge-platform/backend/internal/repository"
)

type Table struct {
	ID        string
	mu        sync.Mutex
	game      *game.Game
	players   map[game.Direction]Player
	gameRepo  *repository.GameRepository
	persisted bool
	boardNum  int
}

func NewTable(id string, gameRepo *repository.GameRepository) *Table {
	return &Table{
		ID:       id,
		players:  make(map[game.Direction]Player),
		gameRepo: gameRepo,
		boardNum: 1,
	}
}

func (t *Table) Sit(p Player, dir game.Direction) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, taken := t.players[dir]; taken {
		return fmt.Errorf("seat %v is taken", dir)
	}

	for d, player := range t.players {
		if player == p {
			return fmt.Errorf("already seated at %v", d)
		}
	}

	t.players[dir] = p
	return nil
}

func (t *Table) SitBot(dir game.Direction, difficulty bot.Difficulty) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, taken := t.players[dir]; taken {
		return fmt.Errorf("seat %v is taken", dir)
	}

	bc := newBotClient(t, dir, bot.New(difficulty))
	t.players[dir] = bc
	go bc.run()
	return nil
}

func (t *Table) Start(seed *int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

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

func (t *Table) RemovePlayer(p Player) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for dir, player := range t.players {
		if player == p {
			delete(t.players, dir)
			break
		}
	}
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
