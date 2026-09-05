package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shard of Greed
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, gain 1 Æmber.
func TestShardOfGreed(t *testing.T) {
	t.Run("gains 1 aember for each friendly shard", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					ShardOfGreed,
					ct.Artifact(ct.Traits(card.Traits.Shard)),
				),
			},
		})

		h.P1.UseAction(ShardOfGreed)

		h.P1.ExpectAmber(2)
	})
}
