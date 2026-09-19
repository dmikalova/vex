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
	entries := len(g.Log)
	// A stun that finds its target already stunned still logs the choice, just
	// without a state change.
	stun.Resolve(ctx)
	if len(g.Log) == entries {
		t.Error("re-stunning an already-stunned creature should still log the choice")
	}

	unstun := Unstun{Target: Target{Kind: TargetEachFriendlyCreature}}
	if unstun.Text() != "unstun each friendly creature" {
		t.Errorf("unstun text = %q", unstun.Text())
	}
	unstun.Resolve(ctx)
	if g.State.Cards[src].Stunned || g.State.Cards[friend].Stunned {
		t.Error("unstun should clear the stun on each friendly creature")
	}
	if got := g.Log[len(g.Log)-1].Text(g); got != "A unstuns src and friend" {
		t.Errorf("unstun log = %q, want %q", got, "A unstuns src and friend")
	}
	// An unstun that frees no one — nothing was stunned — records nothing.
	entries = len(g.Log)
	unstun.Resolve(ctx)
	if len(g.Log) != entries {
		t.Error("unstunning creatures that are not stunned should record nothing")
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
	// The effect renders the past; fightTense puts a Before Fight: ability's line
	// in the present (TestBeforeFightTargetReadsInPresentTense).
	want := "stun each neighbor of the creature " + SelfName + " fought"
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

// TestNeighborsOfFoughtCreatureThatLeftPlay covers the snapshot fallback: when the
// fight destroyed the fought creature, "each neighbor of the fought creature"
// reads the neighbors the fight snapshotted rather than the empty live line, so the
// effect still lands (Smite; see TestSmiteHitsNeighborsWhenFoughtCreatureDies).
func TestNeighborsOfFoughtCreatureThatLeftPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 3), 1)
	mid := g.Register(testCreature("mid", 3), 1)
	g.State.Discard[1].add(mid) // the fought creature died in the fight
	right := g.AddToBattleline(testCreature("right", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: mid, HasIt: true}
	ctx.Produced.Neighbors = []LocalID{left, right}

	Stun{Target: Target{Kind: TargetTheFoughtCreature}.NeighborsOf()}.Resolve(ctx)
	if !g.Stunned(left) || !g.Stunned(right) {
		t.Error("a departed fought creature's snapshotted neighbors should still be reached")
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

// A "you may stun a creature" is one clickable creature, so May drives it by the
// click rather than by a Yes/No, and an empty board is not worth asking about.
func TestStunDeclinable(t *testing.T) {
	chosen := Stun{Target: Target{Kind: TargetChosenCreature}}
	if !chosen.declinable() {
		t.Error("a chosen Stun should be declinable")
	}
	if (Stun{Target: Target{Kind: TargetEachCreature}}).declinable() {
		t.Error("an untargeted Stun should not be declinable")
	}

	empty := NewGame("A", "B", 1)
	if !chosen.vacuous(&EffectContext{Resolver: empty, Controller: 0}) {
		t.Error("a Stun with no creature to stun should be vacuous")
	}

	taken := NewGame("A", "B", 1)
	taken.SetChooser(0, &cardDecliner{})
	foe := taken.AddToBattleline(testCreature("foe", 3), 1)
	if !chosen.resolveOptional(&EffectContext{Resolver: taken, Controller: 0}) {
		t.Error("clicking the creature should report the stun resolved")
	}
	if !taken.Stunned(foe) {
		t.Error("the clicked creature should be stunned")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	spared := declined.AddToBattleline(testCreature("spared", 3), 1)
	if chosen.resolveOptional(&EffectContext{Resolver: declined, Controller: 0}) {
		t.Error("declining should report nothing resolved")
	}
	if declined.Stunned(spared) {
		t.Error("a declined Stun should stun nothing")
	}
}
