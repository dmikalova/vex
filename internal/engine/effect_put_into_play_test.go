package engine

import (
	"slices"
	"testing"
)

func TestPutIntoPlay(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		if got := (PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}, UnderYourControl: true}).Text(); got != "put it into play under your control" {
			t.Errorf("text = %q", got)
		}
		if got := (PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}}).Text(); got != "put it into play" {
			t.Errorf("text = %q", got)
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (PutIntoPlay{}).validate() == nil {
			t.Error("unset target should be invalid")
		}
		if (PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}}).validate() != nil {
			t.Error("a set target should be valid")
		}
	})

	t.Run("puts a discarded creature into play under your control, keeping ownership", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		foe := g.AddToDiscard(testCreature("foe", 3), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0, It: foe, HasIt: true}

		PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}, UnderYourControl: true}.Resolve(ctx)

		if !slices.Contains(g.Battleline(0), foe) {
			t.Error("the creature should enter P0's battleline")
		}
		if g.controller(foe) != 0 {
			t.Errorf("controller = %d, want 0", g.controller(foe))
		}
		if g.owner(foe) != 1 {
			t.Errorf("owner = %d, want 1 (unchanged)", g.owner(foe))
		}
	})

	t.Run("puts a card into play under its owner's control by default", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		c := g.AddToHand(testCreature("c", 3), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0, It: c, HasIt: true}

		PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}}.Resolve(ctx)

		if !slices.Contains(g.Battleline(1), c) {
			t.Error("the creature should enter its owner's battleline")
		}
		if g.controller(c) != 1 {
			t.Errorf("controller = %d, want 1", g.controller(c))
		}
	})

	t.Run("puts an artifact into play in the controller's artifact row", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		art := g.AddToDiscard(NewCard("relic", Mars, Artifact, Common), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0, It: art, HasIt: true}

		PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}, UnderYourControl: true}.Resolve(ctx)

		if !slices.Contains(g.Artifacts(0), art) {
			t.Error("the artifact should enter P0's artifact row")
		}
	})

	t.Run("a card already in play is not moved", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		c := g.AddToBattleline(testCreature("c", 3), 1)
		g.putIntoPlay(c, 0)
		if !slices.Contains(g.Battleline(1), c) {
			t.Error("an in-play card should stay where it is")
		}
	})

	t.Run("with no target, does nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	})
}
