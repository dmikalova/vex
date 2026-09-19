package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mercy, Malkin Queen
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Witch
//
//	Skirmish.
//	After a creature enters play, if it is a friendly Cat creature, ward it.
//	Fight: Ready a friendly Beast creature.
func TestMercyMalkinQueen(t *testing.T) {
	t.Run("wards a friendly Cat that enters play", func(t *testing.T) {
		var mercy, cat ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&mercy, MercyMalkinQueen)),
				Hand: ct.Cards(ct.Bind(&cat,
					ct.Creature(ct.OfHouse(card.House.Untamed), ct.Traits(card.Traits.Cat)))),
			},
		})

		h.P1.Play(cat)

		if !h.Game().Warded(cat.ID()) {
			t.Errorf("%s should be warded", cat.Name())
		}
	})

	t.Run("does not ward a non-Cat creature", func(t *testing.T) {
		var mercy, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&mercy, MercyMalkinQueen)),
				Hand:   ct.Cards(ct.Bind(&plain, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Play(plain)

		if h.Game().Warded(plain.ID()) {
			t.Errorf("%s should not be warded", plain.Name())
		}
	})

	t.Run("readies a friendly Beast on fight", func(t *testing.T) {
		var mercy, beast, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&mercy, MercyMalkinQueen),
					ct.Bind(&beast, ct.Creature(
						ct.OfHouse(card.House.Untamed),
						ct.Traits(card.Traits.Beast),
						ct.Power(4),
					)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2))))},
		})

		h.P1.Reap(beast)
		if !h.Game().Exhausted(beast.ID()) {
			t.Fatalf("%s should be exhausted after reaping", beast.Name())
		}

		h.P1.Fight(mercy, foe)

		if h.Game().Exhausted(beast.ID()) {
			t.Errorf("%s should be readied by Mercy's Fight ability", beast.Name())
		}
	})
}
