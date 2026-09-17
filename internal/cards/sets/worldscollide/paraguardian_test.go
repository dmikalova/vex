package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Paraguardian
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Reap: You may exalt Paraguardian. Ward Paraguardian.
func TestParaguardian(t *testing.T) {
	t.Run("exalts itself and wards its neighbors when accepted", func(t *testing.T) {
		var para, left, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Saurian))),
					ct.Bind(&para, Paraguardian),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Saurian))),
				),
			},
			P2: ct.Side{Amber: 1},
		})

		h.P1.Reap(para)
		h.P1.ClickCard(para)

		h.Expect(para).AmberOn(1)
		if !h.Game().State.Cards[left.ID()].Warded {
			t.Error("left neighbor should be warded")
		}
		if !h.Game().State.Cards[right.ID()].Warded {
			t.Error("right neighbor should be warded")
		}
	})

	t.Run("does nothing when declined", func(t *testing.T) {
		var para, left ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Saurian))),
					ct.Bind(&para, Paraguardian),
				),
			},
			P2: ct.Side{Amber: 1},
		})

		h.P1.Reap(para)
		h.P1.ClickDone()

		h.Expect(para).AmberOn(0)
		if h.Game().State.Cards[left.ID()].Warded {
			t.Error("neighbor should not be warded when declined")
		}
	})
}
