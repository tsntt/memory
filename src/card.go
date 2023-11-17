package src

import "github.com/google/uuid"

type Card struct {
	ID      uuid.UUID
	Content string
	Turned  bool
}

func NewCard(content string) *Card {
	return &Card{
		ID:      uuid.New(),
		Content: content,
		Turned:  false,
	}
}

func MakeCardPairs(content ...string) []Card {
	cards := make([]Card, 0)

	for _, c := range content {
		cardA := NewCard(c)
		cardB := NewCard(c)

		cards = append(cards, *cardA, *cardB)
	}

	shuffleSlice(cards)

	return cards
}

func (c *Card) Turn() {
	c.Turned = !c.Turned
}
