package engine

import "testing"

func TestExhaustCreatures(t *testing.T) {
	t.Run("exhausts creatures the controller chooses, up to the max", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 3), 0)
		b := g.AddToBattleline(testCreature("b", 3), 0)
		g.SetChooser(0, optionPicker{idx: 0}) // always take the first candidate
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := ExhaustCreatures{Max: 3, Target: Target{Kind: TargetEachFriendlyCreature}}
		if e.Text() != "exhaust up to 3 creatures" {
			t.Errorf("text = %q", e.Text())
		}
		e.Resolve(ctx)
		if !g.State.Cards[a].Exhausted || !g.State.Cards[b].Exhausted {
			t.Error("both creatures should be exhausted")
		}
	})

	t.Run("stops when the controller is done", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 3), 0)
		g.SetChooser(0, optionPicker{idx: 5}) // out of range -> Done
		ctx := &EffectContext{Resolver: g, Controller: 0}

		ExhaustCreatures{Max: 3, Target: Target{Kind: TargetEachFriendlyCreature}}.Resolve(ctx)
		if g.State.Cards[a].Exhausted {
			t.Error("declining should exhaust nothing")
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (ExhaustCreatures{Max: 3}).validate() == nil {
			t.Error("unset target should be invalid")
		}
		if (ExhaustCreatures{Target: Target{Kind: TargetEachCreature}}).validate() == nil {
			t.Error("non-positive Max should be invalid")
		}
		if (ExhaustCreatures{Max: 1, Target: Target{Kind: TargetEachCreature}}).validate() != nil {
			t.Error("a valid ExhaustCreatures should pass")
		}
	})
}
