package ws

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/martinzha28/bridge-platform/backend/internal/game"
)

func newTestClient() *Client {
	return &Client{
		send: make(chan []byte, sendBufferSize),
	}
}

func readServerMsg(t *testing.T, c *Client) ServerMessage {
	t.Helper()
	select {
	case data := <-c.send:
		var msg ServerMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("unmarshal server message: %v", err)
		}
		return msg
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
		return ServerMessage{}
	}
}

func drainMessages(c *Client) {
	for {
		select {
		case <-c.send:
		default:
			return
		}
	}
}

func seatAllPlayers(t *testing.T, table *Table) [game.NumDirections]*Client {
	t.Helper()
	dirs := [game.NumDirections]game.Direction{game.North, game.East, game.South, game.West}
	var clients [game.NumDirections]*Client

	for i, dir := range dirs {
		c := newTestClient()
		if err := table.Sit(c, dir); err != nil {
			t.Fatalf("Sit(%v): %v", dir, err)
		}
		clients[i] = c
	}
	return clients
}

func TestTableSit(t *testing.T) {
	table := NewTable("test", nil)
	c := newTestClient()

	if err := table.Sit(c, game.North); err != nil {
		t.Fatalf("Sit(North): %v", err)
	}
}

func TestTableSitOccupied(t *testing.T) {
	table := NewTable("test", nil)
	c1 := newTestClient()
	c2 := newTestClient()

	if err := table.Sit(c1, game.North); err != nil {
		t.Fatalf("Sit(North): %v", err)
	}
	if err := table.Sit(c2, game.North); err == nil {
		t.Error("expected error sitting in occupied seat")
	}
}

func TestTableSitDifferentSeats(t *testing.T) {
	table := NewTable("test", nil)
	c1 := newTestClient()
	c2 := newTestClient()

	if err := table.Sit(c1, game.North); err != nil {
		t.Fatalf("Sit(North): %v", err)
	}
	if err := table.Sit(c2, game.East); err != nil {
		t.Fatalf("Sit(East): %v", err)
	}
}

func TestTableStartNotEnoughPlayers(t *testing.T) {
	table := NewTable("test", nil)
	c := newTestClient()
	table.Sit(c, game.North)

	if err := table.Start(42); err == nil {
		t.Error("expected error starting with 1 player")
	}
}

func TestTableStartSuccess(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)

	if err := table.Start(42); err != nil {
		t.Fatalf("Start: %v", err)
	}

	// All 4 players should receive a game_state message
	for i, c := range clients {
		msg := readServerMsg(t, c)
		if msg.Type != MsgGameState {
			t.Errorf("client %d: got type %q, want %q", i, msg.Type, MsgGameState)
		}
	}
}

func TestTableStartTwice(t *testing.T) {
	table := NewTable("test", nil)
	seatAllPlayers(t, table)

	if err := table.Start(42); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := table.Start(99); err == nil {
		t.Error("expected error starting a second game while one is in progress")
	}
}

func TestTableBidBeforeStart(t *testing.T) {
	table := NewTable("test", nil)

	if err := table.Bid(game.North, game.BidCall(1, game.NoTrump)); err == nil {
		t.Error("expected error bidding before game starts")
	}
}

func TestTableBidBroadcasts(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)
	table.Start(42)

	// Drain the initial game_state messages
	for _, c := range clients {
		drainMessages(c)
	}

	// Board 1 dealer is North
	if err := table.Bid(game.North, game.BidCall(1, game.NoTrump)); err != nil {
		t.Fatalf("Bid: %v", err)
	}

	for i, c := range clients {
		msg := readServerMsg(t, c)
		if msg.Type != MsgGameState {
			t.Errorf("client %d: got type %q, want %q", i, msg.Type, MsgGameState)
		}
	}
}

func TestTableRemoveClient(t *testing.T) {
	table := NewTable("test", nil)
	c := newTestClient()
	table.Sit(c, game.North)

	table.RemoveClient(c)

	// Seat should be free again
	c2 := newTestClient()
	if err := table.Sit(c2, game.North); err != nil {
		t.Errorf("North should be free after RemoveClient: %v", err)
	}
}

func TestTablePlayCardBeforeStart(t *testing.T) {
	table := NewTable("test", nil)

	card := game.NewCard(game.Spades, game.Ace)
	if err := table.PlayCard(game.North, card); err == nil {
		t.Error("expected error playing card before game starts")
	}
}

func TestTablePlayCardDuringAuction(t *testing.T) {
	table := NewTable("test", nil)
	seatAllPlayers(t, table)
	table.Start(42)

	card := game.NewCard(game.Spades, game.Ace)
	err := table.PlayCard(game.North, card)
	if err == nil {
		t.Error("expected error playing card during auction phase")
	}
}

func TestTableBidDuringPlay(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)
	table.Start(42)
	for _, c := range clients {
		drainMessages(c)
	}

	// Complete the auction
	table.Bid(game.North, game.BidCall(1, game.NoTrump))
	table.Bid(game.East, game.Call{Type: game.Pass})
	table.Bid(game.South, game.Call{Type: game.Pass})
	table.Bid(game.West, game.Call{Type: game.Pass})
	for _, c := range clients {
		drainMessages(c)
	}

	// Now in play phase — bidding should fail
	if err := table.Bid(game.East, game.Call{Type: game.Pass}); err == nil {
		t.Error("expected error bidding during play phase")
	}
}

func TestTableBidWrongTurn(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)
	table.Start(42)
	for _, c := range clients {
		drainMessages(c)
	}

	// Board 1 dealer is North — East bidding first should fail
	if err := table.Bid(game.East, game.BidCall(1, game.NoTrump)); err == nil {
		t.Error("expected error for wrong turn")
	}
}

func TestTablePlayCardWrongTurn(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)
	table.Start(42)
	for _, c := range clients {
		drainMessages(c)
	}

	table.Bid(game.North, game.BidCall(1, game.NoTrump))
	table.Bid(game.East, game.Call{Type: game.Pass})
	table.Bid(game.South, game.Call{Type: game.Pass})
	table.Bid(game.West, game.Call{Type: game.Pass})
	for _, c := range clients {
		drainMessages(c)
	}

	// Opening leader is East — North playing should fail
	legal := table.game.LegalCards(game.East)
	if err := table.PlayCard(game.North, legal[0]); err == nil {
		t.Error("expected error for wrong turn")
	}
}

func TestTablePlayCardNotInHand(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)
	table.Start(42)
	for _, c := range clients {
		drainMessages(c)
	}

	table.Bid(game.North, game.BidCall(1, game.NoTrump))
	table.Bid(game.East, game.Call{Type: game.Pass})
	table.Bid(game.South, game.Call{Type: game.Pass})
	table.Bid(game.West, game.Call{Type: game.Pass})
	for _, c := range clients {
		drainMessages(c)
	}

	// Find a card East does NOT have
	eastHand := table.game.Play.RemainingHands[game.East]
	var missing game.Card
	for c := game.Card(0); c < game.NumCards; c++ {
		if !eastHand.Has(c) {
			missing = c
			break
		}
	}

	if err := table.PlayCard(game.East, missing); err == nil {
		t.Error("expected error playing a card not in hand")
	}
}

func TestTableSameClientTwoSeats(t *testing.T) {
	table := NewTable("test", nil)
	c := newTestClient()

	if err := table.Sit(c, game.North); err != nil {
		t.Fatalf("Sit(North): %v", err)
	}
	if err := table.Sit(c, game.East); err == nil {
		t.Error("expected error: same client should not sit in two seats")
	}
}

func TestTableRemoveClientNotAtTable(t *testing.T) {
	table := NewTable("test", nil)
	c := newTestClient()

	// Should not panic
	table.RemoveClient(c)
}

func TestTableStartSeedZero(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)

	// seed=0 should auto-generate a seed
	if err := table.Start(0); err != nil {
		t.Fatalf("Start(0): %v", err)
	}

	for _, c := range clients {
		msg := readServerMsg(t, c)
		if msg.Type != MsgGameState {
			t.Errorf("expected game_state, got %q", msg.Type)
		}
	}
}

func TestTableStartAfterComplete(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)
	table.Start(42)
	for _, c := range clients {
		drainMessages(c)
	}

	// Pass out the hand to reach PhaseComplete quickly
	table.Bid(game.North, game.Call{Type: game.Pass})
	table.Bid(game.East, game.Call{Type: game.Pass})
	table.Bid(game.South, game.Call{Type: game.Pass})
	table.Bid(game.West, game.Call{Type: game.Pass})
	for _, c := range clients {
		drainMessages(c)
	}

	if table.game.Phase != game.PhaseComplete {
		t.Fatalf("Phase = %v, want Complete", table.game.Phase)
	}

	// Starting a new game after completion should work
	if err := table.Start(99); err != nil {
		t.Errorf("Start after complete should succeed: %v", err)
	}
}

func TestTableBroadcastDifferentViews(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)
	table.Start(42)

	// Each player should see their own hand (13 cards) and they
	// should not all be identical (different players have different cards)
	type gameState struct {
		Hand []string `json:"hand"`
		Seat string   `json:"seat"`
	}

	hands := make(map[string]string) // seat -> first card
	for _, c := range clients {
		msg := readServerMsg(t, c)

		payloadBytes, _ := json.Marshal(msg.Payload)
		var gs gameState
		json.Unmarshal(payloadBytes, &gs)

		if len(gs.Hand) != 13 {
			t.Errorf("%s should have 13 cards, got %d", gs.Seat, len(gs.Hand))
		}
		hands[gs.Seat] = gs.Hand[0]
	}

	if len(hands) != 4 {
		t.Fatalf("expected 4 distinct seats, got %d", len(hands))
	}

	// At least two players should have different first cards
	allSame := true
	var prev string
	for _, card := range hands {
		if prev != "" && card != prev {
			allSame = false
			break
		}
		prev = card
	}
	if allSame {
		t.Error("all players received the same hand — ViewFor is not personalizing")
	}
}

func TestTableFullGame(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)

	if err := table.Start(42); err != nil {
		t.Fatalf("Start: %v", err)
	}
	for _, c := range clients {
		drainMessages(c)
	}

	// Auction: North opens 1NT, all pass
	if err := table.Bid(game.North, game.BidCall(1, game.NoTrump)); err != nil {
		t.Fatalf("Bid 1NT: %v", err)
	}
	for _, c := range clients {
		drainMessages(c)
	}

	for _, dir := range []game.Direction{game.East, game.South, game.West} {
		if err := table.Bid(dir, game.Call{Type: game.Pass}); err != nil {
			t.Fatalf("Bid Pass(%v): %v", dir, err)
		}
		for _, c := range clients {
			drainMessages(c)
		}
	}

	// Play all 13 tricks using the game's LegalCards
	for trick := 0; trick < 13; trick++ {
		for card := 0; card < 4; card++ {
			turn, ok := table.game.Turn()
			if !ok {
				t.Fatalf("trick %d card %d: game ended early", trick, card)
			}

			actor := turn
			if turn == table.game.Play.Dummy {
				actor = table.game.Play.Declarer
			}

			legal := table.game.LegalCards(actor)
			if len(legal) == 0 {
				t.Fatalf("trick %d card %d: no legal cards", trick, card)
			}

			if err := table.PlayCard(actor, legal[0]); err != nil {
				t.Fatalf("PlayCard(%v, %v): %v", actor, legal[0], err)
			}
			for _, c := range clients {
				drainMessages(c)
			}
		}
	}

	// Game should be complete — last broadcast should have a result
	if table.game.Phase != game.PhaseComplete {
		t.Errorf("Phase = %v, want Complete", table.game.Phase)
	}
	if table.game.Result == nil {
		t.Error("Result should not be nil")
	}

	t.Logf("Game result: contract=%v, score=%d", table.game.Result.Contract, table.game.Result.Score)
}

func passOutHand(t *testing.T, table *Table, clients [game.NumDirections]*Client) {
	t.Helper()
	dealer := table.game.Board.Dealer
	for i := range game.NumDirections {
		dir := game.Direction((int(dealer) + i) % game.NumDirections)
		if err := table.Bid(dir, game.Call{Type: game.Pass}); err != nil {
			t.Fatalf("Bid Pass(%v): %v", dir, err)
		}
	}
	for _, c := range clients {
		drainMessages(c)
	}
}

func TestTableBoardNumberIncrements(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)

	for board := 1; board <= 4; board++ {
		if table.boardNum != board {
			t.Fatalf("before game %d: boardNum = %d, want %d", board, table.boardNum, board)
		}

		table.Start(int64(board * 100))
		for _, c := range clients {
			drainMessages(c)
		}

		if table.game.Board.Number != board {
			t.Errorf("game %d: Board.Number = %d, want %d", board, table.game.Board.Number, board)
		}

		passOutHand(t, table, clients)

		if table.game.Phase != game.PhaseComplete {
			t.Fatalf("game %d: Phase = %v, want Complete", board, table.game.Phase)
		}
	}

	if table.boardNum != 5 {
		t.Errorf("after 4 games: boardNum = %d, want 5", table.boardNum)
	}
}

func TestTableBoardNumberAffectsDealerRotation(t *testing.T) {
	table := NewTable("test", nil)
	clients := seatAllPlayers(t, table)

	expectedDealers := []game.Direction{game.North, game.East, game.South, game.West}

	for i, wantDealer := range expectedDealers {
		table.Start(int64(i + 1))
		for _, c := range clients {
			drainMessages(c)
		}

		if table.game.Board.Dealer != wantDealer {
			t.Errorf("board %d: Dealer = %v, want %v", i+1, table.game.Board.Dealer, wantDealer)
		}

		passOutHand(t, table, clients)
	}
}
