package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// J43G3R V
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  8
//	Armor:  2
//	Traits: Robot
//
//	Reap: Reap with 2 non-Star Alliance creatures, one at a time.
//	Fight: Fight with 2 non-Star Alliance creatures, one at a time.
func TestJ43G3RV(t *testing.T) {
	var jaeger, ally1, ally2 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.StarAlliance,
			InPlay: ct.Cards(
				ct.Bind(&jaeger, J43G3RV),
				ct.Bind(&ally1, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				ct.Bind(&ally2, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
			),
		},
	})

	h.P1.Reap(jaeger)
	h.P1.ClickCard(ally1)

	// Reaping with J43G3R V then reaps with both non-Star Alliance creatures.
	h.P1.ExpectAmber(3)
	h.Expect(ally1).Exhausted()
	h.Expect(ally2).Exhausted()
}
