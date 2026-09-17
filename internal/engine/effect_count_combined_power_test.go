package engine

import (
	"strings"
	"testing"
)

// TestCombinedPowerOfNeighborsWithout covers the count's live sum and its two
// text renderings, and the blank-able "X" power it drives on a card.
func TestCombinedPowerOfNeighborsWithout(t *testing.T) {
	c := CombinedPowerOfNeighborsWithout{Without: Changeling}
	if got := c.CountText(); got != "combined power of {self}'s non-Changeling neighbors" {
		t.Errorf("CountText = %q", got)
	}
	if got := c.cardinalCountText(); got !=
		"the combined power of {self}'s non-Changeling neighbors" {
		t.Errorf("cardinalCountText = %q", got)
	}

	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 3), 0)
	picaroon := g.AddToBattleline(
		testCreature("Picaroon", 0, WithTraits(Mutant, Changeling), WithPowerX(c)), 0)
	// A Changeling neighbor is excluded from the tally.
	g.AddToBattleline(testCreature("changeling", 4, WithTraits(Changeling)), 0)
	_ = left

	// 3 (non-Changeling left) + 0 (excluded Changeling right) = 3.
	if got := g.Power(picaroon); got != 3 {
		t.Errorf("Picaroon power = %d, want 3", got)
	}

	// Blanking the card drops its X power to 0: the "X is …" line is part of its text.
	BlankEnemyText{}.Resolve(&EffectContext{Resolver: g, Source: picaroon, Controller: 1})
	if !g.textBlanked(picaroon) {
		t.Fatal("precondition: Picaroon should be blanked")
	}
	if got := g.Power(picaroon); got != 0 {
		t.Errorf("blanked Picaroon power = %d, want 0", got)
	}

	// The passive renders into the card's rules text.
	def := NewCard("Picaroon", Dis, Creature, Uncommon, WithPower(0), WithPowerX(c))
	if got := RenderCardRules(&def); !strings.Contains(got,
		"X is the combined power of Picaroon's non-Changeling neighbors.") {
		t.Errorf("rules = %q", got)
	}
}

// TestCombinedPowerOfNeighborsTerminatesOnCycle guards the mutual-reference
// hazard: two creatures whose X power reads each other (two Picaroons that both
// lost Changeling to Grey Aberrant) would recurse forever without the guard. The
// re-entrant creature contributes 0, so power stays finite and Power never
// overflows the stack.
func TestCombinedPowerOfNeighborsTerminatesOnCycle(t *testing.T) {
	c := CombinedPowerOfNeighborsWithout{Without: Changeling}
	g := NewGame("A", "B", 1)
	// Grey Aberrant has stripped Changeling, so these two X-power creatures now each
	// count the other. Line: anchor(3) — p1 — p2.
	anchor := g.AddToBattleline(testCreature("anchor", 3), 0)
	p1 := g.AddToBattleline(testCreature("p1", 0, WithPowerX(c)), 0)
	p2 := g.AddToBattleline(testCreature("p2", 0, WithPowerX(c)), 0)
	_ = anchor

	// p2's only non-Changeling neighbor is p1, which is mid-computation when p2 is
	// reached, so p1 contributes 0 there: p2 = 0. p1 = anchor(3) + p2(0) = 3.
	if got := g.Power(p2); got != 3 {
		t.Errorf("p2 power = %d, want 3", got)
	}
	if got := g.Power(p1); got != 3 {
		t.Errorf("p1 power = %d, want 3", got)
	}
	// The guard stack is balanced back to empty after each top-level read.
	if len(g.powerComputing) != 0 {
		t.Errorf("powerComputing stack left at %d entries, want 0", len(g.powerComputing))
	}
}
