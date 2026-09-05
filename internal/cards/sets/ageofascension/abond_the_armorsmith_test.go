package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Abond the Armorsmith
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Each other friendly creature gains +1 armor.
//	Action: For the remainder of the turn, each other friendly creature gains +1 armor.
func TestAbondTheArmorsmith(t *testing.T) {
	t.Run("constant: other friendly creatures have +1 armor", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					AbondTheArmorsmith,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Armor(0))),
				),
			},
		})

		h.Expect(ally).Armor(1)
		h.Expect(AbondTheArmorsmith).Armor(0) // the source itself is unaffected
	})

	t.Run("action: grants another +1 armor for the turn", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					AbondTheArmorsmith,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Armor(0))),
				),
			},
		})

		h.P1.UseAction(AbondTheArmorsmith)

		h.Expect(ally).Armor(2)               // +1 constant, +1 from the Action
		h.Expect(AbondTheArmorsmith).Armor(0) // the source itself is unaffected
	})
}
