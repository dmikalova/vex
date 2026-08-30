package engine

import (
	"strings"
	"testing"
)

// A constant ability targeting friendly creatures buffs every friendly creature
// and no enemy, and the buff vanishes the moment the source card leaves play.
func TestConstantAbilityFriendlyTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	banner := g.Register(NewCard("Banner", Brobnar, Artifact, Rare,
		WithConstantAbility(ConstantAbility{PowerBonus: 1, Target: Target{Kind: TargetEachFriendlyCreature}})), 0)
	g.State.Artifacts[0].add(banner)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	if got := g.Power(friend); got != 4 {
		t.Errorf("friendly power = %d, want 4 (3 + constant)", got)
	}
	if got := g.Power(enemy); got != 3 {
		t.Errorf("enemy power = %d, want 3 (a friendly target does not reach enemies)", got)
	}

	g.State.Artifacts[0].remove(banner)
	if got := g.Power(friend); got != 3 {
		t.Errorf("friendly power after the source leaves = %d, want 3", got)
	}
}

// A constant ability targeting neighbors reaches only the source's immediate
// battleline neighbors.
func TestConstantAbilityNeighboringTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	bulwark := NewCard("Bulwark", Sanctum, Creature, Common, WithPower(4), WithArmor(2),
		WithConstantAbility(ConstantAbility{ArmorBonus: 2, Target: Target{Kind: TargetEachCreature}.Neighboring()}))
	left := g.AddToBattleline(testCreature("left", 3), 0)
	g.AddToBattleline(bulwark, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	far := g.AddToBattleline(testCreature("far", 3), 0)

	if got := g.Armor(left); got != 2 {
		t.Errorf("left neighbor armor = %d, want 2", got)
	}
	if got := g.Armor(right); got != 2 {
		t.Errorf("right neighbor armor = %d, want 2", got)
	}
	if got := g.Armor(far); got != 0 {
		t.Errorf("non-neighbor armor = %d, want 0", got)
	}
}

// A constant ability with no target reaches every creature in play, including
// the source itself and the enemy.
func TestConstantAbilityNoTargetReachesEveryone(t *testing.T) {
	g := NewGame("A", "B", 1)
	self := g.AddToBattleline(NewCard("Totem", Brobnar, Creature, Common, WithPower(2),
		WithConstantAbility(ConstantAbility{PowerBonus: 1})), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	if got := g.Power(self); got != 3 {
		t.Errorf("source power = %d, want 3 (it buffs itself)", got)
	}
	if got := g.Power(friend); got != 4 {
		t.Errorf("friend power = %d, want 4", got)
	}
	if got := g.Power(enemy); got != 4 {
		t.Errorf("enemy power = %d, want 4 (an untargeted constant reaches everyone)", got)
	}
}

// An untargeted constant ability's scope also reaches artifacts in play, not just
// creatures (creatures have no artifact-relevant stat yet, so this is about the
// scope, verified through constantAffects).
func TestConstantAbilityReachesArtifacts(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(NewCard("Totem", Brobnar, Creature, Common, WithPower(2),
		WithConstantAbility(ConstantAbility{PowerBonus: 1})), 0)
	art := g.AddArtifact(exAutocannon(), 1)

	if !g.constantAffects(src, g.cat.def(src).Constant, art) {
		t.Error("an untargeted constant ability should reach artifacts in play")
	}
}

// A constant ability that grants a trigger gives it only to the creatures its
// Target reaches: hasTrigger sees it on a reached creature and skips a creature
// the grantor does not affect.
func TestHasTriggerFromConstantAbility(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("Grantor", Brobnar, Creature, Common, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target:  Target{Kind: TargetEachFriendlyCreature},
			Granted: []Ability{{Trigger: TriggerAfterReap, Effect: GainAember{Player: Controller, Amount: 1}}},
		})), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	if !g.hasTrigger(friend, TriggerAfterReap) {
		t.Error("a friendly creature should gain the granted Reap trigger")
	}
	if g.hasTrigger(enemy, TriggerAfterReap) {
		t.Error("an enemy creature is not reached by the friendly constant ability")
	}
}

func TestHasKeywordFromConstantAbility(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("Grantor", Untamed, Creature, Common, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target:   Target{Kind: TargetEachFriendlyCreature},
			Keywords: []Keyword{Skirmish},
		})), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	if !g.hasKeyword(friend, Skirmish) {
		t.Error("a friendly creature should gain the granted keyword")
	}
	if g.hasKeyword(friend, Poison) {
		t.Error("only the granted keyword is gained")
	}
	if g.hasKeyword(enemy, Skirmish) {
		t.Error("an enemy creature is not reached by the friendly constant ability")
	}
}

func TestConstantText(t *testing.T) {
	banner := NewCard("Banner", Brobnar, Artifact, Rare,
		WithConstantAbility(ConstantAbility{PowerBonus: 1, Target: Target{Kind: TargetEachFriendlyCreature}}))
	if got := constantText(&banner); got != "Each friendly creature gains +1 power." {
		t.Errorf("friendly constant text = %q", got)
	}
	if !strings.Contains(RenderCardText(&banner), "Each friendly creature gains +1 power.") {
		t.Error("RenderCardText should include the constant-ability line")
	}

	bulwark := NewCard("Bulwark", Sanctum, Creature, Common, WithPower(4),
		WithConstantAbility(ConstantAbility{ArmorBonus: 2, Target: Target{Kind: TargetEachCreature}.Neighboring()}))
	if got := constantText(&bulwark); got != "Each neighboring creature gains +2 armor." {
		t.Errorf("neighboring constant text = %q", got)
	}

	all := NewCard("Totem", Brobnar, Creature, Common, WithPower(2),
		WithConstantAbility(ConstantAbility{PowerBonus: 1}))
	if got := constantText(&all); got != "Each card in play gains +1 power." {
		t.Errorf("untargeted constant text = %q", got)
	}

	table := NewCard("Table", Sanctum, Artifact, Rare,
		WithConstantAbility(ConstantAbility{PowerBonus: 1, Keywords: []Keyword{Taunt}, Target: Target{Kind: TargetEachFriendlyCreature}}))
	if got := constantText(&table); got != "Each friendly creature gains +1 power and taunt." {
		t.Errorf("keyword constant text = %q", got)
	}
	kwOnly := NewCard("Halo", Untamed, Creature, Common, WithPower(4),
		WithConstantAbility(ConstantAbility{Keywords: []Keyword{Skirmish}, Target: Target{Kind: TargetEachFriendlyCreature}}))
	if got := constantText(&kwOnly); got != "Each friendly creature gains skirmish." {
		t.Errorf("keyword-only constant text = %q", got)
	}

	plain := NewCard("Plain", Brobnar, Creature, Common, WithPower(3))
	if got := constantText(&plain); got != "" {
		t.Errorf("no-constant text = %q, want empty", got)
	}
}
