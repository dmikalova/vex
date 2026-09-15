package engine

import "testing"

// DestroyChosen destroys any number of creatures the controller picks from its
// Target pool and tallies them into Produced.Destroyed.
func TestDestroyChosen(t *testing.T) {
	if got := (DestroyChosen{Target: Target{Kind: TargetEachFriendlyCreature}}).Text(); got != "destroy any number of friendly creatures" {
		t.Errorf("text = %q", got)
	}
	if (DestroyChosen{}).validate() == nil {
		t.Error("a targetless DestroyChosen should not validate")
	}
	if (DestroyChosen{Target: Target{Kind: TargetEachFriendlyCreature}}).validate() != nil {
		t.Error("a DestroyChosen with a target should validate")
	}
	if got := (DestroyChosen{Target: Target{Kind: TargetEachFriendlyCreature}, Amount: 2}).Text(); got != "destroy 2 friendly creatures" {
		t.Errorf("fixed-amount text = %q", got)
	}
	if (DestroyChosen{Target: Target{Kind: TargetEachFriendlyCreature}, Amount: -1}).validate() == nil {
		t.Error("a negative Amount should not validate")
	}

	t.Run("a fixed Amount destroys exactly that many", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ids := []LocalID{
			g.AddToBattleline(testCreature("a", 1), 0),
			g.AddToBattleline(testCreature("b", 1), 0),
			g.AddToBattleline(testCreature("c", 1), 0),
		}
		ctx := &EffectContext{Resolver: g, Controller: 0}

		DestroyChosen{Target: Target{Kind: TargetEachFriendlyCreature}, Amount: 2}.Resolve(ctx)

		alive := 0
		for _, id := range ids {
			if g.inPlay(id) {
				alive++
			}
		}
		if alive != 1 {
			t.Errorf("alive = %d, want 1 (2 of 3 destroyed)", alive)
		}
		if ctx.Produced.Destroyed[0] != 2 {
			t.Errorf("Destroyed = %v, want [2 0]", ctx.Produced.Destroyed)
		}
	})

	t.Run("destroys every pick and tallies them", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ids := []LocalID{
			g.AddToBattleline(testCreature("a", 1), 0),
			g.AddToBattleline(testCreature("b", 1), 0),
			g.AddToBattleline(testCreature("c", 1), 0),
		}
		ctx := &EffectContext{Resolver: g, Controller: 0}

		DestroyChosen{Target: Target{Kind: TargetEachFriendlyCreature}}.Resolve(ctx)

		for _, id := range ids {
			if g.inPlay(id) {
				t.Errorf("%s should have been destroyed", g.Name(id))
			}
		}
		if ctx.Produced.Destroyed[0] != 3 {
			t.Errorf("Destroyed = %v, want [3 0]", ctx.Produced.Destroyed)
		}
	})

	t.Run("declining destroys nobody", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		id := g.AddToBattleline(testCreature("a", 1), 0)
		g.SetChooser(0, &cardDecliner{decline: true})
		ctx := &EffectContext{Resolver: g, Controller: 0}

		DestroyChosen{Target: Target{Kind: TargetEachFriendlyCreature}}.Resolve(ctx)

		if !g.inPlay(id) {
			t.Error("a declined DestroyChosen should destroy nobody")
		}
		if ctx.Produced.Destroyed[0] != 0 {
			t.Errorf("Destroyed = %v, want [0 0]", ctx.Produced.Destroyed)
		}
	})
}
