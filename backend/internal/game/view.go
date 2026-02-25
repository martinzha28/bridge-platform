package game

type PlayedCardView struct {
	Seat string `json:"seat"`
	Card string `json:"card"`
}

type ResultView struct {
	Contract  string `json:"contract,omitempty"`
	Declarer  string `json:"declarer,omitempty"`
	TricksNS  int    `json:"tricksNS"`
	TricksEW  int    `json:"tricksEW"`
	Score     int    `json:"score"`
	PassedOut bool   `json:"passedOut"`
}

type PlayerView struct {
	Phase         string `json:"phase"`
	BoardNumber   int    `json:"boardNumber"`
	Dealer        string `json:"dealer"`
	Vulnerability string `json:"vulnerability"`

	Seat string   `json:"seat"`
	Hand []string `json:"hand"`

	Calls    []string `json:"calls"`
	Contract string   `json:"contract,omitempty"`

	Dummy     string   `json:"dummy,omitempty"`
	DummyHand []string `json:"dummyHand,omitempty"`

	CurrentTrick []PlayedCardView `json:"currentTrick"`
	TricksNS     int              `json:"tricksNS"`
	TricksEW     int              `json:"tricksEW"`

	Turn string `json:"turn"`

	Result *ResultView `json:"result,omitempty"`

	LegalCalls []string `json:"legalCalls,omitempty"`
	LegalCards []string `json:"legalCards,omitempty"`
}

// ViewFor returns a filtered view of the game from dir's perspective.
func (g *Game) ViewFor(dir Direction) PlayerView {
	v := PlayerView{
		Phase:         g.Phase.String(),
		BoardNumber:   g.Board.Number,
		Dealer:        g.Board.Dealer.String(),
		Vulnerability: g.Board.Vulnerability.String(),
		Seat:          dir.String(),
	}

	switch g.Phase {
	case PhasePlay:
		v.Hand = cardsToStrings(g.Play.RemainingHands[dir].Cards())
	case PhaseAuction, PhaseComplete:
		v.Hand = cardsToStrings(g.Board.Hands[dir].Cards())
	}

	if g.Auction != nil {
		v.Calls = make([]string, len(g.Auction.Calls))
		for i, c := range g.Auction.Calls {
			v.Calls[i] = c.String()
		}
		if contract, ok := g.Auction.Contract(); ok {
			v.Contract = contract.String()
		}
	}

	if g.Play != nil {
		v.Dummy = g.Play.Dummy.String()

		if g.DummyVisible() {
			v.DummyHand = cardsToStrings(g.Play.RemainingHands[g.Play.Dummy].Cards())
		}

		trick := g.Play.currentTrick()
		trickDir := trick.Leader
		for i := range trick.PlayedSoFar {
			v.CurrentTrick = append(v.CurrentTrick, PlayedCardView{
				Seat: trickDir.String(),
				Card: trick.Cards[i].String(),
			})
			trickDir = trickDir.Next()
		}

		v.TricksNS = g.Play.TricksNS
		v.TricksEW = g.Play.TricksEW
	}

	if turn, ok := g.Turn(); ok {
		v.Turn = turn.String()
	}

	if g.Result != nil {
		rv := &ResultView{
			TricksNS:  g.Result.TricksNS,
			TricksEW:  g.Result.TricksEW,
			Score:     g.Result.Score,
			PassedOut: g.Result.PassedOut(),
		}
		if g.Result.Contract != nil {
			rv.Contract = g.Result.Contract.String()
			rv.Declarer = g.Result.Contract.Declarer.String()
		}
		v.Result = rv
	}

	if lc := g.LegalCalls(dir); lc != nil {
		v.LegalCalls = make([]string, len(lc))
		for i, c := range lc {
			v.LegalCalls[i] = c.String()
		}
	}
	if lp := g.LegalCards(dir); lp != nil {
		v.LegalCards = cardsToStrings(lp)
	}

	return v
}

func cardsToStrings(cards []Card) []string {
	s := make([]string, len(cards))
	for i, c := range cards {
		s[i] = c.String()
	}
	return s
}
