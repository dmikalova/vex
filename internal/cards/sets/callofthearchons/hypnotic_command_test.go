package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hypnotic Command
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: For each friendly Mars creature, an enemy creature captures 1 Æmber from their own side.
func TestHypnoticCommand(t *testing.T) {
	var foe1, foe2 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Mars,
			Hand:  ct.Cards(HypnoticCommand),
			InPlay: ct.Cards(
				ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
				ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4)),
			),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&foe1, ct.Creature(ct.Power(5))),
				ct.Bind(&foe2, ct.Creature(ct.Power(6))),
			),
			Amber: 3,
		},
	})

	h.P1.Play(HypnoticCommand)
	h.P1.ClickCard(foe1) // two friendly Mars creatures -> capture twice
	h.P1.ClickCard(foe2)

	h.Expect(foe1).AmberOn(1)
	h.Expect(foe2).AmberOn(1)
	h.P2.ExpectAmber(1) // 3 in pool - 2 captured onto their own creatures
}
