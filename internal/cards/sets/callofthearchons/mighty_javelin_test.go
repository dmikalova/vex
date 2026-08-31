package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mighty Javelin
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Mighty Javelin. Deal 4 damage to a creature.
func TestMightyJavelin(t *testing.T) {
	t.Run("destroys itself and deals 4 damage to a creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(MightyJavelin)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(6))),
				),
			},
		})

		h.P1.UseAction(MightyJavelin)

		h.Expect(MightyJavelin).At(ct.Discard)
		h.Expect(foe).At(ct.PlayArea).Damage(4)
	})
}
