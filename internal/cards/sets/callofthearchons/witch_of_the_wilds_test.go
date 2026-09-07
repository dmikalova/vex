package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Witch of the Wilds
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast • Witch
//
//	Each turn you may play one Untamed card.
func TestWitchOfTheWilds(t *testing.T) {
	t.Run("allows one Untamed card to be played off-house", func(t *testing.T) {
		var first ct.Card
		var second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(WitchOfTheWilds),
				Hand: ct.Cards(
					ct.Bind(&first, ct.Creature(ct.OfHouse(card.House.Untamed))),
					ct.Bind(&second, ct.Creature(ct.OfHouse(card.House.Untamed))),
				),
			},
		})

		h.P1.Play(first)

		h.Expect(first).At(ct.PlayArea)
		if _, err := h.Game().PlayCreature(0, 0, false); err != engine.ErrWrongHouse {
			t.Fatalf("second off-house Untamed play = %v, want ErrWrongHouse", err)
		}
		h.Expect(second).At(ct.Hand)
	})

	t.Run("is irrelevant when Untamed is the active house", func(t *testing.T) {
		var first ct.Card
		var second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(WitchOfTheWilds),
				Hand: ct.Cards(
					ct.Bind(&first, ct.Creature(ct.OfHouse(card.House.Untamed))),
					ct.Bind(&second, ct.Creature(ct.OfHouse(card.House.Untamed))),
				),
			},
		})

		h.P1.Play(first)
		h.P1.Play(second)

		h.Expect(first).At(ct.PlayArea)
		h.Expect(second).At(ct.PlayArea)
	})

	t.Run("does not allow off-house Untamed cards without Witch", func(t *testing.T) {
		var untamed ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ct.Bind(&untamed, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		if _, err := h.Game().PlayCreature(0, 0, false); err != engine.ErrWrongHouse {
			t.Fatalf("off-house Untamed play without Witch = %v, want ErrWrongHouse", err)
		}
		h.Expect(untamed).At(ct.Hand)
	})
}
