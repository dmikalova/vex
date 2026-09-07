package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Baron Mengevin
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Human • Knight
//
//	After you discard a Sanctum card, Baron Mengevin captures 1 Æmber from your opponent.
func TestBaronMengevin(t *testing.T) {
	t.Run("captures 1 aember when a Sanctum card is discarded", func(t *testing.T) {
		var baron, fodder ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&baron, BaronMengevin)),
				Hand: ct.Cards(
					ct.Bind(&fodder, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Discard(fodder)

		h.Expect(baron).AmberOn(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("stays quiet for an off-house discard", func(t *testing.T) {
		var baron, fodder ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&baron, BaronMengevin)),
				Hand: ct.Cards(
					ct.Bind(&fodder, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Discard(fodder)

		h.Expect(baron).AmberOn(0)
		h.P2.ExpectAmber(3)
	})

	t.Run("each Baron captures on one Sanctum discard", func(t *testing.T) {
		var b1, b2, b3, fodder ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&b1, BaronMengevin),
					ct.Bind(&b2, BaronMengevin),
					ct.Bind(&b3, BaronMengevin),
				),
				Hand: ct.Cards(
					ct.Bind(&fodder, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Discard(fodder)

		h.Expect(b1).AmberOn(1)
		h.Expect(b2).AmberOn(1)
		h.Expect(b3).AmberOn(1)
		h.P2.ExpectAmber(0)
	})
}
