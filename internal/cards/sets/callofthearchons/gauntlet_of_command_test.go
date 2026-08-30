package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gauntlet of Command
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	Action: Ready and fight with a friendly creature.
func TestGauntletOfCommand(t *testing.T) {
	t.Run("readies and fights with a friendly creature", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					GauntletOfCommand,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)))),
			},
		})

		ally.Exhaust()
		h.P1.UseAction(GauntletOfCommand)

		h.Expect(foe).At(ct.Discard)
		h.Expect(ally).At(ct.PlayArea)
	})
}
