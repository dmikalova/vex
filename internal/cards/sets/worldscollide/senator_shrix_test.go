package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Senator Shrix
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Dinosaur • Politician
//
//	You may spend Æmber on Senator Shrix when forging keys.
//	Play/Reap: You may exalt Senator Shrix.
func TestSenatorShrix(t *testing.T) {
	t.Run("may exalt itself when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, Hand: ct.Cards(SenatorShrix)},
			P2: ct.Side{},
		})

		h.P1.Play(SenatorShrix)
		h.P1.ClickOption("Yes")

		h.Expect(SenatorShrix).AmberOn(1)
	})

	t.Run("declining the exalt leaves it bare", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, Hand: ct.Cards(SenatorShrix)},
			P2: ct.Side{},
		})

		h.P1.Play(SenatorShrix)
		h.P1.ClickOption("No")

		h.Expect(SenatorShrix).AmberOn(0)
	})
}
