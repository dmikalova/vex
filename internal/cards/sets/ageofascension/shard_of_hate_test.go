package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shard of Hate
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, stun an enemy Creature.
func TestShardOfHate(t *testing.T) {
	t.Run("stuns an enemy creature for each friendly shard", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ShardOfHate),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.UseAction(ShardOfHate)

		h.Expect(enemy).Stunned(true)
	})
}
