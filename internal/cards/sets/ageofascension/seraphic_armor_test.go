package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Seraphic Armor
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains +1 armor.
//	Play: Fully heal this creature.
func TestSeraphicArmor(t *testing.T) {
	t.Run("fully heals its host and grants armor when played", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(SeraphicArmor),
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(6))),
				),
			},
		})
		host.Damaged(4)

		h.P1.Play(SeraphicArmor)

		h.Expect(host).Damage(0)
		h.Expect(host).Armor(1)
	})
}
