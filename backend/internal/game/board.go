package game

import (
	"math/rand/v2"
	"strings"
)

type Board struct {
	Number        int
	Seed          int64
	Dealer        Direction
	Vulnerability Vulnerability
	Hands         [NumDirections]Hand
}

func DealerForBoard(boardNum int) Direction {
	// 1-16 -> N, E, S, W, N, ...
	return Direction((boardNum - 1) % NumDirections)
}

func VulnerabilityForBoard(boardNum int) Vulnerability {
	//  1: None   5: NS     9: EW   13: Both
	//  2: NS     6: EW    10: Both  14: None
	//  3: EW     7: Both  11: None  15: NS
	//  4: Both   8: None  12: NS    16: EW
	cycle := [16]Vulnerability{
		VulNone, VulNS, VulEW, VulBoth,
		VulNS, VulEW, VulBoth, VulNone,
		VulEW, VulBoth, VulNone, VulNS,
		VulBoth, VulNone, VulNS, VulEW,
	}
	return cycle[(boardNum-1)%16]
}

// Deal from a seed using Fisher-Yates shuffle
func DealFromSeed(seed int64) [NumDirections]Hand {
	var deck [NumCards]Card
	for i := range deck {
		deck[i] = Card(i)
	}

	src := rand.NewChaCha8([32]byte{
		byte(seed), byte(seed >> 8), byte(seed >> 16), byte(seed >> 24),
		byte(seed >> 32), byte(seed >> 40), byte(seed >> 48), byte(seed >> 56),
	})
	rng := rand.New(src)

	for i := NumCards - 1; i > 0; i-- {
		j := rng.IntN(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	var hands [NumDirections]Hand
	for i, c := range deck {
		dir := Direction(i / NumRanks)
		hands[dir] = hands[dir].Add(c)
	}
	return hands
}

func NewBoard(boardNum int, seed int64) Board {
	return Board{
		Number:        boardNum,
		Seed:          seed,
		Dealer:        DealerForBoard(boardNum),
		Vulnerability: VulnerabilityForBoard(boardNum),
		Hands:         DealFromSeed(seed),
	}
}

func (b Board) ToPBN() string {
	result := string(b.Dealer.Letter()) + ":"

	for i := range NumDirections {
		dir := Direction((int(b.Dealer) + i) % NumDirections)
		if i > 0 {
			result += " "
		}
		result += b.Hands[dir].ToPBN()
	}
	return result
}

func PBNToBoard(s string) (Direction, [NumDirections]Hand, bool) {
	if len(s) < 3 || s[1] != ':' {
		return 0, [NumDirections]Hand{}, false
	}

	dealer, ok := ParseDirection(s[0])
	if !ok {
		return 0, [NumDirections]Hand{}, false
	}

	rest := s[2:]
	parts := strings.Fields(rest)
	if len(parts) != NumDirections {
		return 0, [NumDirections]Hand{}, false
	}

	var hands [NumDirections]Hand
	totalCards := 0
	for i, part := range parts {
		dir := Direction((int(dealer) + i) % NumDirections)
		h, ok := PBNToHand(part)
		if !ok {
			return 0, [NumDirections]Hand{}, false
		}
		hands[dir] = h
		totalCards += h.Len()
	}

	if totalCards != NumCards {
		return 0, [NumDirections]Hand{}, false
	}

	return dealer, hands, true
}
