package engine

import "testing"

func TestStunEffects(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	stun := Stun{Target: Target{Kind: TargetEachFriendlyCreature}}
	if stun.Text() != "stun each friendly creature" {
		t.Errorf("stun text = %q", stun.Text())
	}
	stun.Resolve(ctx)
	if !g.State.Cards[src].Stunned || !g.State.Cards[friend].Stunned {
		t.Error("stun should stun each friendly creature")
	}
	if g.State.Cards[foe].Stunned {
		t.Error("stun of friendly creatures should not touch the enemy")
	}

	unstun := Unstun{Target: Target{Kind: TargetEachFriendlyCreature}}
	if unstun.Text() != "unstun each friendly creature" {
		t.Errorf("unstun text = %q", unstun.Text())
	}
	unstun.Resolve(ctx)
	if g.State.Cards[src].Stunned || g.State.Cards[friend].Stunned {
		t.Error("unstun should clear the stun on each friendly creature")
	}
}

func TestStunAndNeighbors(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 3), 1)
	mid := g.AddToBattleline(testCreature("mid", 3), 1)
	right := g.AddToBattleline(testCreature("right", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := Stun{Target: Target{Kind: TargetChosenCreature}.AndNeighbors()}
	if e.Text() != "stun a creature and each of its neighbors" {
		t.Errorf("text = %q", e.Text())
	}
	// The default chooser picks the first candidate (left); it stuns left and its
	// only neighbor, mid, but not right.
	e.Resolve(ctx)
	if !g.Stunned(left) || !g.Stunned(mid) {
		t.Error("the chosen creature and its neighbor should be stunned")
	}
	if g.Stunned(right) {
		t.Error("a non-neighbor should not be stunned")
	}
}

// TestNeighborsOfCreatureFought covers the Before Fight pairing: the creature
// being fought is in context, and NeighborsOf reaches past it to its neighbors
// without touching it.
func TestNeighborsOfCreatureFought(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 3), 1)
	mid := g.AddToBattleline(testCreature("mid", 3), 1)
	right := g.AddToBattleline(testCreature("right", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: mid, HasIt: true}

	e := Stun{Target: Target{Kind: TargetCreatureFought}.NeighborsOf()}
	want := "stun each neighbor of the creature " + SelfName + " fights"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	e.Resolve(ctx)
	if !g.Stunned(left) || !g.Stunned(right) {
		t.Error("both neighbors of the creature fought should be stunned")
	}
	if g.Stunned(mid) {
		t.Error("the creature fought should not itself be stunned")
	}
}

func TestExhaust(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := Exhaust{Target: Target{Kind: TargetThisCreature}}
	if e.Text() != "exhaust "+SelfName {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if !g.State.Cards[src].Exhausted {
		t.Error("Exhaust should exhaust the creature")
	}
}

func TestReady(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.State.Cards[src].Exhausted = true
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := Ready{Target: Target{Kind: TargetThisCreature}}
	if e.Text() != "ready "+SelfName {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Cards[src].Exhausted {
		t.Error("Ready should ready the creature")
	}
}
