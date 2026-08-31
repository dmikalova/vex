package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Yo Mama Mastery
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains taunt.
//	Play: Fully heal this creature.
func TestYoMamaMastery(t *testing.T) {
	t.Run("fully heals its host when played", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(YoMamaMastery),
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
			},
		})
		host.Damaged(4)

		h.P1.Play(YoMamaMastery)

		h.Expect(host).Damage(0) // fully healed
		h.P1.ExpectAmber(1)      // upgrade pip
	})
}
