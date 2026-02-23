package game

import "testing"

func setupHeartsPlay() (*Play, [NumDirections]Hand) {
	hands := [NumDirections]Hand{
		// AKQ.AK.AKQ.AKQ.JT
		North: HandFromCards([]Card{ // dummy
			NewCard(Spades, Ace), NewCard(Spades, King), NewCard(Spades, Queen),
			NewCard(Hearts, Ace), NewCard(Hearts, King),
			NewCard(Diamonds, Ace), NewCard(Diamonds, King), NewCard(Diamonds, Queen),
			NewCard(Clubs, Ace), NewCard(Clubs, King), NewCard(Clubs, Queen),
			NewCard(Clubs, Jack), NewCard(Clubs, Ten),
		}),
		// JT9.QJ.JT9.98765
		East: HandFromCards([]Card{
			NewCard(Spades, Jack), NewCard(Spades, Ten), NewCard(Spades, Nine),
			NewCard(Hearts, Queen), NewCard(Hearts, Jack),
			NewCard(Diamonds, Jack), NewCard(Diamonds, Ten), NewCard(Diamonds, Nine),
			NewCard(Clubs, Nine), NewCard(Clubs, Eight), NewCard(Clubs, Seven),
			NewCard(Clubs, Six), NewCard(Clubs, Five),
		}),
		// 876.T98.8765432.
		South: HandFromCards([]Card{ // declarer
			NewCard(Spades, Eight), NewCard(Spades, Seven), NewCard(Spades, Six),
			NewCard(Hearts, Ten), NewCard(Hearts, Nine), NewCard(Hearts, Eight),
			NewCard(Diamonds, Eight), NewCard(Diamonds, Seven), NewCard(Diamonds, Six),
			NewCard(Diamonds, Five), NewCard(Diamonds, Four), NewCard(Diamonds, Three),
			NewCard(Diamonds, Two),
		}),
		// 5432.765432..432
		West: HandFromCards([]Card{ // leader
			NewCard(Spades, Five), NewCard(Spades, Four), NewCard(Spades, Three), NewCard(Spades, Two),
			NewCard(Hearts, Seven), NewCard(Hearts, Six), NewCard(Hearts, Five),
			NewCard(Hearts, Four), NewCard(Hearts, Three), NewCard(Hearts, Two),
			NewCard(Clubs, Four), NewCard(Clubs, Three), NewCard(Clubs, Two),
		}),
	}

	contract := Contract{
		Level:    4,
		Strain:   HeartStrain,
		Declarer: South,
	}

	return NewPlay(contract, hands), hands
}

func TestPlayOpeningLeader(t *testing.T) {
	p, _ := setupHeartsPlay()

	if p.Turn() != West {
		t.Errorf("opening leader = %v, want West", p.Turn())
	}
	if p.Dummy != North {
		t.Errorf("dummy = %v, want North", p.Dummy)
	}
}

func TestPlayOneTrick(t *testing.T) {
	p, _ := setupHeartsPlay()

	// West leads S5
	must(t, p.PlayCard(West, NewCard(Spades, Five)))
	// North (dummy) plays SA but declarer plays for dummy
	must(t, p.PlayCard(South, NewCard(Spades, Ace)))
	// East plays SJ
	must(t, p.PlayCard(East, NewCard(Spades, Jack)))
	// South plays S8
	must(t, p.PlayCard(South, NewCard(Spades, Eight)))

	if p.TricksNS != 1 {
		t.Errorf("TricksNS = %d, want 1", p.TricksNS)
	}
	if p.TricksEW != 0 {
		t.Errorf("TricksEW = %d, want 0", p.TricksEW)
	}

	// North (dummy) won with SA, so dummy leads next
	if p.Turn() != North {
		t.Errorf("next leader = %v, want North (won the trick)", p.Turn())
	}
}

func TestPlayTrumpWins(t *testing.T) {
	p, _ := setupHeartsPlay()

	// West leads S5
	must(t, p.PlayCard(West, NewCard(Spades, Five)))
	// North plays SQ
	must(t, p.PlayCard(South, NewCard(Spades, Queen)))
	// East plays S9
	must(t, p.PlayCard(East, NewCard(Spades, Nine)))
	// South has spades and must follow suit
	must(t, p.PlayCard(South, NewCard(Spades, Six)))

	// North (dummy) won with SQ. Now test trumping.
	// North leads CA
	must(t, p.PlayCard(South, NewCard(Clubs, Ace)))
	// East plays C9
	must(t, p.PlayCard(East, NewCard(Clubs, Nine)))
	// South has no clubs -> can trump with a heart
	must(t, p.PlayCard(South, NewCard(Hearts, Eight)))
	// West plays C2
	must(t, p.PlayCard(West, NewCard(Clubs, Two)))

	// South trumped with H8 -> South wins
	if p.TricksNS != 2 {
		t.Errorf("TricksNS = %d, want 2", p.TricksNS)
	}
}

func TestPlayMustFollowSuit(t *testing.T) {
	p, _ := setupHeartsPlay()

	// West leads S5
	must(t, p.PlayCard(West, NewCard(Spades, Five)))

	// North (dummy) has spades -> trying to play a club fails
	err := p.PlayCard(South, NewCard(Clubs, Ace))
	if err == nil {
		t.Error("expected error: dummy has spades and must follow suit")
	}

	must(t, p.PlayCard(South, NewCard(Spades, Ace)))
}

func TestPlayCardNotInHand(t *testing.T) {
	p, _ := setupHeartsPlay()

	// West trying to play SA fails
	err := p.PlayCard(West, NewCard(Spades, Ace))
	if err == nil {
		t.Error("expected error: West doesn't have SA")
	}
}

func TestPlayWrongTurn(t *testing.T) {
	p, _ := setupHeartsPlay()

	// East tries to play when it's West's turn
	err := p.PlayCard(East, NewCard(Spades, Jack))
	if err == nil {
		t.Error("expected error: it's West's turn, not East's")
	}
}

func TestPlayDummyCantActForSelf(t *testing.T) {
	p, _ := setupHeartsPlay()

	must(t, p.PlayCard(West, NewCard(Spades, Five)))

	// Dummy's turn -> can't play for themselves
	err := p.PlayCard(North, NewCard(Spades, Ace))
	if err == nil {
		t.Error("expected error: dummy can't play for themselves")
	}
}

func TestPlayNoTrumpHighCardWins(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Spades, Ace)}),
		East:  HandFromCards([]Card{NewCard(Spades, King)}),
		South: HandFromCards([]Card{NewCard(Hearts, Ace)}),
		West:  HandFromCards([]Card{NewCard(Spades, Two)}),
	}

	contract := Contract{Level: 1, Strain: NoTrump, Declarer: South}
	p := NewPlay(contract, hands)

	// West leads S2
	must(t, p.PlayCard(West, NewCard(Spades, Two)))
	// North (dummy) plays SA
	must(t, p.PlayCard(South, NewCard(Spades, Ace)))
	// East plays SK
	must(t, p.PlayCard(East, NewCard(Spades, King)))
	// South has no spades, plays HA -> doesn't win in NT
	must(t, p.PlayCard(South, NewCard(Hearts, Ace)))

	if p.TricksNS != 1 {
		t.Errorf("TricksNS = %d, want 1 (SA wins in NT)", p.TricksNS)
	}
	if !p.IsFinished() {
		t.Error("play should be finished (only 1 card each)")
	}
}

func TestPlayOffSuitLoses(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Clubs, Ace)}),
		East:  HandFromCards([]Card{NewCard(Spades, Two)}),
		South: HandFromCards([]Card{NewCard(Diamonds, Ace)}),
		West:  HandFromCards([]Card{NewCard(Spades, Three)}),
	}

	contract := Contract{Level: 1, Strain: NoTrump, Declarer: South}
	p := NewPlay(contract, hands)

	// West leads S3
	must(t, p.PlayCard(West, NewCard(Spades, Three)))
	// North has no spades, plays CA
	must(t, p.PlayCard(South, NewCard(Clubs, Ace)))
	// East follows with S2
	must(t, p.PlayCard(East, NewCard(Spades, Two)))
	// South has no spades, plays DA
	must(t, p.PlayCard(South, NewCard(Diamonds, Ace)))

	// West wins -> S3 beats off-suit aces in NT
	if p.TricksEW != 1 {
		t.Errorf("TricksEW = %d, want 1 (S3 wins, others are off-suit)", p.TricksEW)
	}
}

func TestPlayAfterFinished(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Spades, Ace)}),
		East:  HandFromCards([]Card{NewCard(Spades, King)}),
		South: HandFromCards([]Card{NewCard(Spades, Queen)}),
		West:  HandFromCards([]Card{NewCard(Spades, Jack)}),
	}

	contract := Contract{Level: 1, Strain: NoTrump, Declarer: South}
	p := NewPlay(contract, hands)

	must(t, p.PlayCard(West, NewCard(Spades, Jack)))
	must(t, p.PlayCard(South, NewCard(Spades, Ace)))
	must(t, p.PlayCard(East, NewCard(Spades, King)))
	must(t, p.PlayCard(South, NewCard(Spades, Queen)))

	if !p.IsFinished() {
		t.Fatal("play should be finished")
	}

	err := p.PlayCard(West, NewCard(Spades, Two))
	if err == nil {
		t.Error("expected error: play is already finished")
	}
}

func TestPlayTrumpBeatsHighOffSuit(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Clubs, Two)}),
		East:  HandFromCards([]Card{NewCard(Spades, King)}),
		South: HandFromCards([]Card{NewCard(Spades, Queen)}),
		West:  HandFromCards([]Card{NewCard(Spades, Ace)}),
	}

	contract := Contract{Level: 1, Strain: ClubStrain, Declarer: South}
	p := NewPlay(contract, hands)

	// West leads SA
	must(t, p.PlayCard(West, NewCard(Spades, Ace)))
	// North (dummy) has no spades, trumps with C2
	must(t, p.PlayCard(South, NewCard(Clubs, Two)))
	// East plays SK
	must(t, p.PlayCard(East, NewCard(Spades, King)))
	// South plays SQ
	must(t, p.PlayCard(South, NewCard(Spades, Queen)))

	// North's C2 (trump) beats SA
	if p.TricksNS != 1 {
		t.Errorf("TricksNS = %d, want 1 (low trump beats high off-suit)", p.TricksNS)
	}
}

func TestPlayHigherTrumpBeatsLowerTrump(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Hearts, Ace)}),
		East:  HandFromCards([]Card{NewCard(Hearts, King)}),
		South: HandFromCards([]Card{NewCard(Spades, Ace)}),
		West:  HandFromCards([]Card{NewCard(Hearts, Two)}),
	}

	contract := Contract{Level: 1, Strain: HeartStrain, Declarer: South}
	p := NewPlay(contract, hands)

	// West leads H2 (trump)
	must(t, p.PlayCard(West, NewCard(Hearts, Two)))
	// North plays HA (trump)
	must(t, p.PlayCard(South, NewCard(Hearts, Ace)))
	// East plays HK (trump)
	must(t, p.PlayCard(East, NewCard(Hearts, King)))
	// South has no hearts, plays SA (off-suit, doesn't beat trump)
	must(t, p.PlayCard(South, NewCard(Spades, Ace)))

	// North's HA wins (highest trump)
	if p.TricksNS != 1 {
		t.Errorf("TricksNS = %d, want 1 (HA highest trump)", p.TricksNS)
	}
}

func TestPlayNotation(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Spades, Ace)}),
		East:  HandFromCards([]Card{NewCard(Spades, King)}),
		South: HandFromCards([]Card{NewCard(Spades, Queen)}),
		West:  HandFromCards([]Card{NewCard(Spades, Jack)}),
	}

	contract := Contract{Level: 1, Strain: NoTrump, Declarer: South}
	p := NewPlay(contract, hands)

	must(t, p.PlayCard(West, NewCard(Spades, Jack)))
	must(t, p.PlayCard(South, NewCard(Spades, Ace)))
	must(t, p.PlayCard(East, NewCard(Spades, King)))
	must(t, p.PlayCard(South, NewCard(Spades, Queen)))

	want := "SJ SA SK SQ"
	if got := p.ToNotation(); got != want {
		t.Errorf("notation = %q, want %q", got, want)
	}
}

func TestPlayFullGame(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{
			NewCard(Spades, Ace), NewCard(Spades, King), NewCard(Spades, Queen),
			NewCard(Spades, Jack), NewCard(Spades, Ten), NewCard(Spades, Nine),
			NewCard(Spades, Eight), NewCard(Spades, Seven), NewCard(Spades, Six),
			NewCard(Spades, Five), NewCard(Spades, Four), NewCard(Spades, Three),
			NewCard(Spades, Two),
		}),
		East: HandFromCards([]Card{
			NewCard(Hearts, Ace), NewCard(Hearts, King), NewCard(Hearts, Queen),
			NewCard(Hearts, Jack), NewCard(Hearts, Ten), NewCard(Hearts, Nine),
			NewCard(Hearts, Eight), NewCard(Hearts, Seven), NewCard(Hearts, Six),
			NewCard(Hearts, Five), NewCard(Hearts, Four), NewCard(Hearts, Three),
			NewCard(Hearts, Two),
		}),
		South: HandFromCards([]Card{
			NewCard(Diamonds, Ace), NewCard(Diamonds, King), NewCard(Diamonds, Queen),
			NewCard(Diamonds, Jack), NewCard(Diamonds, Ten), NewCard(Diamonds, Nine),
			NewCard(Diamonds, Eight), NewCard(Diamonds, Seven), NewCard(Diamonds, Six),
			NewCard(Diamonds, Five), NewCard(Diamonds, Four), NewCard(Diamonds, Three),
			NewCard(Diamonds, Two),
		}),
		West: HandFromCards([]Card{
			NewCard(Clubs, Ace), NewCard(Clubs, King), NewCard(Clubs, Queen),
			NewCard(Clubs, Jack), NewCard(Clubs, Ten), NewCard(Clubs, Nine),
			NewCard(Clubs, Eight), NewCard(Clubs, Seven), NewCard(Clubs, Six),
			NewCard(Clubs, Five), NewCard(Clubs, Four), NewCard(Clubs, Three),
			NewCard(Clubs, Two),
		}),
	}

	contract := Contract{Level: 1, Strain: NoTrump, Declarer: South}
	p := NewPlay(contract, hands)

	for trick := 0; trick < NumTricks; trick++ {
		leader := p.Turn()
		leaderCards := p.RemainingHands[leader].Cards()
		if len(leaderCards) == 0 {
			t.Fatalf("trick %d: leader %v has no cards", trick+1, leader)
		}
		leadCard := leaderCards[0]

		for i := range NumDirections {
			player := Direction((int(leader) + i) % NumDirections)
			cards := p.RemainingHands[player].Cards()
			if len(cards) == 0 {
				t.Fatalf("trick %d: player %v has no cards", trick+1, player)
			}

			actor := player
			if player == p.Dummy {
				actor = p.Declarer
			}

			if i == 0 {
				must(t, p.PlayCard(actor, leadCard))
			} else {
				must(t, p.PlayCard(actor, cards[0]))
			}
		}
	}

	if !p.IsFinished() {
		t.Fatal("play should be finished after 13 tricks")
	}

	total := p.TricksNS + p.TricksEW
	if total != NumTricks {
		t.Errorf("total tricks = %d, want %d", total, NumTricks)
	}

	for d := Direction(0); d < NumDirections; d++ {
		if p.RemainingHands[d].Len() != 0 {
			t.Errorf("%v has %d cards remaining", d, p.RemainingHands[d].Len())
		}
	}
}

func TestPlayEWWinsTrick(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Spades, Two)}),
		East:  HandFromCards([]Card{NewCard(Spades, Ace)}),
		South: HandFromCards([]Card{NewCard(Spades, Three)}),
		West:  HandFromCards([]Card{NewCard(Spades, King)}),
	}

	contract := Contract{Level: 1, Strain: NoTrump, Declarer: South}
	p := NewPlay(contract, hands)

	// West leads SK
	must(t, p.PlayCard(West, NewCard(Spades, King)))
	// North (dummy) plays S2
	must(t, p.PlayCard(South, NewCard(Spades, Two)))
	// East plays SA
	must(t, p.PlayCard(East, NewCard(Spades, Ace)))
	// South plays S3
	must(t, p.PlayCard(South, NewCard(Spades, Three)))

	if p.TricksEW != 1 {
		t.Errorf("TricksEW = %d, want 1", p.TricksEW)
	}
	if p.TricksNS != 0 {
		t.Errorf("TricksNS = %d, want 0", p.TricksNS)
	}
}

func TestPlayOvertrump(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Hearts, Two)}),
		East:  HandFromCards([]Card{NewCard(Hearts, King)}),
		South: HandFromCards([]Card{NewCard(Spades, King)}),
		West:  HandFromCards([]Card{NewCard(Spades, Ace)}),
	}

	contract := Contract{Level: 1, Strain: HeartStrain, Declarer: South}
	p := NewPlay(contract, hands)

	// West leads SA
	must(t, p.PlayCard(West, NewCard(Spades, Ace)))
	// North (dummy) has no spades, ruffs with H2
	must(t, p.PlayCard(South, NewCard(Hearts, Two)))
	// East has no spades, overruffs with HK
	must(t, p.PlayCard(East, NewCard(Hearts, King)))
	// South follows suit with SK
	must(t, p.PlayCard(South, NewCard(Spades, King)))

	if p.TricksEW != 1 {
		t.Errorf("TricksEW = %d, want 1 (HK overruffs H2)", p.TricksEW)
	}
	if p.TricksNS != 0 {
		t.Errorf("TricksNS = %d, want 0", p.TricksNS)
	}
}

func TestPlayDiscardWhenVoid(t *testing.T) {
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{NewCard(Diamonds, Ace)}),
		East:  HandFromCards([]Card{NewCard(Spades, King)}),
		South: HandFromCards([]Card{NewCard(Diamonds, King)}),
		West:  HandFromCards([]Card{NewCard(Spades, Ace)}),
	}

	contract := Contract{Level: 1, Strain: HeartStrain, Declarer: South}
	p := NewPlay(contract, hands)

	// West leads SA
	must(t, p.PlayCard(West, NewCard(Spades, Ace)))
	// North (dummy) has no spades, discards DA
	err := p.PlayCard(South, NewCard(Diamonds, Ace))
	if err != nil {
		t.Errorf("discarding when void should be legal: %v", err)
	}
	// East follows suit with SK
	must(t, p.PlayCard(East, NewCard(Spades, King)))
	// South has no spades, discards DK
	must(t, p.PlayCard(South, NewCard(Diamonds, King)))

	// West wins with SA
	if p.TricksEW != 1 {
		t.Errorf("TricksEW = %d, want 1 (SA wins, others discarded)", p.TricksEW)
	}
}

func TestPlayDeclarerCantPlayDummyCardOnOwnTurn(t *testing.T) {
	p, _ := setupHeartsPlay()

	// West leads S5
	must(t, p.PlayCard(West, NewCard(Spades, Five)))
	// Dummy's turn -> declarer plays SA for dummy
	must(t, p.PlayCard(South, NewCard(Spades, Ace)))
	// East plays SJ
	must(t, p.PlayCard(East, NewCard(Spades, Jack)))

	// Trying to play a card from dummy's hand fails
	err := p.PlayCard(South, NewCard(Spades, King))
	if err == nil {
		t.Error("expected error: SK is in dummy's hand, not declarer's")
	}

	// Playing declarer's own card works
	must(t, p.PlayCard(South, NewCard(Spades, Eight)))
}

func TestPlayLeaderChangesAcrossMultipleTricks(t *testing.T) {
	// 3 tricks where the winner alternates between sides
	hands := [NumDirections]Hand{
		North: HandFromCards([]Card{
			NewCard(Spades, Two),
			NewCard(Hearts, Ace),
			NewCard(Diamonds, Two),
		}),
		East: HandFromCards([]Card{
			NewCard(Spades, Ace),
			NewCard(Hearts, Two),
			NewCard(Diamonds, King),
		}),
		South: HandFromCards([]Card{
			NewCard(Spades, Three),
			NewCard(Hearts, Three),
			NewCard(Diamonds, Three),
		}),
		West: HandFromCards([]Card{
			NewCard(Spades, King),
			NewCard(Hearts, King),
			NewCard(Diamonds, Ace),
		}),
	}

	contract := Contract{Level: 1, Strain: NoTrump, Declarer: South}
	p := NewPlay(contract, hands)

	// Trick 1: West leads SK, East plays SA -> East wins
	must(t, p.PlayCard(West, NewCard(Spades, King)))
	must(t, p.PlayCard(South, NewCard(Spades, Two)))
	must(t, p.PlayCard(East, NewCard(Spades, Ace)))
	must(t, p.PlayCard(South, NewCard(Spades, Three)))

	if p.Turn() != East {
		t.Errorf("trick 2 leader = %v, want East", p.Turn())
	}
	if p.TricksEW != 1 {
		t.Errorf("TricksEW = %d, want 1 after trick 1", p.TricksEW)
	}

	// Trick 2: East leads H2, North (dummy) plays HA -> North wins
	must(t, p.PlayCard(East, NewCard(Hearts, Two)))
	must(t, p.PlayCard(South, NewCard(Hearts, Three)))
	must(t, p.PlayCard(West, NewCard(Hearts, King)))
	must(t, p.PlayCard(South, NewCard(Hearts, Ace)))

	if p.Turn() != North {
		t.Errorf("trick 3 leader = %v, want North", p.Turn())
	}
	if p.TricksNS != 1 {
		t.Errorf("TricksNS = %d, want 1 after trick 2", p.TricksNS)
	}

	// Trick 3: North (dummy) leads D2, West plays DA -> West wins
	must(t, p.PlayCard(South, NewCard(Diamonds, Two)))
	must(t, p.PlayCard(East, NewCard(Diamonds, King)))
	must(t, p.PlayCard(South, NewCard(Diamonds, Three)))
	must(t, p.PlayCard(West, NewCard(Diamonds, Ace)))

	if p.TricksEW != 2 {
		t.Errorf("TricksEW = %d, want 2 after trick 3", p.TricksEW)
	}
	if p.TricksNS != 1 {
		t.Errorf("TricksNS = %d, want 1 after trick 3", p.TricksNS)
	}
	if !p.IsFinished() {
		t.Error("play should be finished after 3 tricks (3 cards each)")
	}
}
