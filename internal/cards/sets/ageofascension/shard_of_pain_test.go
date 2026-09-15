package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shard of Pain
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, deal 1 damage to an enemy creature.
func TestShardOfPain(t *testing.T) {
	t.Run("deals 1 damage to an enemy creature for each friendly shard", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ShardOfPain,
					ct.Artifact(ct.Traits(card.Traits.Shard)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.UseAction(ShardOfPain)

		h.Expect(foe).Damage(2)
	})
}
