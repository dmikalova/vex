package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Eye of Judgment
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Action: Purge a creature from a discard pile.
func TestEyeOfJudgment(t *testing.T) {
	t.Run("purges a creature from a discard pile", func(t *testing.T) {
		var prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(EyeOfJudgment)},
			P2: ct.Side{Discard: ct.Cards(
				ct.Bind(&prey, ct.Creature(ct.OfHouse(card.House.Brobnar))),
			)},
		})

		h.P1.UseAction(EyeOfJudgment)

		h.Expect(prey).At(ct.Purge)
	})
}
