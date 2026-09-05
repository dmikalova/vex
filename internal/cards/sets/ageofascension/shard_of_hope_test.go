package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shard of Hope
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, a friendly creature captures 1 Æmber from your opponent.
func TestShardOfHope(t *testing.T) {
	t.Run("a friendly creature captures 1 aember for each friendly shard", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ShardOfHope,
					ct.Artifact(ct.Traits(card.Traits.Shard)),
					ct.Bind(&ally, ct.Creature(ct.Power(4))),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(ShardOfHope)

		h.Expect(ally).AmberOn(2)
		h.P2.ExpectAmber(3)
	})
}
