package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Blossom Drake
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dragon
//
//	Blossom Drake gains +1 power for each artifact in play.
//	Each artifact's text box is considered blank, except for traits.
func TestBlossomDrake(t *testing.T) {
	t.Run("gets +1 power for each artifact in play", func(t *testing.T) {
		var drake ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&drake, BlossomDrake),
					ct.Artifact(ct.OfHouse(card.House.Untamed)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Artifact(ct.OfHouse(card.House.Brobnar)))},
		})

		// Base 4 power, +1 for each of the two artifacts in play.
		h.Expect(drake).Power(6)
	})
}
