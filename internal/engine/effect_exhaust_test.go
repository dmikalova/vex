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
		if e.Text() != "exhaust up to 3 friendly creatures" {
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

// TestExhaustGate covers Exhaust as a binding result gate (Humble): it exhausts
// the target, binds it in context, and reports progress so a Then hangs off it.
func TestExhaustGate(t *testing.T) {
	t.Run("exhausts, binds the creature, and reports progress", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		c := g.AddToBattleline(testCreature("c", 4), 0)
		ctx := &EffectContext{Resolver: g, Controller: 0}

		e := Exhaust{Target: Target{Kind: TargetThisCreature}, Bind: true}
		if !e.resolveGate(ctx) {
			t.Error("exhausting a creature should report progress")
		}
		if !g.State.Cards[c].Exhausted {
			t.Error("the creature should be exhausted")
		}
		if !ctx.HasIt || ctx.It != c {
			t.Errorf("ctx.It = %v (HasIt %v), want %v bound", ctx.It, ctx.HasIt, c)
		}
	})

	t.Run("no target reports no progress and binds nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		if (Exhaust{Target: Target{Kind: TargetEachFriendlyCreature}, Bind: true}).resolveGate(
			ctx,
		) {
			t.Error("exhausting with no creatures should report no progress")
		}
		if ctx.HasIt {
			t.Error("nothing should be bound when no creature was exhausted")
		}
	})

	t.Run("validate requires a target", func(t *testing.T) {
		if (Exhaust{}).validate() == nil {
			t.Error("unset target should be invalid")
		}
	})
}
