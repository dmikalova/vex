package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Berinon
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  2
//	Traits: Spirit • Knight
//
//	After a creature enters play, if it is a Mutant creature, enrage Berinon.
//	Reap: Berinon captures 2 Æmber from your opponent.
func TestBerinon(t *testing.T) {
	t.Run("enrages after a Mutant creature enters play", func(t *testing.T) {
		var berinon, mutant ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&berinon, Berinon)),
				Hand: ct.Cards(
					ct.Bind(
						&mutant,
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Traits(card.Traits.Mutant)),
					),
				),
			},
		})

		h.P1.Play(mutant)

		if !h.Game().Enraged(berinon.ID()) {
			t.Errorf("%s should be enraged", berinon.Name())
		}
	})

	t.Run("does not enrage after a non-Mutant creature enters play", func(t *testing.T) {
		var berinon, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&berinon, Berinon)),
				Hand:   ct.Cards(ct.Bind(&plain, ct.Creature(ct.OfHouse(card.House.Sanctum)))),
			},
		})

		h.P1.Play(plain)

		if h.Game().Enraged(berinon.ID()) {
			t.Errorf("%s should not be enraged", berinon.Name())
		}
	})

	t.Run("captures 2 Æmber on reap", func(t *testing.T) {
		var berinon ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&berinon, Berinon)),
			},
			P2: ct.Side{Amber: 2},
		})

		h.P1.Reap(berinon)

		h.Expect(berinon).AmberOn(2)
		h.P2.ExpectAmber(0)
	})
}
