package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Reassembling Automaton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Robot • Experiment
//
//	Destroyed: If you have any other creatures in play, instead of destroying Reassembling Automaton, fully heal it, and exhaust it, and move it to either flank of its controller's battleline.
func TestReassemblingAutomaton(t *testing.T) {
	t.Run("with another creature, survives fully healed on a flank", func(t *testing.T) {
		var automaton, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(
				ct.Bind(&automaton, ReassemblingAutomaton),
				ct.Creature(),
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(automaton, enemy)
		h.P1.ClickOption("left flank")

		h.Expect(automaton).At(ct.PlayArea)
		h.Expect(automaton).Damage(0)
	})

	t.Run("alone, is actually destroyed", func(t *testing.T) {
		var automaton, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(
				ct.Bind(&automaton, ReassemblingAutomaton),
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(automaton, enemy)

		h.Expect(automaton).At(ct.Discard)
	})
}
