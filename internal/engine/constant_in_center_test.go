package engine

import (
	"strings"
	"testing"
)

// InCenterOfBattleline is true only for the single middle creature of an
// odd-sized battleline; an even-sized line has no center, and a lone creature is
// its own center.
func TestInCenterOfBattleline(t *testing.T) {
	g := NewGame("A", "B", 1)

	// A lone creature is its own center.
	lone := g.AddToBattleline(testCreature("lone", 2), 0)
	if !g.InCenterOfBattleline(lone) {
		t.Error("a lone creature should be its own center")
	}

	// A second creature makes the line even: no center.
	second := g.AddToBattleline(testCreature("second", 2), 0)
	if g.InCenterOfBattleline(lone) || g.InCenterOfBattleline(second) {
		t.Error("an even-sized battleline should have no center")
	}

	// A third makes it odd again: only the middle creature is centered.
	third := g.AddToBattleline(testCreature("third", 2), 0)
	if !g.InCenterOfBattleline(second) {
		t.Error("the middle creature of an odd line should be centered")
	}
	if g.InCenterOfBattleline(lone) || g.InCenterOfBattleline(third) {
		t.Error("a flank creature of an odd line should not be centered")
	}
}

// A WhileInCenter constant self-buff applies only while the source is the middle
// creature of an odd-sized battleline, covering Kaloch Stonefather.
func TestConstantAbilityWhileInCenter(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Kaloch", Brobnar, Creature, Rare, WithPower(6),
		WithConstantAbility(ConstantAbility{
			Target:        Target{Kind: TargetEachFriendlyCreature},
			Keywords:      []Keyword{Skirmish},
			WhileInCenter: true,
		}))

	if !strings.Contains(
		RenderCardRules(&def),
		"While Kaloch is in the center of your battleline, each friendly Creature gains skirmish.",
	) {
		t.Error("card rules should render the while-in-center line")
	}

	left := g.AddToBattleline(testCreature("left", 3), 0)
	kaloch := g.AddToBattleline(def, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)

	// Kaloch is the middle creature, so the grant is active.
	if !g.hasKeyword(left, Skirmish) {
		t.Error("centered Kaloch should grant friendly creatures skirmish")
	}

	// Remove the right creature so the line is even: the grant lifts.
	g.State.Battleline[0].remove(right)
	if g.hasKeyword(left, Skirmish) {
		t.Error("off-center Kaloch should not grant skirmish")
	}
	_ = kaloch
}

// A WhileInCenter granted ability renders and reaches its target only while
// centered, covering The Shadow Council and Eldest Bear.
func TestConstantAbilityWhileInCenterGranted(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Council", Shadows, Creature, Rare, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target:        Target{Kind: TargetThisCreature},
			WhileInCenter: true,
			Granted: []Ability{{
				Trigger: TriggerAction,
				Effect:  StealAember{Amount: 2},
			}},
		}))

	if !strings.Contains(RenderCardRules(&def),
		`While Council is in the center of your battleline, it gains, "Action: Steal 2 Æmber."`) {
		t.Errorf("card rules should render the while-in-center granted line, got %q",
			RenderCardRules(&def))
	}

	left := g.AddToBattleline(testCreature("left", 3), 0)
	council := g.AddToBattleline(def, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	_ = left

	if !g.HasTrigger(council, TriggerAction) {
		t.Error("centered Council should have the granted action")
	}

	g.State.Battleline[0].remove(right)
	if g.HasTrigger(council, TriggerAction) {
		t.Error("off-center Council should not have the granted action")
	}
}

// SourceInCenterOfBattleline is met only while the source is the middle creature
// of an odd-sized battleline.
func TestSourceInCenterOfBattlelineCondition(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 0)
	mid := g.AddToBattleline(testCreature("mid", 2), 0)
	g.AddToBattleline(testCreature("right", 2), 0)

	c := SourceInCenterOfBattleline{}
	if c.CondText() != "if "+SelfName+" is in the center of your battleline" {
		t.Errorf("CondText = %q", c.CondText())
	}
	if !c.Met(&EffectContext{Resolver: g, Source: mid}) {
		t.Error("middle creature should satisfy SourceInCenterOfBattleline")
	}
	if c.Met(&EffectContext{Resolver: g, Source: left}) {
		t.Error("flank creature should not satisfy SourceInCenterOfBattleline")
	}
}
