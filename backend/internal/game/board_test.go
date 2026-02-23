package game

import "testing"

func TestDealerForBoard(t *testing.T) {
	expected := []Direction{North, East, South, West, North, East}
	for i, want := range expected {
		got := DealerForBoard(i + 1)
		if got != want {
			t.Errorf("DealerForBoard(%d) = %v, want %v", i+1, got, want)
		}
	}
}

func TestVulnerabilityForBoard(t *testing.T) {
	tests := []struct {
		board int
		want  Vulnerability
	}{
		{1, VulNone}, {2, VulNS}, {3, VulEW}, {4, VulBoth},
		{5, VulNS}, {6, VulEW}, {7, VulBoth}, {8, VulNone},
		{9, VulEW}, {10, VulBoth}, {11, VulNone}, {12, VulNS},
		{13, VulBoth}, {14, VulNone}, {15, VulNS}, {16, VulEW},
		// Cycle repeats
		{17, VulNone}, {32, VulEW},
	}
	for _, tt := range tests {
		got := VulnerabilityForBoard(tt.board)
		if got != tt.want {
			t.Errorf("VulnerabilityForBoard(%d) = %v, want %v", tt.board, got, tt.want)
		}
	}
}

func TestDealFromSeed(t *testing.T) {
	hands := DealFromSeed(42)

	// Every hand should have exactly 13 cards
	for d := Direction(0); d < NumDirections; d++ {
		if hands[d].Len() != 13 {
			t.Errorf("%v has %d cards, want 13", d, hands[d].Len())
		}
	}

	// All 52 cards should be dealt exactly once
	var all Hand
	for d := Direction(0); d < NumDirections; d++ {
		overlap := all & hands[d]
		if overlap != 0 {
			t.Errorf("%v has cards already dealt to another hand", d)
		}
		all |= hands[d]
	}
	if all.Len() != NumCards {
		t.Errorf("total cards dealt = %d, want %d", all.Len(), NumCards)
	}
}

func TestDealFromSeedDeterministic(t *testing.T) {
	hands1 := DealFromSeed(12345)
	hands2 := DealFromSeed(12345)

	for d := Direction(0); d < NumDirections; d++ {
		if hands1[d] != hands2[d] {
			t.Errorf("same seed produced different hands for %v", d)
		}
	}
}

func TestDealFromSeedDifferentSeeds(t *testing.T) {
	hands1 := DealFromSeed(1)
	hands2 := DealFromSeed(2)

	same := true
	for d := Direction(0); d < NumDirections; d++ {
		if hands1[d] != hands2[d] {
			same = false
			break
		}
	}
	if same {
		t.Error("different seeds produced identical deals")
	}
}

func TestDealPBNRoundTrip(t *testing.T) {
	for _, seed := range []int64{0, 1, 42, 99999, -1} {
		b := NewBoard(1, seed)
		pbn := b.ToPBN()

		dealer, hands, ok := PBNToBoard(pbn)
		if !ok {
			t.Errorf("PBNToBoard failed for seed %d: %q", seed, pbn)
			continue
		}
		if dealer != b.Dealer {
			t.Errorf("seed %d: dealer = %v, want %v", seed, dealer, b.Dealer)
		}
		for d := Direction(0); d < NumDirections; d++ {
			if hands[d] != b.Hands[d] {
				t.Errorf("seed %d: %v hand mismatch after round-trip", seed, d)
			}
		}
	}
}

func TestNewBoard(t *testing.T) {
	b := NewBoard(5, 100)

	if b.Number != 5 {
		t.Errorf("Number = %d, want 5", b.Number)
	}
	if b.Seed != 100 {
		t.Errorf("Seed = %d, want 100", b.Seed)
	}
	if b.Dealer != North {
		t.Errorf("Dealer = %v, want North (board 5)", b.Dealer)
	}
	if b.Vulnerability != VulNS {
		t.Errorf("Vulnerability = %v, want NS (board 5)", b.Vulnerability)
	}

	total := 0
	for d := Direction(0); d < NumDirections; d++ {
		total += b.Hands[d].Len()
	}
	if total != NumCards {
		t.Errorf("total cards = %d, want %d", total, NumCards)
	}
}

func TestPBNToBoardInvalid(t *testing.T) {
	invalids := []string{
		"",
		"X:AK.QJ.T9.87 65.43.2A.KQ 65.43.2A.KQ 65.43.2A.KQ",
		"N AK.QJ.T9.87 65.43.2A.KQ 65.43.2A.KQ 65.43.2A.KQ", // missing colon
		"N:",          // no hands
		"N:AK.QJ.T9", // only 1 hand
	}
	for _, s := range invalids {
		if _, _, ok := PBNToBoard(s); ok {
			t.Errorf("PBNToBoard(%q) should have failed", s)
		}
	}
}

func TestDealDistribution(t *testing.T) {
	// Verify deals look random: across many seeds, every card should appear
	// in each direction at least once.
	seen := [NumDirections][NumCards]bool{}
	for seed := int64(0); seed < 200; seed++ {
		hands := DealFromSeed(seed)
		for d := Direction(0); d < NumDirections; d++ {
			for c := Card(0); c < NumCards; c++ {
				if hands[d].Has(c) {
					seen[d][c] = true
				}
			}
		}
	}
	for d := Direction(0); d < NumDirections; d++ {
		for c := Card(0); c < NumCards; c++ {
			if !seen[d][c] {
				t.Errorf("card %v never dealt to %v across 200 deals", c, d)
			}
		}
	}
}
