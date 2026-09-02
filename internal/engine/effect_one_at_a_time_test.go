package engine

import "testing"

// TestOneAtATimeText covers the rendered phrase and the validation of its bounds.
func TestOneAtATimeText(t *testing.T) {
	e := OneAtATime{
		Times:  3,
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
	if err := (OneAtATime{Times: 3}).validate(); err == nil {
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
		Times:  3,
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
		Times:  2,
		Target: Target{Kind: TargetChosenFriendlyCreature},
		Verbs:  []CreatureVerb{ReadyVerb{}},
	}.Resolve(ctx)

	if !g.State.Cards[a].Exhausted || !g.State.Cards[b].Exhausted {
		t.Error("a declined first pass should ready nobody")
	}
}
