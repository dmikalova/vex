package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ragnarok
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: You cannot use creatures to reap for the remainder of the turn. For the remainder of the turn, each time a friendly creature fights, gain 1 Æmber. At the end of the turn, destroy each creature.
func TestRagnarok(t *testing.T) {
	t.Run("bars friendly creatures from reaping this turn", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Ragnarok),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
			},
		})

		h.P1.Play(Ragnarok)

		h.P1.ExpectCannotUseTo(ally, card.UseKind.Reap)
	})

	t.Run("gains 1 Æmber each time a friendly creature fights this turn", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Ragnarok),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Play(Ragnarok)
		h.P1.Fight(ally, foe)

		h.P1.ExpectAmber(1)
	})

	t.Run("destroys each creature at the end of the turn", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(Ragnarok),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Play(Ragnarok)
		// Creatures survive until the turn actually ends.
		h.Expect(ally).At(ct.PlayArea)
		h.Expect(foe).At(ct.PlayArea)

		h.P1.EndTurn()

		h.Expect(ally).At(ct.Discard)
		h.Expect(foe).At(ct.Discard)
	})
}
