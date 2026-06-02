package bot

type Difficulty int

const (
	Easy Difficulty = iota + 1
	// Medium and Hard reserved for future levels
)

// Bot selects a call or card given the legal options.
type Bot interface {
	ChooseCall(legalCalls []string) string
	ChooseCard(legalCards []string) string
}

// New returns the bot implementation for the given difficulty.
func New(d Difficulty) Bot {
	switch d {
	default:
		return &randomBot{}
	}
}

// ParseDifficulty converts a client-supplied int to a Difficulty.
// Returns Easy and false if the value is unrecognised.
func ParseDifficulty(n int) (Difficulty, bool) {
	switch Difficulty(n) {
	case Easy:
		return Easy, true
	default:
		return Easy, false
	}
}
