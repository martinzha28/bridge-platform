package game

import (
	"fmt"
	"strings"
)

const NumTricks = 13

type Trick struct {
	Leader      Direction
	Cards       [NumDirections]Card
	PlayedSoFar int
}

type Play struct {
	Contract       Contract
	RemainingHands [NumDirections]Hand
	Declarer       Direction
	Dummy          Direction

	Tricks   []Trick
	TricksNS int
	TricksEW int

	turn     Direction
	finished bool
}

func NewPlay(contract Contract, hands [NumDirections]Hand) *Play {
	leader := contract.Declarer.Next()

	return &Play{
		Contract:       contract,
		RemainingHands: hands,
		Declarer:       contract.Declarer,
		Dummy:          contract.Declarer.Partner(),
		Tricks:         []Trick{{Leader: leader}},
		turn:           leader,
	}
}

func (p *Play) Turn() Direction  { return p.turn }
func (p *Play) IsFinished() bool { return p.finished }

func (p *Play) currentTrick() *Trick {
	return &p.Tricks[len(p.Tricks)-1]
}

// lastCompletedTrick returns the trick that was most recently finished
// and is still sitting on the table — i.e. the next trick hasn't been
// led yet, or the hand is over. Returns nil otherwise.
func (p *Play) lastCompletedTrick() *Trick {
	n := len(p.Tricks)
	if n == 0 {
		return nil
	}
	if p.finished {
		return &p.Tricks[n-1]
	}
	if n >= 2 && p.Tricks[n-1].PlayedSoFar == 0 {
		return &p.Tricks[n-2]
	}
	return nil
}

func (p *Play) PlayCard(actor Direction, card Card) error {
	if p.finished {
		return fmt.Errorf("play is already finished")
	}

	if err := p.validateActor(actor); err != nil {
		return err
	}

	hand := &p.RemainingHands[p.turn]
	if !hand.Has(card) {
		return fmt.Errorf("%v does not have %v", p.turn, card)
	}

	if err := p.validateFollowSuit(card); err != nil {
		return err
	}

	trick := p.currentTrick()
	trick.Cards[trick.PlayedSoFar] = card
	trick.PlayedSoFar++
	*hand = hand.Remove(card)

	if trick.PlayedSoFar == NumDirections {
		p.resolveTrick()
	} else {
		p.turn = p.turn.Next()
	}

	return nil
}

func (p *Play) validateActor(actor Direction) error {
	if p.turn == p.Dummy {
		if actor != p.Declarer {
			return fmt.Errorf("only declarer (%v) can play dummy's cards", p.Declarer)
		}
		return nil
	}
	if actor != p.turn {
		return fmt.Errorf("not %v's turn (expected %v)", actor, p.turn)
	}
	return nil
}

// validateFollowSuit ensures the player follows suit if they can.
func (p *Play) validateFollowSuit(card Card) error {
	trick := p.currentTrick()

	if trick.PlayedSoFar == 0 {
		return nil
	}

	ledSuit := trick.Cards[0].Suit()

	if card.Suit() != ledSuit && p.RemainingHands[p.turn].HasSuit(ledSuit) {
		return fmt.Errorf("must follow suit (%v)", ledSuit)
	}
	return nil
}

func (p *Play) resolveTrick() {
	trick := p.currentTrick()
	winner := p.trickWinner(trick)

	if winner == North || winner == South {
		p.TricksNS++
	} else {
		p.TricksEW++
	}

	if p.RemainingHands[winner].Len() == 0 {
		p.finished = true
		return
	}

	p.turn = winner
	p.Tricks = append(p.Tricks, Trick{Leader: winner})
}

func (p *Play) trickWinner(trick *Trick) Direction {
	ledSuit := trick.Cards[0].Suit()
	trump := p.strainToTrump()

	winnerIdx := 0
	winnerCard := trick.Cards[0]

	for i := 1; i < NumDirections; i++ {
		candidate := trick.Cards[i]

		if beats(candidate, winnerCard, ledSuit, trump) {
			winnerIdx = i
			winnerCard = candidate
		}
	}

	dir := trick.Leader

	for range winnerIdx {
		dir = dir.Next()
	}
	return dir
}

func (p *Play) strainToTrump() *Suit {
	switch p.Contract.Strain {
	case ClubStrain:
		s := Clubs
		return &s

	case DiamondStrain:
		s := Diamonds
		return &s

	case HeartStrain:
		s := Hearts
		return &s

	case SpadeStrain:
		s := Spades
		return &s

	// return nil for NoTrump
	default:
		return nil
	}
}

func beats(a, b Card, ledSuit Suit, trump *Suit) bool {
	// checks if cards are trump, which is optional
	aIsTrump := trump != nil && a.Suit() == *trump
	bIsTrump := trump != nil && b.Suit() == *trump

	switch {
	case aIsTrump && !bIsTrump:
		return true

	case !aIsTrump && bIsTrump:
		return false

	case aIsTrump && bIsTrump:
		return a.Rank() > b.Rank()

	default:
		if a.Suit() == ledSuit && b.Suit() != ledSuit {
			return true
		}

		if a.Suit() != ledSuit && b.Suit() == ledSuit {
			return false
		}

		if a.Suit() == ledSuit && b.Suit() == ledSuit {
			return a.Rank() > b.Rank()
		}

		return false
	}
}

func (p *Play) ToNotation() string {
	var parts []string
	for _, trick := range p.Tricks {
		var cards []string
		for i := range trick.PlayedSoFar {
			cards = append(cards, trick.Cards[i].String())
		}
		parts = append(parts, strings.Join(cards, " "))
	}
	return strings.Join(parts, " | ")
}
