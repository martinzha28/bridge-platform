package game

import "testing"

func TestNewCard(t *testing.T) {
	tests := []struct {
		suit Suit
		rank Rank
		want Card
		str  string
	}{
		{Clubs, Two, Card(0), "C2"},
		{Clubs, Ace, Card(12), "CA"},
		{Diamonds, Two, Card(13), "D2"},
		{Hearts, King, Card(37), "HK"},
		{Spades, Ace, Card(51), "SA"},
	}
	for _, tt := range tests {
		c := NewCard(tt.suit, tt.rank)
		if c != tt.want {
			t.Errorf("NewCard(%v, %v) = %d, want %d", tt.suit, tt.rank, c, tt.want)
		}
		if c.Suit() != tt.suit {
			t.Errorf("Card(%d).Suit() = %v, want %v", c, c.Suit(), tt.suit)
		}
		if c.Rank() != tt.rank {
			t.Errorf("Card(%d).Rank() = %v, want %v", c, c.Rank(), tt.rank)
		}
		if c.String() != tt.str {
			t.Errorf("Card(%d).String() = %q, want %q", c, c.String(), tt.str)
		}
	}
}

func TestParseCard(t *testing.T) {
	for suit := Suit(0); suit < NumSuits; suit++ {
		for rank := Rank(0); rank < NumRanks; rank++ {
			c := NewCard(suit, rank)
			parsed, ok := ParseCard(c.String())
			if !ok {
				t.Errorf("ParseCard(%q) failed", c.String())
				continue
			}
			if parsed != c {
				t.Errorf("ParseCard(%q) = %d, want %d", c.String(), parsed, c)
			}
		}
	}
}

func TestParseCardCaseInsensitive(t *testing.T) {
	tests := []struct {
		input string
		want  Card
	}{
		{"sa", NewCard(Spades, Ace)},
		{"hk", NewCard(Hearts, King)},
		{"dq", NewCard(Diamonds, Queen)},
		{"cj", NewCard(Clubs, Jack)},
		{"St", NewCard(Spades, Ten)},
	}
	for _, tt := range tests {
		got, ok := ParseCard(tt.input)
		if !ok {
			t.Errorf("ParseCard(%q) failed", tt.input)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseCard(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseCardInvalid(t *testing.T) {
	invalids := []string{"", "A", "XY", "ZZ", "S1", "XX", "ABC", "S"}
	for _, s := range invalids {
		if _, ok := ParseCard(s); ok {
			t.Errorf("ParseCard(%q) should have failed", s)
		}
	}
}

func TestParseSuit(t *testing.T) {
	tests := []struct {
		input byte
		want  Suit
		ok    bool
	}{
		{'C', Clubs, true},
		{'D', Diamonds, true},
		{'H', Hearts, true},
		{'S', Spades, true},
		{'c', Clubs, true},
		{'d', Diamonds, true},
		{'h', Hearts, true},
		{'s', Spades, true},
		{'X', 0, false},
		{'1', 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseSuit(tt.input)
		if ok != tt.ok {
			t.Errorf("ParseSuit(%q) ok = %v, want %v", tt.input, ok, tt.ok)
			continue
		}
		if ok && got != tt.want {
			t.Errorf("ParseSuit(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseRank(t *testing.T) {
	tests := []struct {
		input byte
		want  Rank
		ok    bool
	}{
		{'2', Two, true},
		{'9', Nine, true},
		{'T', Ten, true},
		{'t', Ten, true},
		{'J', Jack, true},
		{'j', Jack, true},
		{'Q', Queen, true},
		{'K', King, true},
		{'A', Ace, true},
		{'a', Ace, true},
		{'1', 0, false},
		{'0', 0, false},
		{'X', 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseRank(tt.input)
		if ok != tt.ok {
			t.Errorf("ParseRank(%q) ok = %v, want %v", tt.input, ok, tt.ok)
			continue
		}
		if ok && got != tt.want {
			t.Errorf("ParseRank(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestHCP(t *testing.T) {
	tests := []struct {
		rank Rank
		want int
	}{
		{Two, 0}, {Nine, 0}, {Ten, 0},
		{Jack, 1}, {Queen, 2}, {King, 3}, {Ace, 4},
	}
	for _, tt := range tests {
		if got := tt.rank.HCP(); got != tt.want {
			t.Errorf("Rank(%v).HCP() = %d, want %d", tt.rank, got, tt.want)
		}
	}
}
