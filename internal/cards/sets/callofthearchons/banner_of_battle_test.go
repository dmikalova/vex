package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Banner of Battle
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Each friendly creature gains +1 power.
func TestBannerOfBattle(t *testing.T) {
	t.Run("gives each friendly creature +1 power while in play", func(t *testing.T) {
		var friend, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(BannerOfBattle),
				InPlay: ct.Cards(
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(BannerOfBattle)

		h.Expect(friend).Power(4) // +1 from Banner
		h.Expect(enemy).Power(3)  // Banner buffs only friendly creatures
	})
}
