package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Little Rapscal
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Goblin
//
//	Elusive.
//	Creatures must fight when used, if able.
func TestLittleRapscal(t *testing.T) {
	t.Run("forces creatures with a fight target to fight rather than reap", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(LittleRapscal, ct.Bind(&ally, ct.Creature(ct.Power(4)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(ct.Power(3)))},
		})

		h.P1.ExpectCannotUseTo(ally, engine.ReapUse)
	})

	t.Run("leaves a creature with nothing to fight free to reap", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(LittleRapscal, ct.Bind(&ally, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.Reap(ally)
		h.Expect(ally).Exhausted()
	})
}
