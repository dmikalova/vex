package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Matter Maker
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	You may play upgrades as if they were in the active house.
func TestMatterMaker(t *testing.T) {
	upgrade := ct.Upgrade(ct.OfHouse(card.House.Logos))
	host := ct.Creature(ct.OfHouse(card.House.StarAlliance))

	t.Run("lets an off-house upgrade be played as if in the active house", func(t *testing.T) {
		var up ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(MatterMaker, host),
				Hand:   ct.Cards(ct.Bind(&up, upgrade)),
			},
			P2: ct.Side{},
		})

		h.P1.Play(up)
		h.Expect(up).At(ct.Attached)
	})

	t.Run("without Matter Maker the off-house upgrade cannot be played", func(t *testing.T) {
		var up ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(host),
				Hand:   ct.Cards(ct.Bind(&up, upgrade)),
			},
			P2: ct.Side{},
		})

		h.P1.ExpectCannotPlay(up)
	})
}
