package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gleeful Mayhem
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each house, deal 5 damage to a creature of that house.
func TestGleefulMayhem(t *testing.T) {
	var mayhem, dis, logos, otherDis ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand:  ct.Cards(ct.Bind(&mayhem, GleefulMayhem)),
			InPlay: ct.Cards(
				ct.Bind(&dis, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(10))),
				ct.Bind(&otherDis, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(10))),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(10))),
			),
		},
	})

	h.P1.Play(mayhem)
	// One Dis creature is chosen for the Dis house; the lone Logos creature is
	// forced for the Logos house. No other house has a creature in play.
	h.P1.ClickCard(dis)

	h.Expect(dis).Damage(5)
	h.Expect(otherDis).Damage(0)
	h.Expect(logos).Damage(5)
}
