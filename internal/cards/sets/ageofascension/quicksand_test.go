package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Quicksand
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy the most powerful Creature controlled by each player who does not have a friendly ready Untamed Creature in play.
func TestQuicksand(t *testing.T) {
	var mine, theirBig, theirSmall ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand:  ct.Cards(Quicksand),
			// P1 controls a ready Untamed creature, so its board is spared.
			InPlay: ct.Cards(
				ct.Bind(&mine, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(8))),
			),
		},
		P2: ct.Side{
			// P2 controls no ready Untamed creature, so its most powerful is destroyed.
			InPlay: ct.Cards(
				ct.Bind(&theirBig, ct.Creature(ct.Power(6))),
				ct.Bind(&theirSmall, ct.Creature(ct.Power(3))),
			),
		},
	})

	h.P1.Play(Quicksand)

	h.Expect(mine).At(ct.PlayArea)
	h.Expect(theirBig).At(ct.Discard)
	h.Expect(theirSmall).At(ct.PlayArea)
}
