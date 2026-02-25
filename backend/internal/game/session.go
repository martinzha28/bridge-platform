package game

import "fmt"

type Phase uint8

const (
	PhaseDeal Phase = iota
	PhaseAuction
	PhasePlay
	PhaseComplete
)

var phaseNames = [4]string{"Deal", "Auction", "Play", "Complete"}

func (p Phase) String() string {
	if p <= PhaseComplete {
		return phaseNames[p]
	}
	return fmt.Sprintf("Phase(%d)", p)
}

type GameResult struct {
	Contract *Contract
	TricksNS int
	TricksEW int
	Score    int
}

func (r *GameResult) PassedOut() bool {
	return r.Contract == nil
}

type Game struct {
	Board   Board
	Phase   Phase
	Auction *Auction
	Play    *Play
	Result  *GameResult
}

func NewGame(boardNum int) *Game {
	return &Game{
		Board: Board{
			Number:        boardNum,
			Dealer:        DealerForBoard(boardNum),
			Vulnerability: VulnerabilityForBoard(boardNum),
		},
		Phase: PhaseDeal,
	}
}

func (g *Game) Deal(seed int64) error {
	if g.Phase != PhaseDeal {
		return fmt.Errorf("cannot deal: game is in %v phase", g.Phase)
	}
	g.Board.Seed = seed
	g.Board.Hands = DealFromSeed(seed)
	g.Phase = PhaseAuction
	g.Auction = NewAuction(g.Board.Dealer)
	return nil
}

// TODO: Implement hand constraints or preset deals
func (g *Game) SetHands(hands [NumDirections]Hand) error {
	if g.Phase != PhaseDeal {
		return fmt.Errorf("cannot set hands: game is in %v phase", g.Phase)
	}
	if err := validateHands(hands); err != nil {
		return err
	}
	g.Board.Hands = hands
	g.Phase = PhaseAuction
	g.Auction = NewAuction(g.Board.Dealer)
	return nil
}

func validateHands(hands [NumDirections]Hand) error {
	var all Hand
	for i, h := range hands {
		if h.Len() != NumRanks {
			return fmt.Errorf("%v has %d cards, want %d", Direction(i), h.Len(), NumRanks)
		}
		if all&h != 0 {
			return fmt.Errorf("duplicate cards in deal")
		}
		all |= h
	}
	return nil
}

func (g *Game) Turn() (Direction, bool) {
	switch g.Phase {
	case PhaseAuction:
		return g.Auction.Turn(), true
	case PhasePlay:
		return g.Play.Turn(), true
	default:
		return 0, false
	}
}

func (g *Game) MakeCall(dir Direction, call Call) error {
	if g.Phase != PhaseAuction {
		return fmt.Errorf("cannot make call: game is in %v phase", g.Phase)
	}

	if err := g.Auction.MakeCall(dir, call); err != nil {
		return err
	}
	if g.Auction.IsFinished() {
		g.advanceFromAuction()
	}
	return nil
}

func (g *Game) PlayCard(actor Direction, card Card) error {
	if g.Phase != PhasePlay {
		return fmt.Errorf("cannot play card: game is in %v phase", g.Phase)
	}
	if err := g.Play.PlayCard(actor, card); err != nil {
		return err
	}
	if g.Play.IsFinished() {
		g.advanceFromPlay()
	}
	return nil
}

func (g *Game) advanceFromAuction() {
	if g.Auction.PassedOut() {
		g.Phase = PhaseComplete
		g.Result = &GameResult{}
		return
	}

	contract, _ := g.Auction.Contract()
	g.Play = NewPlay(contract, g.Board.Hands)
	g.Phase = PhasePlay
}

func (g *Game) advanceFromPlay() {
	contract := g.Play.Contract
	vul := g.Board.Vulnerability.IsVulnerable(contract.Declarer)

	tricksTaken := g.Play.TricksNS
	if contract.Declarer == East || contract.Declarer == West {
		tricksTaken = g.Play.TricksEW
	}

	g.Phase = PhaseComplete
	g.Result = &GameResult{
		Contract: &contract,
		TricksNS: g.Play.TricksNS,
		TricksEW: g.Play.TricksEW,
		Score:    Score(contract, vul, tricksTaken),
	}
}

func (g *Game) DummyVisible() bool {
	return g.Play != nil && len(g.Play.Tricks) > 0 && g.Play.Tricks[0].PlayedSoFar > 0
}

func (g *Game) LegalCalls(dir Direction) []Call {
	if g.Phase != PhaseAuction || g.Auction.Turn() != dir {
		return nil
	}

	var calls []Call
	calls = append(calls, passCall)

	if g.Auction.validateCall(doubleCall) == nil {
		calls = append(calls, doubleCall)
	}
	if g.Auction.validateCall(redoubleCall) == nil {
		calls = append(calls, redoubleCall)
	}

	for level := uint8(1); level <= 7; level++ {
		for strain := Strain(0); strain < NumStrains; strain++ {
			call := BidCall(level, strain)
			if g.Auction.validateCall(call) == nil {
				calls = append(calls, call)
			}
		}
	}
	return calls
}

func (g *Game) LegalCards(dir Direction) []Card {
	if g.Phase != PhasePlay || g.Play.IsFinished() {
		return nil
	}

	whoseTurn := g.Play.Turn()
	if whoseTurn == g.Play.Dummy {
		if dir != g.Play.Declarer {
			return nil
		}
	} else if dir != whoseTurn {
		return nil
	}

	hand := g.Play.RemainingHands[whoseTurn]
	trick := g.Play.currentTrick()

	if trick.PlayedSoFar == 0 {
		return hand.Cards()
	}

	ledSuit := trick.Cards[0].Suit()
	if hand.HasSuit(ledSuit) {
		return hand.SuitHolding(ledSuit).Cards()
	}
	return hand.Cards()
}
