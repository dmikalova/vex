package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bonerot Venom
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "After this creature is used, deal 2 damage to this creature."
func TestBonerotVenom(t *testing.T) {
	t.Run("deals 2 damage to its host after the host reaps", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
						BonerotVenom,
					),
				),
			},
		})

		h.P1.Reap(host)

		h.Expect(host).Damage(2)
	})
}
