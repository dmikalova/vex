package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Safe House
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	Action: Archive a friendly creature from play.
func TestSafeHouse(t *testing.T) {
	t.Run("archives a friendly creature from play as an action", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					SafeHouse,
					ct.Bind(&ally, ct.Creature(ct.Power(4))),
				),
			},
		})

		h.P1.UseAction(SafeHouse)

		h.Expect(ally).At(ct.Archives)
	})
}
