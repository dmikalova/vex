package engine

import (
	"strings"
	"testing"
)

// A directional constant ability reaches only the creatures on the named side of
// its source in the battleline: Panpaca, Anga buffs those to its right, Jaga
// those to its left.
func TestConstantAbilityDirectional(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 4), 0)
	anga := g.AddToBattleline(NewCard("Anga", Untamed, Creature, Common,
		WithPower(5), WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetEachCreature}.ToRightOfSource(),
			PowerBonus: 2,
		})), 0)
	right := g.AddToBattleline(testCreature("right", 4), 0)

	if got := g.Power(right); got != 6 {
		t.Errorf("right creature power = %d, want 6", got)
	}
	if got := g.Power(left); got != 4 {
		t.Errorf("left creature power = %d, want 4 (not to Anga's right)", got)
	}
	if got := g.Power(anga); got != 5 {
		t.Errorf("Anga power = %d, want 5 (does not buff itself)", got)
	}

	def := g.cat.def(anga)
	if !strings.Contains(RenderCardRules(def),
		"to the right of Anga") {
		t.Errorf("Anga rules should name the right side: %q", RenderCardRules(def))
	}
}

// The left-facing variant reaches only creatures to its source's left, and a
// creature off the source's battleline is on neither side.
func TestConstantAbilityDirectionalLeft(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 4), 0)
	jaga := g.AddToBattleline(NewCard("Jaga", Untamed, Creature, Common,
		WithPower(3), WithConstantAbility(ConstantAbility{
			Target:   Target{Kind: TargetEachCreature}.ToLeftOfSource(),
			Keywords: []Keyword{Skirmish},
		})), 0)
	right := g.AddToBattleline(testCreature("right", 4), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 4), 1)

	if !g.hasKeyword(left, Skirmish) {
		t.Error("creature to Jaga's left should gain skirmish")
	}
	if g.hasKeyword(right, Skirmish) {
		t.Error("creature to Jaga's right should not gain skirmish")
	}
	if g.hasKeyword(enemy, Skirmish) {
		t.Error("enemy in another battleline is on neither side of Jaga")
	}
	if !strings.Contains(RenderCardRules(g.cat.def(jaga)), "to the left of Jaga") {
		t.Error("Jaga rules should name the left side")
	}
}
