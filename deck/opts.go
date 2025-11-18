package deck

import (
	"math/rand"
	"time"
)

func Shuffle(cards []Card) []Card {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	shuffled := make([]Card, len(cards))
	copy(shuffled, cards)

	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

func AddJokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{Suit: Joker})
		}
		return cards
	}
}
