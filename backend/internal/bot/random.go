package bot

import "math/rand/v2"

// randomBot passes every bid and plays a random legal card.
type randomBot struct{}

func (b *randomBot) ChooseCall(legalCalls []string) string {
	return "P"
}

func (b *randomBot) ChooseCard(legalCards []string) string {
	return legalCards[rand.IntN(len(legalCards))]
}
