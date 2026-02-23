package game

import "fmt"

type Direction uint8
type Vulnerability uint8

const (
	VulNone Vulnerability = 0
	VulNS   Vulnerability = 1
	VulEW   Vulnerability = 2
	VulBoth Vulnerability = 3
)

const (
	North Direction = 0
	East  Direction = 1
	South Direction = 2
	West  Direction = 3

	NumDirections = 4
)

var directionNames = [NumDirections]string{"North", "East", "South", "West"}
var directionLetters = [NumDirections]byte{'N', 'E', 'S', 'W'}
var vulNames = [4]string{"None", "NS", "EW", "Both"}

func (d Direction) String() string {
	if d < NumDirections {
		return directionNames[d]
	}
	return fmt.Sprintf("Direction(%d)", d)
}

func (d Direction) Letter() byte {
	return directionLetters[d]
}

func (d Direction) Next() Direction {
	return (d + 1) % NumDirections
}

func (d Direction) Partner() Direction {
	return (d + 2) % NumDirections
}

func ParseDirection(b byte) (Direction, bool) {
	switch b {
	case 'N', 'n':
		return North, true
	case 'E', 'e':
		return East, true
	case 'S', 's':
		return South, true
	case 'W', 'w':
		return West, true
	default:
		return 0, false
	}
}

func (v Vulnerability) String() string {
	if v <= VulBoth {
		return vulNames[v]
	}
	return fmt.Sprintf("Vulnerability(%d)", v)
}

func (v Vulnerability) IsVulnerable(d Direction) bool {
	switch v {
	case VulNS:
		return d == North || d == South
	case VulEW:
		return d == East || d == West
	case VulBoth:
		return true
	default:
		return false
	}
}
