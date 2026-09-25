package anomalyexpansion

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Shard of Glory
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Connected
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, exalt an enemy creature.
func TestShardOfGlory(t *testing.T) {
	t.Run("exalts an enemy creature for each friendly shard", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ShardOfGlory,
					ct.Artifact(ct.Traits(card.Traits.Shard)),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.Power(4))),
				),
			},
		})

		h.P1.UseAction(ShardOfGlory)

		h.Expect(enemy).AmberOn(2) // two friendly Shards exalt the lone enemy twice
	})
}
