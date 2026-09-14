package engine

import "testing"

// TestOneAtATimeText covers the rendered phrase and the validation of its bounds.
func TestOneAtATimeText(t *testing.T) {
	e := OneAtATime{
		Times:  Fixed(3),
		Target: Target{Kind: TargetChosenFriendlyCreature},
		Verbs:  []CreatureVerb{ReadyVerb{}, FightVerb{}},
	}
	want := "ready and fight with up to 3 different friendly creatures, one at a time"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}

	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (OneAtATime{Times: Fixed(3)}).validate(); err == nil {
		t.Error("a targetless OneAtATime should not validate")
	}
	if err := (OneAtATime{Target: e.Target}).validate(); err == nil {
		t.Error("a OneAtATime with no passes should not validate")
	}
}

// TestOneAtATimeActsOnDifferentCreatures checks each pass picks a creature no
// earlier pass took, and that it stops once the pool runs dry.
func TestOneAtATimeActsOnDifferentCreatures(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := []LocalID{
		g.AddToBattleline(testCreature("a", 1), 0),
		g.AddToBattleline(testCreature("b", 1), 0),
	}
	for _, id := range mine {
		g.State.Cards[id].Exhausted = true
	}
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Three passes over only two creatures: both are readied, then the third pass
	// finds nobody left and stops.
	OneAtATime{
		Times:  Fixed(3),
		Target: Target{Kind: TargetChosenFriendlyCreature},
		Verbs:  []CreatureVerb{ReadyVerb{}},
	}.Resolve(ctx)

	for _, id := range mine {
		if g.State.Cards[id].Exhausted {
			t.Errorf("%s should have been readied", g.Name(id))
		}
	}
}

// TestOneAtATimeStopsWhenDeclined checks a declined pass ends the whole effect,
// leaving the untouched creatures alone.
func TestOneAtATimeStopsWhenDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 1), 0)
	b := g.AddToBattleline(testCreature("b", 1), 0)
	g.State.Cards[a].Exhausted = true
	g.State.Cards[b].Exhausted = true
	g.SetChooser(0, &cardDecliner{decline: true})
	ctx := &EffectContext{Resolver: g, Controller: 0}

	OneAtATime{
		Times:  Fixed(2),
		Target: Target{Kind: TargetChosenFriendlyCreature},
		Verbs:  []CreatureVerb{ReadyVerb{}},
	}.Resolve(ctx)

	if !g.State.Cards[a].Exhausted || !g.State.Cards[b].Exhausted {
		t.Error("a declined first pass should ready nobody")
	}
}

// TestOneAtATimeEachSetText covers the whole-set mode's rendered phrase and that
// it needs no Times bound to validate — Ghosthawk reaps with each neighbor.
func TestOneAtATimeEachSetText(t *testing.T) {
	e := OneAtATime{
		Target: Target{Kind: TargetEachNeighbor},
		Verbs:  []CreatureVerb{ReapVerb{}},
	}
	want := "reap with each of " + SelfName + "'s neighbors, one at a time"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil (whole-set mode needs no Times)", err)
	}
}

// TestOneAtATimeEachSetActsOnEveryone covers the whole-set mode acting on each
// creature the target names — both of the source's neighbors reap.
func TestOneAtATimeEachSetActsOnEveryone(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 1), 0)
	mid := g.AddToBattleline(testCreature("mid", 1), 0)
	right := g.AddToBattleline(testCreature("right", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, Source: mid}

	OneAtATime{
		Target: Target{Kind: TargetEachNeighbor},
		Verbs:  []CreatureVerb{ReapVerb{}},
	}.Resolve(ctx)

	if !g.Exhausted(left) || !g.Exhausted(right) {
		t.Error("both neighbors should have reaped")
	}
	if got := g.Aember(0); got != 2 {
		t.Errorf("Æmber = %d, want 2 from two reaps", got)
	}
}

// TestOneAtATimeEachSetStopsWhenDeclined covers the decline path: refusing the
// first pick ends the whole-set effect, so neither neighbor reaps.
func TestOneAtATimeEachSetStopsWhenDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 1), 0)
	mid := g.AddToBattleline(testCreature("mid", 1), 0)
	right := g.AddToBattleline(testCreature("right", 1), 0)
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{Resolver: g, Controller: 0, Source: mid}

	OneAtATime{
		Target: Target{Kind: TargetEachNeighbor},
		Verbs:  []CreatureVerb{ReapVerb{}},
	}.Resolve(ctx)

	if g.Exhausted(left) || g.Exhausted(right) {
		t.Error("a declined first pass should reap nobody")
	}
}

// TestOneAtATimeEachSetStopsWhenPoolLeaves covers the guard for a named creature
// that leaves play mid-effect: the first neighbor's reap destroys the other, so
// the loop finds nobody left to act on and stops without forcing a dead creature.
func TestOneAtATimeEachSetStopsWhenPoolLeaves(t *testing.T) {
	g := NewGame("A", "B", 1)
	wipe := NewCard("wipe", Brobnar, Creature, Common, WithPower(1),
		WithAbility(TriggerAfterReap,
			Destroy{Target: Target{Kind: TargetEachOtherFriendlyCreature}}))
	left := g.AddToBattleline(wipe, 0)
	mid := g.AddToBattleline(testCreature("mid", 1), 0)
	right := g.AddToBattleline(testCreature("right", 1), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{left}})
	ctx := &EffectContext{Resolver: g, Controller: 0, Source: mid}

	OneAtATime{
		Target: Target{Kind: TargetEachNeighbor},
		Verbs:  []CreatureVerb{ReapVerb{}},
	}.Resolve(ctx)

	if g.inPlay(right) {
		t.Error("the right neighbor should have been destroyed by the first reap")
	}
	if !g.Exhausted(left) {
		t.Error("the first neighbor should have reaped")
	}
}
