package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Good of the Many
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy each Creature that does not share a trait with another Creature in its controller's battleline.
func TestGoodOfTheMany(t *testing.T) {
	// Two friendly Beasts share a trait (survive). A friendly Knight shares a trait
	// only with an enemy Knight (dies — the enemy does not count). A friendly
	// traitless creature is a loner (dies). The enemy Knight is itself a loner on
	// its own side (dies).
	var beastA, beastB, knight, lone, enemyKnight ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			Hand:  ct.Cards(GoodOfTheMany),
			InPlay: ct.Cards(
				ct.Bind(&beastA, ct.Creature(ct.Traits(card.Traits.Beast))),
				ct.Bind(&beastB, ct.Creature(ct.Traits(card.Traits.Beast))),
				ct.Bind(&knight, ct.Creature(ct.Traits(card.Traits.Knight))),
				ct.Bind(&lone, ct.Creature()),
			),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&enemyKnight, ct.Creature(ct.Traits(card.Traits.Knight))),
		)},
	})

	h.P1.Play(GoodOfTheMany)

	h.Expect(beastA).At(ct.PlayArea)
	h.Expect(beastB).At(ct.PlayArea)
	h.Expect(knight).At(ct.Discard)
	h.Expect(lone).At(ct.Discard)
	h.Expect(enemyKnight).At(ct.Discard)
}
