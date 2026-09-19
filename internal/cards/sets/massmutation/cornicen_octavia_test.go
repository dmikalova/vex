package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cornicen Octavia
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Action: Cornicen Octavia captures 2 Æmber from your opponent.
func TestCornicenOctavia(t *testing.T) {
	t.Run("captures 2 Æmber from the opponent", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(CornicenOctavia),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(CornicenOctavia)

		h.Expect(CornicenOctavia).AmberOn(2)
		h.P2.ExpectAmber(1)
	})
}
