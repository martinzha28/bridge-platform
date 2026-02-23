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

func TestPBNToHandVoids(t *testing.T) {
	tests := []struct {
		pbn       string
		voidSuits []Suit
	}{
		{"AKQ.AKQ.AKQ.AKQJ", nil},
		{".AKQJT98765432..", []Suit{Spades, Diamonds, Clubs}},
		{"AKQJT98765432...", []Suit{Hearts, Diamonds, Clubs}},
		{"...AKQJT98765432", []Suit{Spades, Hearts, Diamonds}},
		{"..AKQJT98765432.", []Suit{Spades, Hearts, Clubs}},
	}

	for _, tt := range tests {
		h, ok := PBNToHand(tt.pbn)
		if !ok {
			t.Errorf("PBNToHand(%q) failed", tt.pbn)
			continue
		}
		if h.Len() != 13 {
			t.Errorf("PBNToHand(%q).Len() = %d, want 13", tt.pbn, h.Len())
		}
		for _, suit := range tt.voidSuits {
			if h.HasSuit(suit) {
				t.Errorf("PBNToHand(%q) should be void in %v", tt.pbn, suit)
			}
		}
	}
}

func TestPBNToHandSpecificCards(t *testing.T) {
	h, ok := PBNToHand("AK.QJ.T9.8765432")
	if !ok {
		t.Fatal("PBNToHand failed")
	}

	hasCards := []Card{
		NewCard(Spades, Ace), NewCard(Spades, King),
		NewCard(Hearts, Queen), NewCard(Hearts, Jack),
		NewCard(Diamonds, Ten), NewCard(Diamonds, Nine),
		NewCard(Clubs, Eight), NewCard(Clubs, Seven),
	}
	for _, c := range hasCards {
		if !h.Has(c) {
			t.Errorf("hand should contain %v", c)
		}
	}

	missingCards := []Card{
		NewCard(Spades, Queen),
		NewCard(Hearts, Ace),
		NewCard(Diamonds, Ace),
		NewCard(Clubs, Ace),
	}
	for _, c := range missingCards {
		if h.Has(c) {
			t.Errorf("hand should not contain %v", c)
		}
	}
}

func TestPBNToHandInvalid(t *testing.T) {
	invalids := []struct {
		pbn    string
		reason string
	}{
		{"AKQ", "only 1 part"},
		{"AKQ.KJ3.T98", "only 3 parts"},
		{"AKQ.KJ3.T98.Q65.2", "5 parts"},
		{"AKX.KJ3.T98.Q65", "invalid rank X"},
		{"AAK.KJ3.T98.Q65", "duplicate card within suit"},
		{"", "empty string"},
		{"...", "all voids (valid — 0 cards)"},
	}
	for _, tt := range invalids {
		h, ok := PBNToHand(tt.pbn)
		// "..." is actually valid (0 cards), the others should fail
		if tt.pbn == "..." {
			if !ok {
				t.Errorf("PBNToHand(%q) should succeed: %s", tt.pbn, tt.reason)
			}
			if h.Len() != 0 {
				t.Errorf("PBNToHand(%q).Len() = %d, want 0", tt.pbn, h.Len())
			}
			continue
		}
		if ok {
			t.Errorf("PBNToHand(%q) should have failed (%s)", tt.pbn, tt.reason)
		}
	}
}
