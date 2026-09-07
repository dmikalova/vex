package engine

import (
	"strings"
	"testing"
)

// A constant self-buff gated on the source being damaged applies only while it
// carries damage — Gron Nine-Toes's "+4 power while it is damaged."
func TestConstantAbilityWhileDamaged(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Gron", Brobnar, Creature, Rare, WithPower(8),
		WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetThisCreature}.Damaged(),
			PowerBonus: 4,
		}))

	if !strings.Contains(RenderCardRules(&def),
		"Gron gains +4 power while it is damaged.") {
		t.Errorf("card rules missing the while-damaged line:\n%s", RenderCardRules(&def))
	}

	gron := g.AddToBattleline(def, 0)

	// Undamaged: the buff is suspended.
	if got := g.Power(gron); got != 8 {
		t.Errorf("undamaged power = %d, want 8 (buff suspended)", got)
	}

	// One point of damage turns the buff on.
	g.applyRawDamage(DamageTarget{ID: gron, Amount: 1, IgnoreArmor: true})
	if got := g.Power(gron); got != 12 {
		t.Errorf("damaged power = %d, want 12 (8 + 4)", got)
	}
}

// A constant ability scaled PerTarget reads the count separately for each creature
// it reaches — Tribune Pompitus's "+2 power for each Æmber on it."
func TestConstantAbilityPerTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Tribune", Saurian, Creature, Rare, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetEachFriendlyCreature},
			PowerBonus: 2,
			PerTarget:  AemberOnIt,
		}))

	if !strings.Contains(RenderCardRules(&def),
		"Each friendly creature gains +2 power for each Æmber on it.") {
		t.Errorf("card rules missing the per-target line:\n%s", RenderCardRules(&def))
	}

	tribune := g.AddToBattleline(def, 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)

	// No Æmber on either: the bonus scales to zero.
	if got := g.Power(friend); got != 3 {
		t.Errorf("friend power with no Æmber = %d, want 3", got)
	}
	if got := g.Power(tribune); got != 4 {
		t.Errorf("tribune power with no Æmber = %d, want 4", got)
	}

	// Three Æmber on the friend: +2 per Æmber reaches only that creature.
	g.addAmberOn(friend, 3)
	if got := g.Power(friend); got != 9 {
		t.Errorf("friend power with 3 Æmber = %d, want 9 (3 + 2*3)", got)
	}
	if got := g.Power(tribune); got != 4 {
		t.Errorf("tribune power should still be 4 (its own Æmber is 0), got %d", got)
	}
}

// A constant ability scaled Per a source count names the source in its text and
// applies the same bonus to every creature it reaches — Primus Unguis's "+2 power
// for each Æmber on Primus Unguis."
func TestConstantAbilityPerSourceReachingOthers(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Primus Unguis", Saurian, Creature, Rare, WithPower(5),
		WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetEachFriendlyCreature},
			PowerBonus: 2,
			Per:        AemberOnThis{},
		}))

	if !strings.Contains(RenderCardRules(&def),
		"Each friendly creature gains +2 power for each Æmber on Primus Unguis.") {
		t.Errorf("card rules missing the per-source line:\n%s", RenderCardRules(&def))
	}

	primus := g.AddToBattleline(def, 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)

	// No Æmber on the source: the bonus scales to zero for everyone.
	if got := g.Power(friend); got != 3 {
		t.Errorf("friend power with no Æmber on source = %d, want 3", got)
	}

	// Two Æmber on the source reaches every friendly creature equally.
	g.addAmberOn(primus, 2)
	if got := g.Power(friend); got != 7 {
		t.Errorf("friend power = %d, want 7 (3 + 2*2)", got)
	}
	if got := g.Power(primus); got != 9 {
		t.Errorf("primus power = %d, want 9 (5 + 2*2)", got)
	}
}
