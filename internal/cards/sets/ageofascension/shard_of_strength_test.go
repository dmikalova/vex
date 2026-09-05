package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shard of Strength
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, give a friendly creature 3 +1 power counters.
func TestShardOfStrength(t *testing.T) {
	t.Run(
		"gives a friendly creature 3 +1 power counters for each friendly shard",
		func(t *testing.T) {
			var ally ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Brobnar,
					InPlay: ct.Cards(
						ShardOfStrength,
						ct.Artifact(ct.Traits(card.Traits.Shard)),
						ct.Bind(&ally, ct.Creature(ct.Power(4))),
					),
				},
			})

			h.P1.UseAction(ShardOfStrength)

			// Two friendly Shards, 3 counters each: 4 + 6 = 10.
			h.Expect(ally).Power(10)
		},
	)
}
