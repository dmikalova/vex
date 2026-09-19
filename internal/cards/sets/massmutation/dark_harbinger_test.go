package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dark Harbinger
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant • Witch
//
//	After you play an Untamed tactic, ready Dark Harbinger.
func TestDarkHarbinger(t *testing.T) {
	t.Run("readies after you play an Untamed Tactic", func(t *testing.T) {
		var harbinger ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(RapidEvolution),
				InPlay: ct.Cards(
					ct.Bind(&harbinger, DarkHarbinger),
					ct.Creature(ct.OfHouse(card.House.Untamed)),
				),
			},
		})

		h.P1.Reap(harbinger)
		h.Expect(harbinger).Exhausted()

		// Rapid Evolution is an Untamed Tactic; it needs a target creature.
		h.P1.Play(RapidEvolution)
		h.P1.ClickCard(harbinger)

		h.Expect(harbinger).Ready()
	})

	t.Run("does not ready after a non-Tactic is played", func(t *testing.T) {
		var harbinger, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed)))),
				InPlay: ct.Cards(
					ct.Bind(&harbinger, DarkHarbinger),
				),
			},
		})

		h.P1.Reap(harbinger)
		h.Expect(harbinger).Exhausted()

		h.P1.Play(ally)

		h.Expect(harbinger).Exhausted()
	})
}
