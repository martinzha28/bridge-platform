package game

import "testing"

func TestViewForDealPhase(t *testing.T) {
	g := NewGame(1)
	v := g.ViewFor(North)

	if v.Phase != "Deal" {
		t.Errorf("Phase = %q, want Deal", v.Phase)
	}
	if v.Seat != "North" {
		t.Errorf("Seat = %q, want North", v.Seat)
	}
	if len(v.Hand) != 0 {
		t.Errorf("Hand should be empty during deal, got %d cards", len(v.Hand))
	}
	if v.Turn != "" {
		t.Errorf("Turn should be empty during deal, got %q", v.Turn)
	}
}

func TestViewForAuctionPhase(t *testing.T) {
	g := newDealtGame(t, 1, 42)
	v := g.ViewFor(North)

	if v.Phase != "Auction" {
		t.Errorf("Phase = %q, want Auction", v.Phase)
	}
	if len(v.Hand) != 13 {
		t.Errorf("Hand should have 13 cards, got %d", len(v.Hand))
	}
	if v.Turn != "North" {
		t.Errorf("Turn = %q, want North", v.Turn)
	}
	if len(v.LegalCalls) != 36 {
		t.Errorf("LegalCalls = %d, want 36", len(v.LegalCalls))
	}

	// East (not their turn) should have no legal calls
	ev := g.ViewFor(East)
	if len(ev.LegalCalls) != 0 {
		t.Errorf("East should have 0 legal calls, got %d", len(ev.LegalCalls))
	}
}

func TestViewForPlayPhase(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	sessionMustCall(t, g, North, BidCall(1, NoTrump))
	sessionMustCall(t, g, East, passCall)
	sessionMustCall(t, g, South, passCall)
	sessionMustCall(t, g, West, passCall)

	v := g.ViewFor(East)
	if v.Phase != "Play" {
		t.Errorf("Phase = %q, want Play", v.Phase)
	}
	if v.Dummy != "South" {
		t.Errorf("Dummy = %q, want South", v.Dummy)
	}
	if v.DummyHand != nil {
		t.Error("DummyHand should be nil before opening lead")
	}
	if len(v.LegalCards) != 13 {
		t.Errorf("East (leader) should have 13 legal cards, got %d", len(v.LegalCards))
	}
	if v.Contract != "1NT by North" {
		t.Errorf("Contract = %q, want 1NT by North", v.Contract)
	}

	// Play opening lead
	legal := g.LegalCards(East)
	sessionMustPlay(t, g, East, legal[0])

	v = g.ViewFor(North)
	if v.DummyHand == nil {
		t.Error("DummyHand should be visible after opening lead")
	}
	if len(v.CurrentTrick) != 1 {
		t.Errorf("CurrentTrick should have 1 card, got %d", len(v.CurrentTrick))
	}
}

func TestViewLastTrick(t *testing.T) {
	g := newDealtGame(t, 1, 42)
	sessionMustCall(t, g, North, BidCall(1, NoTrump))
	sessionMustCall(t, g, East, passCall)
	sessionMustCall(t, g, South, passCall)
	sessionMustCall(t, g, West, passCall)

	playCards(t, g, 3) // three into the first trick
	if len(g.ViewFor(North).LastTrick) != 0 {
		t.Fatal("LastTrick should be empty mid-trick")
	}

	playCards(t, g, 1) // complete the trick
	v := g.ViewFor(North)
	if len(v.LastTrick) != 4 {
		t.Fatalf("LastTrick = %d cards, want 4 after a completed trick", len(v.LastTrick))
	}
	if len(v.CurrentTrick) != 0 {
		t.Errorf("CurrentTrick = %d, want 0 between tricks", len(v.CurrentTrick))
	}
	if v.LastTrick[0].Seat != "East" {
		t.Errorf("LastTrick[0].Seat = %q, want East (opening leader)", v.LastTrick[0].Seat)
	}

	playCards(t, g, 1) // lead the next trick
	if len(g.ViewFor(North).LastTrick) != 0 {
		t.Error("LastTrick should clear once the next trick is led")
	}
}

// playCards plays n legal cards in turn order (declarer plays for dummy).
func playCards(t *testing.T, g *Game, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		turn, ok := g.Turn()
		if !ok {
			t.Fatalf("playCards: game ended after %d cards", i)
		}
		actor := turn
		if turn == g.Play.Dummy {
			actor = g.Play.Declarer
		}
		sessionMustPlay(t, g, actor, g.LegalCards(actor)[0])
	}
}

func TestViewForCompletePhase(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	sessionMustCall(t, g, North, BidCall(1, NoTrump))
	sessionMustCall(t, g, East, passCall)
	sessionMustCall(t, g, South, passCall)
	sessionMustCall(t, g, West, passCall)
	playAllTricks(t, g)

	v := g.ViewFor(North)
	if v.Phase != "Complete" {
		t.Errorf("Phase = %q, want Complete", v.Phase)
	}
	if v.Result == nil {
		t.Fatal("Result should not be nil")
	}
	if v.Result.TricksNS+v.Result.TricksEW != 13 {
		t.Errorf("Tricks don't add up: NS=%d + EW=%d", v.Result.TricksNS, v.Result.TricksEW)
	}
	if len(v.Hand) != 13 {
		t.Errorf("Complete phase should show original 13 cards, got %d", len(v.Hand))
	}
	if len(v.LegalCalls) != 0 {
		t.Errorf("should have no legal calls after game, got %d", len(v.LegalCalls))
	}
}

func TestViewCallHistory(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	sessionMustCall(t, g, North, BidCall(1, HeartStrain))
	sessionMustCall(t, g, East, passCall)

	v := g.ViewFor(South)
	if len(v.Calls) != 2 {
		t.Fatalf("Calls should have 2 entries, got %d", len(v.Calls))
	}
	if v.Calls[0] != "1H" {
		t.Errorf("Calls[0] = %q, want 1H", v.Calls[0])
	}
	if v.Calls[1] != "P" {
		t.Errorf("Calls[1] = %q, want P", v.Calls[1])
	}
}

func TestViewPassedOut(t *testing.T) {
	g := newDealtGame(t, 1, 42)

	for _, dir := range []Direction{North, East, South, West} {
		sessionMustCall(t, g, dir, passCall)
	}

	v := g.ViewFor(North)
	if v.Result == nil {
		t.Fatal("Result should not be nil")
	}
	if !v.Result.PassedOut {
		t.Error("Result.PassedOut should be true")
	}
}
