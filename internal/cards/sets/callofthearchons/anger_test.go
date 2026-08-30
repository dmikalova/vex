package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Anger
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Ready and fight with a friendly creature.
func TestAnger(t *testing.T) {
	setup := func(t *testing.T) (h *ct.Harness, troll, witch ct.Card) {
		h = ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6)), // a bystander
					ct.Bind(&troll, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(8))),
				),
				Hand: ct.Cards(Anger),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&witch, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5)),
				),
			},
		})
		return h, troll, witch
	}

	t.Run("readies and fights with the chosen creature", func(t *testing.T) {
		h, troll, witch := setup(t)

		h.P1.Play(Anger)
		h.P1.ExpectAmber(1) // Anger's Æmber pip is gained on play
		h.P1.ExpectPrompt("Choose a friendly creature").Source("Anger")
		h.P1.ClickCard(troll)
		h.P1.ExpectPrompt("Choose a creature to fight").Source("Anger")
		h.P1.ClickCard(witch)

		h.Expect(troll).Exhausted().Damage(3).At(ct.PlayArea)
		h.Expect(witch).At(ct.Discard)
	})

	t.Run("can fight with an already-exhausted creature", func(t *testing.T) {
		h, troll, witch := setup(t)
		troll.Exhaust()

		h.P1.Play(Anger)
		h.P1.ClickCard(troll)
		h.P1.ClickCard(witch)

		h.Expect(troll).Exhausted() // readied, then fought, so exhausted again
		h.Expect(witch).At(ct.Discard)
	})
}
