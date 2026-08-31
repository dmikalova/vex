package engine

import "testing"

func TestReadyIfFirstUse(t *testing.T) {
	e := ReadyIfFirstUse{Target: Target{Kind: TargetThisCreature}}
	if got := e.Text(); got != "if this is the first time {self} was used this turn, ready it" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.State.Cards[src].Exhausted = true
	g.State.Cards[src].TimesUsedThisTurn = 1
	e.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
	if g.Exhausted(src) {
		t.Error("first use should ready the creature")
	}

	g.State.Cards[src].Exhausted = true
	g.State.Cards[src].TimesUsedThisTurn = 2
	e.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
	if !g.Exhausted(src) {
		t.Error("later use should not ready the creature")
	}
}

func TestReadyCreatures(t *testing.T) {
	// Text: leads with the Max count and renders the bare singular noun.
	e := ReadyCreatures{
		Max:    CardsRevealed{},
		Target: Target{Kind: TargetEachFriendlyCreature}.OfHouse(Mars),
	}
	if got := e.Text(); got != "for each card revealed this way, ready a Mars creature" {
		t.Errorf("text = %q", got)
	}
	// Text with no Max renders a single ready with no "for each" clause.
	if got := (ReadyCreatures{Target: Target{Kind: TargetEachFriendlyCreature}}).Text(); got != "ready a creature" {
		t.Errorf("no-max text = %q", got)
	}

	// Resolve readies up to Max exhausted friendly Mars creatures.
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(NewCard("a", Mars, Creature, Common, WithPower(3)), 0)
	b := g.AddToBattleline(NewCard("b", Mars, Creature, Common, WithPower(3)), 0)
	g.State.Cards[a].Exhausted = true
	g.State.Cards[b].Exhausted = true
	ctx := &EffectContext{Resolver: g, Controller: 0, Produced: Produced{Revealed: 2}}
	e.Resolve(ctx)
	if g.Exhausted(a) || g.Exhausted(b) {
		t.Errorf("both should be readied: a=%v b=%v", g.Exhausted(a), g.Exhausted(b))
	}

	// Max nil readies exactly one.
	g2 := NewGame("A", "B", 1)
	c := g2.AddToBattleline(NewCard("c", Mars, Creature, Common, WithPower(3)), 0)
	d := g2.AddToBattleline(NewCard("d", Mars, Creature, Common, WithPower(3)), 0)
	g2.State.Cards[c].Exhausted = true
	g2.State.Cards[d].Exhausted = true
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	(ReadyCreatures{Target: Target{Kind: TargetEachFriendlyCreature}}).Resolve(ctx2)
	if g2.Exhausted(c) == g2.Exhausted(d) {
		t.Errorf("exactly one should be readied: c=%v d=%v", g2.Exhausted(c), g2.Exhausted(d))
	}

	// No exhausted candidate: nothing happens.
	g3 := NewGame("A", "B", 1)
	g3.AddToBattleline(NewCard("e", Mars, Creature, Common, WithPower(3)), 0) // ready
	ctx3 := &EffectContext{Resolver: g3, Controller: 0, Produced: Produced{Revealed: 1}}
	e.Resolve(ctx3) // no panic, no-op

	// Declining the choice readies nothing (two candidates so the chooser is asked).
	g4 := NewGame("A", "B", 1)
	f := g4.AddToBattleline(NewCard("f", Mars, Creature, Common, WithPower(3)), 0)
	h := g4.AddToBattleline(NewCard("h", Mars, Creature, Common, WithPower(3)), 0)
	g4.State.Cards[f].Exhausted = true
	g4.State.Cards[h].Exhausted = true
	g4.SetChooser(0, orderRejectChooser{}) // ChooseCreature declines
	ctx4 := &EffectContext{Resolver: g4, Controller: 0, Produced: Produced{Revealed: 1}}
	e.Resolve(ctx4)
	if !g4.Exhausted(f) || !g4.Exhausted(h) {
		t.Error("declining should ready nothing")
	}
}

func TestSingularNoun(t *testing.T) {
	cases := map[string]string{
		"each friendly Mars creature":  "Mars creature",
		"each enemy creature":          "creature",
		"each other friendly creature": "creature",
		"each creature":                "creature",
		"a friendly creature":          "creature",
		"an enemy creature":            "creature",
		"a creature":                   "creature",
		"an artifact":                  "artifact",
		"it":                           "it", // no known prefix, unchanged
	}
	for in, want := range cases {
		if got := singularNoun(in); got != want {
			t.Errorf("singularNoun(%q) = %q, want %q", in, got, want)
		}
	}
}
