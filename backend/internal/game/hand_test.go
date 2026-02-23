package game

import "testing"

func TestHandAddRemoveHas(t *testing.T) {
	var h Hand
	sa := NewCard(Spades, Ace)
	hk := NewCard(Hearts, King)

	h = h.Add(sa)
	if !h.Has(sa) {
		t.Error("hand should contain SA after Add")
	}
	if h.Len() != 1 {
		t.Errorf("hand length = %d, want 1", h.Len())
	}

	h = h.Add(hk)
	if h.Len() != 2 {
		t.Errorf("hand length = %d, want 2", h.Len())
	}

	h = h.Remove(sa)
	if h.Has(sa) {
		t.Error("hand should not contain SA after Remove")
	}
	if !h.Has(hk) {
		t.Error("hand should still contain HK after removing SA")
	}
	if h.Len() != 1 {
		t.Errorf("hand length = %d, want 1", h.Len())
	}
}

func TestHandFromCards(t *testing.T) {
	cards := []Card{
		NewCard(Spades, Ace), NewCard(Spades, King),
		NewCard(Hearts, Queen), NewCard(Clubs, Two),
	}
	h := HandFromCards(cards)
	if h.Len() != 4 {
		t.Errorf("hand length = %d, want 4", h.Len())
	}
	for _, c := range cards {
		if !h.Has(c) {
			t.Errorf("hand should contain %v", c)
		}
	}
}

func TestHandSuitOperations(t *testing.T) {
	cards := []Card{
		NewCard(Spades, Ace), NewCard(Spades, King), NewCard(Spades, Queen),
		NewCard(Hearts, Jack), NewCard(Hearts, Ten),
		NewCard(Diamonds, Five),
	}
	h := HandFromCards(cards)

	if h.SuitLen(Spades) != 3 {
		t.Errorf("spades length = %d, want 3", h.SuitLen(Spades))
	}
	if h.SuitLen(Hearts) != 2 {
		t.Errorf("hearts length = %d, want 2", h.SuitLen(Hearts))
	}
	if h.SuitLen(Diamonds) != 1 {
		t.Errorf("diamonds length = %d, want 1", h.SuitLen(Diamonds))
	}
	if h.SuitLen(Clubs) != 0 {
		t.Errorf("clubs length = %d, want 0", h.SuitLen(Clubs))
	}
	if !h.HasSuit(Spades) {
		t.Error("hand should have spades")
	}
	if h.HasSuit(Clubs) {
		t.Error("hand should not have clubs")
	}
}

func TestHandHCP(t *testing.T) {
	cards := []Card{
		NewCard(Spades, Ace),   // 4
		NewCard(Hearts, King),  // 3
		NewCard(Diamonds, Two), // 0
	}
	h := HandFromCards(cards)
	if h.HCP() != 7 {
		t.Errorf("HCP = %d, want 7", h.HCP())
	}
}

func TestHandCards(t *testing.T) {
	cards := []Card{
		NewCard(Clubs, Two),
		NewCard(Spades, Ace),
		NewCard(Hearts, Queen),
		NewCard(Diamonds, Five),
	}
	h := HandFromCards(cards)
	got := h.Cards()

	if len(got) != len(cards) {
		t.Fatalf("Cards() length = %d, want %d", len(got), len(cards))
	}

	expected := []Card{
		NewCard(Spades, Ace),
		NewCard(Hearts, Queen),
		NewCard(Diamonds, Five),
		NewCard(Clubs, Two),
	}
	for i, c := range got {
		if c != expected[i] {
			t.Errorf("Cards()[%d] = %v, want %v", i, c, expected[i])
		}
	}
}

func TestHandPBNRoundTrip(t *testing.T) {
	tests := []struct {
		pbn  string
		hcp  int
		size int
	}{
		{"AKQ2.KJ3.T98.Q65", 15, 13},
		{"AKQJT98765432...", 10, 13},
		{"..AKQJT98765432.", 10, 13},
		{"...", 0, 0},
	}

	for _, tt := range tests {
		h, ok := PBNToHand(tt.pbn)
		if !ok {
			t.Errorf("PBNToHand(%q) failed", tt.pbn)
			continue
		}
		if h.Len() != tt.size {
			t.Errorf("PBNToHand(%q).Len() = %d, want %d", tt.pbn, h.Len(), tt.size)
		}
		if h.HCP() != tt.hcp {
			t.Errorf("PBNToHand(%q).HCP() = %d, want %d", tt.pbn, h.HCP(), tt.hcp)
		}
		got := h.ToPBN()
		if got != tt.pbn {
			t.Errorf("round-trip PBN: got %q, want %q", got, tt.pbn)
		}
	}
}

func TestPBNToHandInvalid(t *testing.T) {
	invalids := []string{
		"AKQ",               // only 1 part
		"AKQ.KJ3.T98",       // only 3 parts
		"AKQ.KJ3.T98.Q65.2", // 5 parts
		"AKX.KJ3.T98.Q65",   // invalid rank 'X'
		"AAK.KJ3.T98.Q65",   // duplicate card
	}
	for _, s := range invalids {
		if _, ok := PBNToHand(s); ok {
			t.Errorf("PBNToHand(%q) should have failed", s)
		}
	}
}
