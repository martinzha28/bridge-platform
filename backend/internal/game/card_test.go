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

func TestParseCardInvalid(t *testing.T) {
	invalids := []string{"", "A", "XY", "ZZ", "S1", "XX"}
	for _, s := range invalids {
		if _, ok := ParseCard(s); ok {
			t.Errorf("ParseCard(%q) should have failed", s)
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
