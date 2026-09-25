package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Project Z.Y.X.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Cyborg • Mutant
//
//	Fight/Reap: You may play a card from your archives.
func TestProjectZYX(t *testing.T) {
	t.Run("reap plays a card from the archives", func(t *testing.T) {
		var zyx, archived ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&zyx, ProjectZYX)),
				Archives: ct.Cards(
					ct.Bind(&archived, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(zyx)
		h.P1.ClickCard(zyx)

		h.Expect(archived).At(ct.PlayArea)
	})
}
