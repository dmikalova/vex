package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Self-Bolstering Automata
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Robot
//
//	Destroyed: If there is another friendly Creature in play, instead of destroying Self-Bolstering Automata, fully heal it, exhaust it, and move it to either flank of its controller's battleline -> give it two +1 power counters.
func TestSelfBolsteringAutomata(t *testing.T) {
	t.Run("with another creature, survives with two +1 power counters", func(t *testing.T) {
		var automata, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(
				ct.Bind(&automata, SelfBolsteringAutomata),
				ct.Creature(),
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(automata, enemy)
		h.P1.ClickOption("left flank")

		h.Expect(automata).At(ct.PlayArea)
		h.Expect(automata).Damage(0)
		if got := h.Game().Power(automata.ID()); got != 3 {
			t.Errorf("power = %d, want 3 (1 base + two +1 counters)", got)
		}
	})

	t.Run("alone, is actually destroyed", func(t *testing.T) {
		var automata, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, InPlay: ct.Cards(
				ct.Bind(&automata, SelfBolsteringAutomata),
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(automata, enemy)

		h.Expect(automata).At(ct.Discard)
	})
}
