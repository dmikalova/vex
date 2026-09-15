package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Teleporter Chief Tink
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Alien
//
//	Elusive.
//	Action: Swap this creature with another friendly creature in your battleline. Use the other creature.
func TestTeleporterChiefTink(t *testing.T) {
	t.Run("swaps with a friendly creature and uses it", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					TeleporterChiefTink,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(TeleporterChiefTink)

		h.Expect(ally).Exhausted() // it was used to reap
		h.P1.ExpectAmber(1)
	})
}
