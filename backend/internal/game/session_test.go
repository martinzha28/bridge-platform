package game

import "testing"

func TestNewGame(t *testing.T) {
	g := NewGame(1)

	if g.Phase != PhaseDeal {
		t.Errorf("Phase = %v, want Deal", g.Phase)
	}
	if g.Board.Dealer != North {
		t.Errorf("Dealer = %v, want North", g.Board.Dealer)
	}
	if g.Auction != nil {
		t.Error("Auction should be nil before deal")
	}
	if _, ok := g.Turn(); ok {
		t.Error("Turn should return false during deal phase")
	}
}

func TestDeal(t *testing.T) {
	g := NewGame(1)
	if err := g.Deal(42); err != nil {
		t.Fatalf("Deal: %v", err)
	}

	if g.Phase != PhaseAuction {
		t.Errorf("Phase = %v, want Auction", g.Phase)
	}
	if g.Auction == nil {
		t.Fatal("Auction is nil after deal")
	}

	turn, ok := g.Turn()
	if !ok || turn != North {
		t.Errorf("Turn = %v, %v; want North, true", turn, ok)
	}

	// Can't deal again
	if err := g.Deal(99); err == nil {
		t.Error("expected error dealing twice")
	}
}

func TestSetHands(t *testing.T) {
	g := NewGame(1)

	hands := DealFromSeed(42)
	if err := g.SetHands(hands); err != nil {
		t.Fatalf("SetHands: %v", err)
	}

	if g.Phase != PhaseAuction {
		t.Errorf("Phase = %v, want Auction", g.Phase)
	}
	for i, h := range g.Board.Hands {
		if h != hands[i] {
			t.Errorf("Hand[%d] mismatch", i)
		}
	}
}

func TestSetHandsValidation(t *testing.T) {
	g := NewGame(1)

	// Wrong number of cards
	badHands := [NumDirections]Hand{}
	badHands[North] = HandFromCards([]Card{NewCard(Spades, Ace)})
	if err := g.SetHands(badHands); err == nil {
		t.Error("expected error for wrong card count")
	}

	// Duplicate cards across hands
	full13 := DealFromSeed(42)
	dup := full13
	dup[East] = dup[North] // same 13 cards in two hands
	g2 := NewGame(1)
	if err := g2.SetHands(dup); err == nil {
		t.Error("expected error for duplicate cards")
	}
}

func TestDealPhaseEnforcement(t *testing.T) {
	g := NewGame(1)

	if err := g.MakeCall(North, passCall); err == nil {
		t.Error("expected error making call during deal phase")
	}
	if err := g.PlayCard(North, NewCard(Spades, Ace)); err == nil {
		t.Error("expected error playing card during deal phase")
	}
	if calls := g.LegalCalls(North); calls != nil {
		t.Error("expected nil legal calls during deal phase")
	}
	if cards := g.LegalCards(North); cards != nil {
		t.Error("expected nil legal cards during deal phase")
	}
}

func TestPassedOut(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	for _, dir := range []Direction{North, East, South, West} {
		sessionMustCall(t, g, dir, passCall)
	}

	if g.Phase != PhaseComplete {
		t.Errorf("Phase = %v, want Complete", g.Phase)
	}
	if g.Result == nil || !g.Result.PassedOut() {
		t.Error("expected passed-out result")
	}
	if g.Play != nil {
		t.Error("Play should be nil for passed-out hand")
	}
	if _, ok := g.Turn(); ok {
		t.Error("Turn should return false after completion")
	}
}

func TestFullGame(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	sessionMustCall(t, g, North, BidCall(1, NoTrump))
	sessionMustCall(t, g, East, passCall)
	sessionMustCall(t, g, South, passCall)
	sessionMustCall(t, g, West, passCall)

	if g.Phase != PhasePlay {
		t.Fatalf("Phase = %v, want Play", g.Phase)
	}

	c := g.Play.Contract
	if c.Level != 1 || c.Strain != NoTrump || c.Declarer != North {
		t.Fatalf("Contract = %v, want 1NT by North", c)
	}

	turn, _ := g.Turn()
	if turn != East {
		t.Fatalf("Opening leader = %v, want East", turn)
	}

	playAllTricks(t, g)

	if g.Phase != PhaseComplete {
		t.Fatalf("Phase = %v, want Complete", g.Phase)
	}
	if g.Result == nil {
		t.Fatal("Result is nil")
	}
	if g.Result.PassedOut() {
		t.Error("should not be passed out")
	}
	if g.Result.TricksNS+g.Result.TricksEW != 13 {
		t.Errorf("NS=%d + EW=%d != 13", g.Result.TricksNS, g.Result.TricksEW)
	}
	if g.Result.Contract == nil {
		t.Fatal("Result.Contract is nil")
	}

	t.Logf("Contract: %v", g.Result.Contract)
	t.Logf("Tricks: NS=%d EW=%d", g.Result.TricksNS, g.Result.TricksEW)
	t.Logf("Score: %d", g.Result.Score)
}

func TestPhaseEnforcement(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	if err := g.PlayCard(North, NewCard(Spades, Ace)); err == nil {
		t.Error("expected error playing card during auction")
	}

	sessionMustCall(t, g, North, BidCall(1, NoTrump))
	sessionMustCall(t, g, East, passCall)
	sessionMustCall(t, g, South, passCall)
	sessionMustCall(t, g, West, passCall)

	if err := g.MakeCall(East, passCall); err == nil {
		t.Error("expected error making call during play")
	}
}

func TestLegalCallsOpening(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	// Before any bids: pass + 35 bids (7 levels * 5 strains)
	calls := g.LegalCalls(North)
	if len(calls) != 36 {
		t.Errorf("len(LegalCalls) = %d, want 36", len(calls))
	}
	if calls[0].Type != Pass {
		t.Error("first legal call should be Pass")
	}
}

func TestLegalCallsAfterBid(t *testing.T) {
	g := newDealtGame(t, 1, 42)
	sessionMustCall(t, g, North, BidCall(1, ClubStrain))

	calls := g.LegalCalls(East)

	hasPass, hasDouble := false, false
	bidCount := 0
	for _, c := range calls {
		switch c.Type {
		case Pass:
			hasPass = true
		case Double:
			hasDouble = true
		case Bid:
			bidCount++
		}
	}

	if !hasPass {
		t.Error("expected pass in legal calls")
	}
	if !hasDouble {
		t.Error("expected double in legal calls for opponent")
	}
	// 1D through 7NT = 34 bids
	if bidCount != 34 {
		t.Errorf("bidCount = %d, want 34", bidCount)
	}
}

func TestLegalCallsWrongTurn(t *testing.T) {
	g := newDealtGame(t, 1, 42)
	sessionMustCall(t, g, North, BidCall(1, ClubStrain))

	if calls := g.LegalCalls(North); calls != nil {
		t.Errorf("expected nil for wrong turn, got %d calls", len(calls))
	}
}

func TestLegalCardsLeader(t *testing.T) {
	g := newDealtGame(t, 1, 42)
	sessionMustCall(t, g, North, BidCall(1, ClubStrain))
	sessionMustCall(t, g, East, passCall)
	sessionMustCall(t, g, South, passCall)
	sessionMustCall(t, g, West, passCall)

	legal := g.LegalCards(East)
	if len(legal) != 13 {
		t.Errorf("leader should have 13 legal cards, got %d", len(legal))
	}
}

func TestLegalCardsFollowSuit(t *testing.T) {
	g := newDealtGame(t, 1, 42)
	sessionMustCall(t, g, North, BidCall(1, ClubStrain))
	sessionMustCall(t, g, East, passCall)
	sessionMustCall(t, g, South, passCall)
	sessionMustCall(t, g, West, passCall)

	// East leads
	legal := g.LegalCards(East)
	leadCard := legal[0]
	sessionMustPlay(t, g, East, leadCard)

	// Next player must follow suit if they can
	turn, _ := g.Turn()
	actor := turn
	if turn == g.Play.Dummy {
		actor = g.Play.Declarer
	}

	legal = g.LegalCards(actor)
	ledSuit := leadCard.Suit()
	hand := g.Play.RemainingHands[turn]

	if hand.HasSuit(ledSuit) {
		for _, c := range legal {
			if c.Suit() != ledSuit {
				t.Errorf("must follow suit %v but %v is legal", ledSuit, c)
			}
		}
	}
}

func TestDummyVisible(t *testing.T) {
	g := NewGame(1)

	if g.DummyVisible() {
		t.Error("dummy should not be visible during deal")
	}

	if err := g.Deal(42); err != nil {
		t.Fatalf("Deal: %v", err)
	}

	if g.DummyVisible() {
		t.Error("dummy should not be visible during auction")
	}

	sessionMustCall(t, g, North, BidCall(1, NoTrump))
	sessionMustCall(t, g, East, passCall)
	sessionMustCall(t, g, South, passCall)
	sessionMustCall(t, g, West, passCall)

	if g.DummyVisible() {
		t.Error("dummy should not be visible before opening lead")
	}

	legal := g.LegalCards(East)
	sessionMustPlay(t, g, East, legal[0])

	if !g.DummyVisible() {
		t.Error("dummy should be visible after opening lead")
	}
}

func TestDoubledContract(t *testing.T) {
	// Board 2: dealer East, vul NS
	g := newDealtGame(t, 2, 99)

	sessionMustCall(t, g, East, BidCall(1, SpadeStrain))
	sessionMustCall(t, g, South, doubleCall)
	sessionMustCall(t, g, West, passCall)
	sessionMustCall(t, g, North, passCall)
	sessionMustCall(t, g, East, passCall)

	if g.Phase != PhasePlay {
		t.Fatalf("Phase = %v, want Play", g.Phase)
	}

	c := g.Play.Contract
	if !c.Doubled {
		t.Error("contract should be doubled")
	}
	if c.Declarer != East {
		t.Errorf("Declarer = %v, want East", c.Declarer)
	}

	playAllTricks(t, g)

	if g.Phase != PhaseComplete {
		t.Fatalf("Phase = %v, want Complete", g.Phase)
	}
	if !g.Result.Contract.Doubled {
		t.Error("result contract should be doubled")
	}

	t.Logf("Doubled contract: %v, score: %d", g.Result.Contract, g.Result.Score)
}

func TestCompletedGameActions(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	for _, dir := range []Direction{North, East, South, West} {
		sessionMustCall(t, g, dir, passCall)
	}

	if err := g.MakeCall(North, passCall); err == nil {
		t.Error("expected error making call after game complete")
	}
	if err := g.PlayCard(North, NewCard(Spades, Ace)); err == nil {
		t.Error("expected error playing card after game complete")
	}
	if calls := g.LegalCalls(North); calls != nil {
		t.Error("expected nil legal calls after game complete")
	}
	if cards := g.LegalCards(North); cards != nil {
		t.Error("expected nil legal cards after game complete")
	}
}

// --- helpers ---

func newDealtGame(t *testing.T, boardNum int, seed int64) *Game {
	t.Helper()
	g := NewGame(boardNum)
	if err := g.Deal(seed); err != nil {
		t.Fatalf("Deal(%d): %v", seed, err)
	}
	return g
}

func playAllTricks(t *testing.T, g *Game) {
	t.Helper()
	for trick := 0; trick < NumTricks; trick++ {
		for card := 0; card < NumDirections; card++ {
			turn, ok := g.Turn()
			if !ok {
				t.Fatalf("trick %d card %d: game ended early", trick, card)
			}
			actor := turn
			if turn == g.Play.Dummy {
				actor = g.Play.Declarer
			}
			legal := g.LegalCards(actor)
			if len(legal) == 0 {
				t.Fatalf("trick %d card %d: no legal cards for %v (turn=%v)", trick, card, actor, turn)
			}
			sessionMustPlay(t, g, actor, legal[0])
		}
	}
}

func sessionMustCall(t *testing.T, g *Game, dir Direction, call Call) {
	t.Helper()
	if err := g.MakeCall(dir, call); err != nil {
		t.Fatalf("MakeCall(%v, %v): %v", dir, call, err)
	}
}

func sessionMustPlay(t *testing.T, g *Game, actor Direction, card Card) {
	t.Helper()
	if err := g.PlayCard(actor, card); err != nil {
		t.Fatalf("PlayCard(%v, %v): %v", actor, card, err)
	}
}
