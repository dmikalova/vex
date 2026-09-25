package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Grimlocus Dux
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  11
//	Armor:  2
//	Traits: Dinosaur • Soldier
//
//	Taunt.
//	Play: Exalt Grimlocus Dux 2 times.
func TestGrimlocusDux(t *testing.T) {
	t.Run("exalts itself twice when played", func(t *testing.T) {
		var dux ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(ct.Bind(&dux, GrimlocusDux)),
			},
		})

		h.P1.Play(GrimlocusDux)

		h.Expect(dux).AmberOn(2)
	})
}
