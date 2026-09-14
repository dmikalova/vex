package anomalyexpansion

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shard of Unity
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Connected
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, use a friendly Creature.
func TestShardOfUnity(t *testing.T) {
	t.Run("uses a friendly creature for each friendly shard", func(t *testing.T) {
		var friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ShardOfUnity,
					ct.Bind(&friend, ct.Creature(
						ct.OfHouse(card.House.StarAlliance),
						ct.Power(4),
					)),
					// A second ready creature makes which one is used a real choice.
					ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
				),
			},
		})

		h.P1.UseAction(ShardOfUnity)
		h.P1.ClickCard(friend) // Shard of Unity is the one friendly Shard: use one creature.

		h.P1.ExpectAmber(1) // the used creature reaps for 1 Æmber
	})
}
