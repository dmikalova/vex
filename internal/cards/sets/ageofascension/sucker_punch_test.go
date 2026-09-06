package ageofascension_test

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/cards/sets/ageofascension"
)

// Sucker Punch
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Alpha.
//	Play: Deal 2 damage to an enemy creature. If this damage destroys that creature, archive Sucker Punch.
func TestSuckerPunch(t *testing.T) {
	t.Run("archives itself when the damage destroys the creature", func(t *testing.T) {
		var punch ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ct.Bind(&punch, ageofascension.SuckerPunch)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Creature(ct.Power(2))),
			},
		})

		h.P1.Play(punch)

		g := h.Game()
		if slices.Contains(g.Discard(0), punch.ID()) {
			t.Error("Sucker Punch should not be in the discard pile")
		}
		if !slices.Contains(g.Archives(0), punch.ID()) {
			t.Error("Sucker Punch should be archived after destroying the creature")
		}
	})

	t.Run("goes to discard when the creature survives", func(t *testing.T) {
		var punch ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ct.Bind(&punch, ageofascension.SuckerPunch)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Creature(ct.Power(5))),
			},
		})

		h.P1.Play(punch)

		g := h.Game()
		if !slices.Contains(g.Discard(0), punch.ID()) {
			t.Error("Sucker Punch should be discarded when the creature survives")
		}
		if slices.Contains(g.Archives(0), punch.ID()) {
			t.Error("Sucker Punch should not be archived when the creature survives")
		}
	})
}
