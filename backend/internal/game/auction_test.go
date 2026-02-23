package game

import "testing"

func TestAuctionSimpleBidding(t *testing.T) {
	// N opens 1H, everyone passes -> contract 1H by North
	a := NewAuction(North)

	calls := []struct {
		dir  Direction
		call Call
	}{
		{North, BidCall(1, HeartStrain)},
		{East, passCall},
		{South, passCall},
		{West, passCall},
	}
	for _, c := range calls {
		if err := a.MakeCall(c.dir, c.call); err != nil {
			t.Fatalf("MakeCall(%v, %v) failed: %v", c.dir, c.call, err)
		}
	}

	if !a.IsFinished() {
		t.Fatal("auction should be finished")
	}
	if a.PassedOut() {
		t.Fatal("auction should not be passed out")
	}

	contract, ok := a.Contract()
	if !ok {
		t.Fatal("Contract() should succeed")
	}
	if contract.Level != 1 || contract.Strain != HeartStrain {
		t.Errorf("contract = %d%s, want 1H", contract.Level, contract.Strain)
	}
	if contract.Declarer != North {
		t.Errorf("declarer = %v, want North", contract.Declarer)
	}
	if contract.Doubled || contract.Redoubled {
		t.Error("contract should not be doubled")
	}
}

func TestAuctionPassedOut(t *testing.T) {
	a := NewAuction(East)

	for _, dir := range []Direction{East, South, West, North} {
		if err := a.MakeCall(dir, passCall); err != nil {
			t.Fatalf("MakeCall(%v, Pass) failed: %v", dir, err)
		}
	}

	if !a.IsFinished() {
		t.Fatal("auction should be finished")
	}
	if !a.PassedOut() {
		t.Fatal("auction should be passed out")
	}
	if _, ok := a.Contract(); ok {
		t.Error("passed-out auction should not have a contract")
	}
}

func TestAuctionCompetitiveBidding(t *testing.T) {
	// N:1C - 1S - 2H - P - 4H - P - P - P -> 4H by North (first to bid hearts on NS side)
	a := NewAuction(North)

	calls := []struct {
		dir  Direction
		call Call
	}{
		{North, BidCall(1, ClubStrain)},
		{East, BidCall(1, SpadeStrain)},
		{South, BidCall(2, HeartStrain)},
		{West, passCall},
		{North, BidCall(4, HeartStrain)},
		{East, passCall},
		{South, passCall},
		{West, passCall},
	}
	for _, c := range calls {
		if err := a.MakeCall(c.dir, c.call); err != nil {
			t.Fatalf("MakeCall(%v, %v) failed: %v", c.dir, c.call, err)
		}
	}

	contract, ok := a.Contract()
	if !ok {
		t.Fatal("Contract() should succeed")
	}
	if contract.Level != 4 || contract.Strain != HeartStrain {
		t.Errorf("contract = %d%s, want 4H", contract.Level, contract.Strain)
	}
	// South bid hearts first on the NS side
	if contract.Declarer != South {
		t.Errorf("declarer = %v, want South (first to bid hearts on NS side)", contract.Declarer)
	}

	want := "1C 1S 2H P 4H P P P"
	if got := a.ToNotation(); got != want {
		t.Errorf("notation = %q, want %q", got, want)
	}
}

func TestAuctionDoubled(t *testing.T) {
	// N:1NT - X - P - P - P -> 1NT doubled by North
	a := NewAuction(North)

	calls := []struct {
		dir  Direction
		call Call
	}{
		{North, BidCall(1, NoTrump)},
		{East, doubleCall},
		{South, passCall},
		{West, passCall},
		{North, passCall},
	}
	for _, c := range calls {
		if err := a.MakeCall(c.dir, c.call); err != nil {
			t.Fatalf("MakeCall(%v, %v) failed: %v", c.dir, c.call, err)
		}
	}

	contract, ok := a.Contract()
	if !ok {
		t.Fatal("Contract() should succeed")
	}
	if !contract.Doubled {
		t.Error("contract should be doubled")
	}
	if contract.Redoubled {
		t.Error("contract should not be redoubled")
	}
	if contract.Declarer != North {
		t.Errorf("declarer = %v, want North", contract.Declarer)
	}
}

func TestAuctionRedoubled(t *testing.T) {
	// N:1H - X - XX - P - P - P -> 1H redoubled by North
	a := NewAuction(North)

	calls := []struct {
		dir  Direction
		call Call
	}{
		{North, BidCall(1, HeartStrain)},
		{East, doubleCall},
		{South, redoubleCall},
		{West, passCall},
		{North, passCall},
		{East, passCall},
	}
	for _, c := range calls {
		if err := a.MakeCall(c.dir, c.call); err != nil {
			t.Fatalf("MakeCall(%v, %v) failed: %v", c.dir, c.call, err)
		}
	}

	contract, ok := a.Contract()
	if !ok {
		t.Fatal("Contract() should succeed")
	}
	if !contract.Redoubled {
		t.Error("contract should be redoubled")
	}
}

func TestAuctionNewBidCancelsDouble(t *testing.T) {
	// N:1H - X - 2H - P - P - P -> double is cancelled by the new bid
	a := NewAuction(North)

	must(t, a.MakeCall(North, BidCall(1, HeartStrain)))
	must(t, a.MakeCall(East, doubleCall))
	must(t, a.MakeCall(South, BidCall(2, HeartStrain)))
	must(t, a.MakeCall(West, passCall))
	must(t, a.MakeCall(North, passCall))
	must(t, a.MakeCall(East, passCall))

	contract, ok := a.Contract()
	if !ok {
		t.Fatal("Contract() should succeed")
	}
	if contract.Doubled {
		t.Error("contract should not be doubled (new bid cancels double)")
	}
	if contract.Declarer != North {
		t.Errorf("declarer = %v, want North (first to bid hearts on NS)", contract.Declarer)
	}
}

func TestAuctionWrongTurn(t *testing.T) {
	a := NewAuction(North)
	err := a.MakeCall(East, passCall)
	if err == nil {
		t.Error("expected error for wrong turn")
	}
}

func TestAuctionBidTooLow(t *testing.T) {
	a := NewAuction(North)
	must(t, a.MakeCall(North, BidCall(2, HeartStrain)))

	err := a.MakeCall(East, BidCall(1, SpadeStrain))
	if err == nil {
		t.Error("expected error: 1S is lower than 2H")
	}

	err = a.MakeCall(East, BidCall(2, DiamondStrain))
	if err == nil {
		t.Error("expected error: 2D is lower than 2H")
	}

	// Same level, higher strain should work
	if err := a.MakeCall(East, BidCall(2, SpadeStrain)); err != nil {
		t.Errorf("2S should be legal after 2H: %v", err)
	}
}

func TestAuctionBidOutOfRange(t *testing.T) {
	a := NewAuction(North)

	if err := a.MakeCall(North, BidCall(0, ClubStrain)); err == nil {
		t.Error("expected error for level 0")
	}
	if err := a.MakeCall(North, BidCall(8, ClubStrain)); err == nil {
		t.Error("expected error for level 8")
	}
}

func TestAuctionDoubleOwnBid(t *testing.T) {
	a := NewAuction(North)
	must(t, a.MakeCall(North, BidCall(1, ClubStrain)))
	must(t, a.MakeCall(East, passCall))

	// South is North's partner — can't double own side's bid
	err := a.MakeCall(South, doubleCall)
	if err == nil {
		t.Error("expected error: cannot double own side's bid")
	}
}

func TestAuctionDoubleNoBid(t *testing.T) {
	a := NewAuction(North)
	err := a.MakeCall(North, doubleCall)
	if err == nil {
		t.Error("expected error: cannot double with no bid")
	}
}

func TestAuctionRedoubleWithoutDouble(t *testing.T) {
	a := NewAuction(North)
	must(t, a.MakeCall(North, BidCall(1, ClubStrain)))

	err := a.MakeCall(East, redoubleCall)
	if err == nil {
		t.Error("expected error: cannot redouble without double")
	}
}

func TestAuctionRedoubleWrongSide(t *testing.T) {
	// N:1C - X - P - P - P -> West (same side as East who doubled) can't redouble
	a := NewAuction(North)
	must(t, a.MakeCall(North, BidCall(1, ClubStrain)))
	must(t, a.MakeCall(East, doubleCall))
	must(t, a.MakeCall(South, passCall))

	err := a.MakeCall(West, redoubleCall)
	if err == nil {
		t.Error("expected error: wrong side can't redouble")
	}
}

func TestAuctionCallAfterFinished(t *testing.T) {
	a := NewAuction(North)
	must(t, a.MakeCall(North, passCall))
	must(t, a.MakeCall(East, passCall))
	must(t, a.MakeCall(South, passCall))
	must(t, a.MakeCall(West, passCall))

	err := a.MakeCall(North, passCall)
	if err == nil {
		t.Error("expected error: auction is finished")
	}
}

func TestAuctionGrandSlam(t *testing.T) {
	a := NewAuction(South)
	must(t, a.MakeCall(South, BidCall(7, NoTrump)))
	must(t, a.MakeCall(West, passCall))
	must(t, a.MakeCall(North, passCall))
	must(t, a.MakeCall(East, passCall))

	contract, ok := a.Contract()
	if !ok {
		t.Fatal("Contract() should succeed")
	}
	if contract.Level != 7 || contract.Strain != NoTrump {
		t.Errorf("contract = %d%s, want 7NT", contract.Level, contract.Strain)
	}
	if contract.Declarer != South {
		t.Errorf("declarer = %v, want South", contract.Declarer)
	}
}

func TestAuctionDoubleAfterPassOverBid(t *testing.T) {
	// N:1H - P - P - X -> West can double (opponent's bid, passes don't change that)
	a := NewAuction(North)
	must(t, a.MakeCall(North, BidCall(1, HeartStrain)))
	must(t, a.MakeCall(East, passCall))
	must(t, a.MakeCall(South, passCall))

	if err := a.MakeCall(West, doubleCall); err != nil {
		t.Errorf("West should be able to double North's bid: %v", err)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
