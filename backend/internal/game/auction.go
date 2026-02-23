package game

import (
	"fmt"
	"strings"
)

type Strain uint8
type CallType uint8

type Call struct {
	Type   CallType
	Level  uint8 // 1-7
	Strain Strain
}

type Contract struct {
	Level     uint8
	Strain    Strain
	Doubled   bool
	Redoubled bool
	Declarer  Direction
}

type Auction struct {
	Dealer Direction
	Calls  []Call

	// Cached state updated on each call
	turn      Direction
	lastBid   *Call
	lastBidBy Direction
	doubled   bool
	redoubled bool
	passCount int
	bidMade   bool
	finished  bool
}

const (
	ClubStrain Strain = iota
	DiamondStrain
	HeartStrain
	SpadeStrain
	NoTrump
)

const NumStrains = 5

const (
	Pass CallType = iota
	Bid
	Double
	Redouble
)

var strainLetters = [NumStrains]string{"C", "D", "H", "S", "NT"}
var passCall = Call{Type: Pass}
var doubleCall = Call{Type: Double}
var redoubleCall = Call{Type: Redouble}

func (s Strain) String() string {
	if s < NumStrains {
		return strainLetters[s]
	}
	return fmt.Sprintf("Strain(%d)", s)
}

func ParseStrain(s string) (Strain, bool) {
	switch strings.ToUpper(s) {
	case "C":
		return ClubStrain, true
	case "D":
		return DiamondStrain, true
	case "H":
		return HeartStrain, true
	case "S":
		return SpadeStrain, true
	case "NT", "N":
		return NoTrump, true
	default:
		return 0, false
	}
}

func BidCall(level uint8, strain Strain) Call {
	return Call{Type: Bid, Level: level, Strain: strain}
}

func (c Call) String() string {
	switch c.Type {
	case Pass:
		return "P"
	case Double:
		return "X"
	case Redouble:
		return "XX"
	case Bid:
		return fmt.Sprintf("%d%s", c.Level, c.Strain)
	default:
		return "?"
	}
}

func (c Call) bidIndex() int {
	return (int(c.Level)-1)*NumStrains + int(c.Strain)
}

func (c Contract) String() string {
	s := fmt.Sprintf("%d%s", c.Level, c.Strain)
	if c.Redoubled {
		s += "XX"
	} else if c.Doubled {
		s += "X"
	}
	s += " by " + c.Declarer.String()
	return s
}

func NewAuction(dealer Direction) *Auction {
	return &Auction{
		Dealer: dealer,
		turn:   dealer,
	}
}

func (a *Auction) Turn() Direction  { return a.turn }
func (a *Auction) IsFinished() bool { return a.finished }

func (a *Auction) MakeCall(dir Direction, call Call) error {
	if a.finished {
		return fmt.Errorf("auction is already finished")
	}
	if dir != a.turn {
		return fmt.Errorf("not %v's turn (expected %v)", dir, a.turn)
	}
	if err := a.validateCall(call); err != nil {
		return err
	}

	a.Calls = append(a.Calls, call)

	switch call.Type {
	case Pass:
		a.passCount++

	case Bid:
		a.lastBid = &Call{Type: Bid, Level: call.Level, Strain: call.Strain}
		a.lastBidBy = dir
		a.doubled = false
		a.redoubled = false
		a.bidMade = true
		a.passCount = 0

	case Double:
		a.doubled = true
		a.passCount = 0

	case Redouble:
		a.redoubled = true
		a.passCount = 0
	}

	a.turn = a.turn.Next()

	if (!a.bidMade && a.passCount == 4) || (a.bidMade && a.passCount == 3) {
		a.finished = true
	}

	return nil
}

func (a *Auction) validateCall(call Call) error {
	switch call.Type {
	case Pass:
		return nil

	case Bid:
		if call.Level < 1 || call.Level > 7 {
			return fmt.Errorf("bid level %d out of range (1-7)", call.Level)
		}
		if call.Strain >= NumStrains {
			return fmt.Errorf("invalid strain %d", call.Strain)
		}
		if a.lastBid != nil && call.bidIndex() <= a.lastBid.bidIndex() {
			return fmt.Errorf("%s is not higher than the current bid %s", call, a.lastBid)
		}
		return nil

	case Double:
		if a.lastBid == nil {
			return fmt.Errorf("cannot double: no bid to double")
		}
		if a.doubled || a.redoubled {
			return fmt.Errorf("cannot double: already doubled")
		}
		if a.isSameSide(a.turn, a.lastBidBy) {
			return fmt.Errorf("cannot double your own side's bid")
		}
		return nil

	case Redouble:
		if !a.doubled {
			return fmt.Errorf("cannot redouble: not doubled")
		}
		if a.redoubled {
			return fmt.Errorf("cannot redouble: already redoubled")
		}
		if !a.isSameSide(a.turn, a.lastBidBy) {
			return fmt.Errorf("cannot redouble: the double is on the opponent's bid")
		}
		return nil

	default:
		return fmt.Errorf("unknown call type %d", call.Type)
	}
}

func (a *Auction) isSameSide(d1, d2 Direction) bool {
	return d1 == d2 || d1 == d2.Partner()
}

func (a *Auction) PassedOut() bool {
	return a.finished && !a.bidMade
}

func (a *Auction) Contract() (Contract, bool) {
	if !a.finished || a.PassedOut() {
		return Contract{}, false
	}

	declarerSide := a.lastBidBy
	strain := a.lastBid.Strain

	declarer := declarerSide
	caller := a.Dealer
	for _, c := range a.Calls {
		if c.Type == Bid && c.Strain == strain && a.isSameSide(caller, declarerSide) {
			declarer = caller
			break
		}
		caller = caller.Next()
	}

	return Contract{
		Level:     a.lastBid.Level,
		Strain:    strain,
		Doubled:   a.doubled,
		Redoubled: a.redoubled,
		Declarer:  declarer,
	}, true
}

func ParseCall(s string) (Call, bool) {
	s = strings.ToUpper(s)
	switch s {
	case "P", "PASS":
		return passCall, true
	case "X", "DBL":
		return doubleCall, true
	case "XX", "RDBL":
		return redoubleCall, true
	}

	if len(s) < 2 || s[0] < '1' || s[0] > '7' {
		return Call{}, false
	}

	level := uint8(s[0] - '0')
	strain, ok := ParseStrain(s[1:])
	if !ok {
		return Call{}, false
	}
	return BidCall(level, strain), true
}

func (a *Auction) ToNotation() string {
	parts := make([]string, len(a.Calls))
	for i, c := range a.Calls {
		parts[i] = c.String()
	}
	return strings.Join(parts, " ")
}
