package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Destructive Analysis
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to a creature and purge any number of cards from your archives, and for each card purged this way, deal 2 damage to it.
func TestDestructiveAnalysis(t *testing.T) {
	t.Run("deals 2, then 2 more per card purged from archives", func(t *testing.T) {
		var analysis, foe, one, two ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:    card.House.Mars,
				Hand:     ct.Cards(ct.Bind(&analysis, DestructiveAnalysis)),
				Archives: ct.Cards(ct.Bind(&one, ct.Creature()), ct.Bind(&two, ct.Creature())),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Play(analysis)
		// Purge both archived cards, then stop.
		h.P1.ClickCard(one)
		h.P1.ClickCard(two)

		// 2 base + 2 per purged (2 purged) = 6.
		h.Expect(foe).Damage(6)
	})

	t.Run("purging nothing deals only the base 2 damage", func(t *testing.T) {
		var analysis, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:    card.House.Mars,
				Hand:     ct.Cards(ct.Bind(&analysis, DestructiveAnalysis)),
				Archives: ct.Cards(ct.Creature()),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Play(analysis)
		h.P1.ClickDone()

		h.Expect(foe).Damage(2)
	})

	t.Run("empty archives deals only the base 2 damage", func(t *testing.T) {
		var analysis, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&analysis, DestructiveAnalysis)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Play(analysis)

		h.Expect(foe).Damage(2)
	})
}
