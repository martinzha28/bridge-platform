package game

import "fmt"

type Suit uint8
type Rank uint8

// Card is a single playing card encoded as suit*13 + rank. 0-51 from club 2 to spade A
type Card uint8

const NumCards = 52

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

const NumSuits = 4

const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

const NumRanks = 13

var suitNames = [NumSuits]string{"Clubs", "Diamonds", "Hearts", "Spades"}
var suitLetters = [NumSuits]byte{'C', 'D', 'H', 'S'}
var rankLetters = [NumRanks]byte{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A'}

func (s Suit) String() string {
	if s < NumSuits {
		return suitNames[s]
	}
	return fmt.Sprintf("Suit(%d)", s)
}

func (s Suit) Letter() byte {
	return suitLetters[s]
}

func (r Rank) String() string {
	if r < NumRanks {
		return string(rankLetters[r])
	}
	return fmt.Sprintf("Rank(%d)", r)
}

func (r Rank) Letter() byte {
	return rankLetters[r]
}

func (c Card) String() string {
	return string([]byte{c.Suit().Letter(), c.Rank().Letter()})
}

func ParseSuit(b byte) (Suit, bool) {
	switch b {
	case 'C', 'c':
		return Clubs, true
	case 'D', 'd':
		return Diamonds, true
	case 'H', 'h':
		return Hearts, true
	case 'S', 's':
		return Spades, true
	default:
		return 0, false
	}
}

func ParseRank(b byte) (Rank, bool) {
	switch b {
	case '2':
		return Two, true
	case '3':
		return Three, true
	case '4':
		return Four, true
	case '5':
		return Five, true
	case '6':
		return Six, true
	case '7':
		return Seven, true
	case '8':
		return Eight, true
	case '9':
		return Nine, true
	case 'T', 't':
		return Ten, true
	case 'J', 'j':
		return Jack, true
	case 'Q', 'q':
		return Queen, true
	case 'K', 'k':
		return King, true
	case 'A', 'a':
		return Ace, true
	default:
		return 0, false
	}
}

func ParseCard(s string) (Card, bool) {
	if len(s) != 2 {
		return 0, false
	}
	suit, ok := ParseSuit(s[0])
	if !ok {
		return 0, false
	}
	rank, ok := ParseRank(s[1])
	if !ok {
		return 0, false
	}
	return NewCard(suit, rank), true
}

// TODO: Implement better hand evaluation
func (r Rank) HCP() int {
	switch r {
	case Jack:
		return 1
	case Queen:
		return 2
	case King:
		return 3
	case Ace:
		return 4
	default:
		return 0
	}
}

func NewCard(suit Suit, rank Rank) Card {
	return Card(uint8(suit)*NumRanks + uint8(rank))
}

func (c Card) Suit() Suit { return Suit(uint8(c) / NumRanks) }
func (c Card) Rank() Rank { return Rank(uint8(c) % NumRanks) }
func (c Card) HCP() int   { return c.Rank().HCP() }
