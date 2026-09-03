package engine

import (
	"strings"
	"testing"
)

func TestLoseKeyword(t *testing.T) {
	e := LoseKeyword{Keyword: Elusive}
	if got := e.Text(); got != "for the remainder of the turn, each creature loses elusive" {
		t.Errorf("text = %q", got)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate: %v", err)
	}
	// A keyword outside the rulebook's set is now a compile error, so the only
	// unlosable keyword left to reject is the unset zero value.
	if err := (LoseKeyword{}).validate(); err == nil {
		t.Error("an unset keyword should not be losable")
	}

	g := started(t)
	hider := g.AddToBattleline(testCreature("hider", 3, WithKeywords(Elusive)), 1)
	if !g.HasKeyword(hider, Elusive) {
		t.Fatal("the creature should start elusive")
	}
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.HasKeyword(hider, Elusive) {
		t.Error("elusive should be lost for the remainder of the turn")
	}
	g.EndPlayPhase(0)
	if !g.HasKeyword(hider, Elusive) {
		t.Error("elusive should come back when the turn ends")
	}
}

func TestAttackIgnoresKeywords(t *testing.T) {
	ape := func() CardDefinition {
		return testCreature("ape", 3, WithAttackIgnores(Taunt, Elusive))
	}
	if got := attackIgnoresText(&CardDefinition{
		Name:          "Niffle Ape",
		AttackIgnores: []Keyword{Taunt, Elusive},
	}); got != "While Niffle Ape is attacking, ignore taunt and elusive." {
		t.Errorf("text = %q", got)
	}
	if got := attackIgnoresText(&CardDefinition{Name: "plain"}); got != "" {
		t.Errorf("a creature that ignores nothing should print nothing, got %q", got)
	}

	g := started(t)
	attacker := g.AddToBattleline(ape(), 0)
	def := ape()
	if got := RenderCardText(&def); !strings.HasSuffix(
		got, "While ape is attacking, ignore taunt and elusive.") {
		t.Errorf("card text = %q", got)
	}
	g.AddToBattleline(testCreature("taunter", 4, WithKeywords(Taunt)), 1)
	shielded := g.AddToBattleline(testCreature("shielded", 1), 1)
	hidden := g.AddToBattleline(testCreature("hidden", 1, WithKeywords(Elusive)), 1)

	// A neighbor of a taunter is out of reach for an ordinary attacker but not for
	// one that ignores taunt.
	plain := g.AddToBattleline(testCreature("plain", 3), 0)
	if !g.protectedByTaunt(plain, shielded) {
		t.Error("taunt should shield the neighbor from an ordinary attacker")
	}
	if g.protectedByTaunt(attacker, shielded) {
		t.Error("an attacker ignoring taunt should reach the neighbor")
	}

	// Elusive does not stop the ape's damage, so the 1-power creature dies.
	g.fight(attacker, hidden)
	if g.inPlay(hidden) {
		t.Error("an attacker ignoring elusive should destroy the hidden creature")
	}
}

func TestShoulderArmorStatic(t *testing.T) {
	armor := NewCard("Shoulder Armor", Sanctum, Upgrade, Uncommon,
		WithStatic(StaticModifier{PowerBonus: 2, ArmorBonus: 2, WhileOnFlank: true}))
	if got := staticText(armor.Static); got !=
		"While this creature is on a flank, it gains +2 power and +2 armor." {
		t.Errorf("text = %q", got)
	}
	if got := upgradeStaticLines(&armor, true); len(got) != 1 ||
		got[0] != "While on a flank, +2 power and +2 armor." {
		t.Errorf("hosted lines = %v", got)
	}

	g := started(t)
	flanker := g.AddToBattleline(testCreature("flanker", 3), 0)
	middle := g.AddToBattleline(testCreature("middle", 3), 0)
	g.AddToBattleline(testCreature("other flank", 3), 0)
	g.AttachUpgrade(flanker, g.Register(armor, 0))
	g.AttachUpgrade(middle, g.Register(armor, 0))

	if got := g.Power(flanker); got != 5 {
		t.Errorf("flanker power = %d, want 5", got)
	}
	if got := g.Armor(flanker); got != 2 {
		t.Errorf("flanker armor = %d, want 2", got)
	}
	if got := g.Power(middle); got != 3 {
		t.Errorf("middle power = %d, want 3", got)
	}
	if got := g.Armor(middle); got != 0 {
		t.Errorf("middle armor = %d, want 0", got)
	}
}

func TestArtifactOrFlankCreatureText(t *testing.T) {
	got := Target{Kind: TargetChosenCreatureOrArtifact}.OnFlank().Text()
	if got != "an artifact or flank creature" {
		t.Errorf("text = %q", got)
	}
}
