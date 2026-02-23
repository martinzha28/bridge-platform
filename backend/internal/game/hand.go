package game

import (
	"math/bits"
	"strings"
)

// Hand is a bitmask of cards. Bit i corresponds to Card(i).
// Bits 0–12 = Clubs, 13–25 = Diamonds, 26–38 = Hearts, 39–51 = Spades.
type Hand uint64

const (
	suitMaskWidth = NumRanks
	allRankBits   = (1 << suitMaskWidth) - 1

	ClubsMask    Hand = Hand(allRankBits) << (uint(Clubs) * suitMaskWidth)
	DiamondsMask Hand = Hand(allRankBits) << (uint(Diamonds) * suitMaskWidth)
	HeartsMask   Hand = Hand(allRankBits) << (uint(Hearts) * suitMaskWidth)
	SpadesMask   Hand = Hand(allRankBits) << (uint(Spades) * suitMaskWidth)
)

var suitMasks = [NumSuits]Hand{ClubsMask, DiamondsMask, HeartsMask, SpadesMask}

func HandFromCards(cards []Card) Hand {
	var h Hand

	for _, c := range cards {
		h |= 1 << c
	}
	return h
}

func (h Hand) Has(c Card) bool {
	return h&(1<<c) != 0
}

func (h Hand) Add(c Card) Hand {
	return h | (1 << c)
}

func (h Hand) Remove(c Card) Hand {
	return h &^ (1 << c)
}

func (h Hand) Len() int {
	return bits.OnesCount64(uint64(h))
}

func (h Hand) SuitHolding(s Suit) Hand {
	return h & suitMasks[s]
}

func (h Hand) SuitLen(s Suit) int {
	return bits.OnesCount64(uint64(h.SuitHolding(s)))
}

func (h Hand) HasSuit(s Suit) bool {
	return h.SuitHolding(s) != 0
}

func (h Hand) HCP() int {
	total := 0

	for c := Card(0); c < NumCards; c++ {
		if h.Has(c) {
			total += c.HCP()
		}
	}
	return total
}

func (h Hand) Cards() []Card {
	cards := make([]Card, 0, h.Len())

	for s := Suit(NumSuits - 1); ; s-- {
		holding := uint64(h.SuitHolding(s)) >> (uint(s) * suitMaskWidth)

		for r := Rank(NumRanks - 1); ; r-- {
			if holding&(1<<r) != 0 {
				cards = append(cards, NewCard(s, r))
			}
			if r == 0 {
				break
			}
		}
		if s == 0 {
			break
		}
	}
	return cards
}

func (h Hand) ToPBN() string {
	suits := [NumSuits]Suit{Spades, Hearts, Diamonds, Clubs}
	var buf [NumCards + 3]byte // 52 max rank chars + 3 dots
	n := 0

	for i, s := range suits {
		if i > 0 {
			buf[n] = '.'
			n++
		}

		shift := uint(s) * suitMaskWidth
		holding := uint64(h) >> shift & allRankBits

		for r := Rank(NumRanks - 1); ; r-- {
			if holding&(1<<r) != 0 {
				buf[n] = r.Letter()
				n++
			}
			if r == 0 {
				break
			}
		}
	}
	return string(buf[:n])
}

func PBNToHand(s string) (Hand, bool) {
	parts := strings.Split(s, ".")

	if len(parts) != NumSuits {
		return 0, false
	}

	suits := [NumSuits]Suit{Spades, Hearts, Diamonds, Clubs}
	var h Hand

	for i, part := range parts {
		for j := 0; j < len(part); j++ {
			rank, ok := ParseRank(part[j])
			if !ok {
				return 0, false
			}

			c := NewCard(suits[i], rank)

			if h.Has(c) {
				return 0, false
			}
			h = h.Add(c)
		}
	}
	return h, true
}
