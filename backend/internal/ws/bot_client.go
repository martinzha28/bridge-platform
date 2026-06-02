package ws

import (
	"encoding/json"
	"log"

	"github.com/martinzha28/bridge-platform/backend/internal/bot"
	"github.com/martinzha28/bridge-platform/backend/internal/game"
)

type botClient struct {
	table *Table
	dir   game.Direction
	bot   bot.Bot
	send  chan []byte
}

func newBotClient(t *Table, dir game.Direction, b bot.Bot) *botClient {
	return &botClient{
		table: t,
		dir:   dir,
		bot:   b,
		send:  make(chan []byte, sendBufferSize),
	}
}

func (b *botClient) Send(data []byte) {
	select {
	case b.send <- data:
	default:
	}
}

// run reads game_state messages and acts when it is the bot's turn.
// Runs in its own goroutine, started by Table.SitBot.
func (b *botClient) run() {
	for data := range b.send {
		var msg struct {
			Type    string          `json:"type"`
			Payload game.PlayerView `json:"payload"`
		}
		if err := json.Unmarshal(data, &msg); err != nil || msg.Type != MsgGameState {
			continue
		}

		view := msg.Payload
		if view.Turn != view.Seat {
			continue
		}

		switch view.Phase {
		case "Auction":
			if len(view.LegalCalls) == 0 {
				continue
			}
			callStr := b.bot.ChooseCall(view.LegalCalls)
			call, ok := game.ParseCall(callStr)
			if !ok {
				log.Printf("bot: invalid call %q", callStr)
				continue
			}
			if err := b.table.Bid(b.dir, call); err != nil {
				log.Printf("bot: bid error: %v", err)
			}

		case "Play":
			if len(view.LegalCards) == 0 {
				continue
			}
			cardStr := b.bot.ChooseCard(view.LegalCards)
			card, ok := game.ParseCard(cardStr)
			if !ok {
				log.Printf("bot: invalid card %q", cardStr)
				continue
			}
			if err := b.table.PlayCard(b.dir, card); err != nil {
				log.Printf("bot: play error: %v", err)
			}
		}
	}
}
