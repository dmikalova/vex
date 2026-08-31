package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Key to Dis
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Key to Dis and each creature.
func TestKeyToDis(t *testing.T) {
	t.Run("sacrifices itself to destroy every creature in play", func(t *testing.T) {
		var friend, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					KeyToDis,
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.UseAction(KeyToDis)

		h.Expect(KeyToDis).At(ct.Discard) // sacrifices itself
		h.Expect(friend).At(ct.Discard)   // and destroys every creature
		h.Expect(foe).At(ct.Discard)
	})
}
