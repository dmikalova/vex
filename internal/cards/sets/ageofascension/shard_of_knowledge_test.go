package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shard of Knowledge
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, draw a card.
func TestShardOfKnowledge(t *testing.T) {
	t.Run("draws a card for each friendly shard", func(t *testing.T) {
		var top, next ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ShardOfKnowledge,
					ct.Artifact(ct.Traits(card.Traits.Shard)),
				),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.Power(1))),
					ct.Bind(&next, ct.Creature(ct.Power(1))),
				),
			},
		})

		h.P1.UseAction(ShardOfKnowledge)

		h.Expect(top).At(ct.Hand)
		h.Expect(next).At(ct.Hand)
	})
}
