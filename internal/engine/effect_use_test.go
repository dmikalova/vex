package engine

import "testing"

func TestUseTextAndValidation(t *testing.T) {
	pool := Target{Kind: TargetEachFriendlyInPlay}.OfHouse(Mars).Other()
	if got := (Use{Max: 2, Target: pool}).Text(); got != "use 2 other Mars cards, one at a time" {
		t.Errorf("text = %q", got)
	}
	if err := (Use{Max: 2}).validate(); err == nil {
		t.Error("unset target should be rejected")
	}
	if err := (Use{Target: pool}).validate(); err == nil {
		t.Error("non-positive max should be rejected")
	}
	if err := (Use{Max: 2, Target: pool}).validate(); err != nil {
		t.Errorf("valid effect should pass, got %v", err)
	}
}

func TestUseUsesCreaturesSequentially(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("source", Mars, Artifact, Uncommon), 0)
	first := g.AddToBattleline(NewCard("first", Mars, Creature, Common, WithPower(3)), 0)
	second := g.AddToBattleline(NewCard("second", Mars, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(NewCard("off-house", Logos, Creature, Common, WithPower(3)), 0)
	exhausted := g.AddToBattleline(NewCard("exhausted", Mars, Creature, Common, WithPower(3)), 0)
	g.State.Cards[exhausted].Exhausted = true
	g.SetChooser(0, idChooser{id: second})
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	Use{Max: 2, Target: Target{Kind: TargetEachFriendlyInPlay}.OfHouse(Mars).Other()}.Resolve(ctx)

	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2", g.Aember(0))
	}
	if !g.Exhausted(first) || !g.Exhausted(second) {
		t.Error("both ready Mars creatures should have been used")
	}
	if g.Exhausted(exhausted) != true {
		t.Error("the already-exhausted Mars creature should remain exhausted")
	}
}

func TestUseUsesArtifactAction(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("source", Mars, Artifact, Uncommon), 0)
	actionArtifact := g.AddArtifact(NewCard("action artifact", Mars, Artifact, Common,
		WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 3})), 0)
	g.AddArtifact(NewCard("blank artifact", Mars, Artifact, Common), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	Use{Max: 1, Target: Target{Kind: TargetEachFriendlyInPlay}.OfHouse(Mars).Other()}.Resolve(ctx)

	if g.Aember(0) != 3 {
		t.Errorf("aember = %d, want 3", g.Aember(0))
	}
	if !g.Exhausted(actionArtifact) {
		t.Error("the artifact action should exhaust the artifact")
	}
}

func TestUseStopsWhenNoChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("source", Mars, Artifact, Uncommon), 0)
	c := g.AddToBattleline(NewCard("c", Mars, Creature, Common, WithPower(3)), 0)
	other := g.AddToBattleline(NewCard("other", Mars, Creature, Common, WithPower(3)), 0)
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	Use{Max: 1, Target: Target{Kind: TargetEachFriendlyInPlay}.OfHouse(Mars).Other()}.Resolve(ctx)

	if g.Exhausted(c) || g.Exhausted(other) {
		t.Error("a rejected choice should not use the card")
	}
}

func TestUseStopsWhenNoneUsable(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("source", Mars, Artifact, Uncommon), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	Use{Max: 1, Target: Target{Kind: TargetEachFriendlyInPlay}.OfHouse(Mars).Other()}.Resolve(ctx)

	if g.TimesUsedThisTurn(src) != 0 {
		t.Errorf("Use must not use its own source; used %d times", g.TimesUsedThisTurn(src))
	}
}

func TestSentenceEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	e := Sentence{Effect: GainAember{Player: Controller, Amount: 1}}

	if got := e.Text(); got != "gain 1 Æmber." {
		t.Errorf("text = %q", got)
	}
	e.Resolve(ctx)
	if g.Aember(0) != 1 {
		t.Errorf("aember = %d, want 1", g.Aember(0))
	}
	if err := (Sentence{Effect: Destroy{}}).validate(); err == nil {
		t.Error("sentence should surface child validation")
	}

	seq := Sequence{Effects: []Effect{
		Sentence{Effect: Destroy{Target: Target{Kind: TargetThisCreature}}},
		Use{Max: 2, Target: Target{Kind: TargetEachFriendlyInPlay}.OfHouse(Mars).Other()},
	}}
	if got := seq.Text(); got != "destroy "+SelfName+". Use 2 other Mars cards, one at a time" {
		t.Errorf("sequence text = %q", got)
	}
}
