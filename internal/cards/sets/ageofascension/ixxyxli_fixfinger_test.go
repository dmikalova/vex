package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ixxyxli Fixfinger
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  2
//	Traits: Martian • Scientist
//
//	Elusive.
//	Each other friendly Mars creature gains +1 armor.
func TestIxxyxliFixfinger(t *testing.T) {
	t.Run("gives other Martian creatures +1 armor while in play", func(t *testing.T) {
		var ixxyxli, martian, offhouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&ixxyxli, IxxyxliFixfinger),
					ct.Bind(&martian, ct.Creature(ct.OfHouse(card.House.Mars), ct.Armor(0))),
					ct.Bind(&offhouse, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Armor(0))),
				),
			},
			P2: ct.Side{},
		})

		h.Expect(martian).Armor(1)  // other Martian creature buffed
		h.Expect(offhouse).Armor(0) // non-Martian creature unaffected
		h.Expect(ixxyxli).Armor(2)  // does not buff itself
	})
}
