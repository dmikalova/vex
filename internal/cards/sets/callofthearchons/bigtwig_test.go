package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Bigtwig
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Beast
//
//	Bigtwig can only fight stunned creatures.
//	Reap: Stun and exhaust a creature.
func TestBigtwig(t *testing.T) {
	t.Run("reap stuns and exhausts a chosen creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(Bigtwig)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
		})

		h.P1.Reap(Bigtwig)
		h.P1.ClickCard(foe)

		h.Expect(foe).Stunned(true).Exhausted()
	})

	t.Run("can only fight stunned enemy creatures", func(t *testing.T) {
		var big, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(ct.Bind(&big, Bigtwig))},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
				),
			},
		})

		// An unstunned enemy is not a legal target.
		if err := h.Game().Fight(0, big.ID(), foe.ID()); err != engine.ErrNoTarget {
			t.Errorf("fight unstunned = %v, want ErrNoTarget", err)
		}

		foe.Stun()
		if err := h.Game().Fight(0, big.ID(), foe.ID()); err != nil {
			t.Fatalf("fight stunned: %v", err)
		}
		h.Expect(foe).At(ct.Discard)
	})
}
