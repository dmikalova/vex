package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Groggins
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  8
//	Traits: Giant
//
//	Groggins can only fight flank creatures.
func TestGroggins(t *testing.T) {
	t.Run("can only fight flank creatures", func(t *testing.T) {
		var groggins, left, middle ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&groggins, Groggins)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(20))),
					ct.Bind(&middle, ct.Creature(ct.Power(4))),
					ct.Creature(ct.Power(4)),
				),
			},
		})

		// A creature not on a flank is not a legal target.
		if err := h.Game().Fight(0, groggins.ID(), middle.ID()); err != engine.ErrNoTarget {
			t.Errorf("fight middle = %v, want ErrNoTarget", err)
		}

		// A flank creature is a legal target.
		if err := h.Game().Fight(0, groggins.ID(), left.ID()); err != nil {
			t.Fatalf("fight flank: %v", err)
		}
		h.Expect(left).Damage(8)
	})
}
