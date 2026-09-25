package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Purifier of Souls
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  2
//	Traits: Human • Priest
//
//	Destroyed effects cannot trigger.
func TestPurifierOfSouls(t *testing.T) {
	t.Run("Destroyed abilities do not fire while Purifier of Souls is in play", func(t *testing.T) {
		var brabble ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(PurifierOfSouls, ct.Bind(&brabble, Brabble)),
			},
			P2: ct.Side{Amber: 5},
		})

		h.Game().DestroyEach(0, []engine.LocalID{brabble.ID()})

		h.P2.ExpectAmber(5)
	})
}
