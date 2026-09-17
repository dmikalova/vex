package engine

import "testing"

// These tests cover COPYING PRINTED STATS — the mechanic behind Cyber-Clone: a
// creature takes on another card's printed power (as an override) and gains its
// printed armor, keywords, and traits until the copier leaves play.

// TestCopiedStatsSource reports the copied source, or none before a copy.
func TestCopiedStatsSource(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(testCreature("src", 6), 0)
	recipient := g.AddToBattleline(testCreature("recipient", 2), 0)

	if _, ok := g.copiedStatsSource(recipient); ok {
		t.Error("no copy should report no source")
	}
	g.CopyStats(recipient, src)
	if got, ok := g.copiedStatsSource(recipient); !ok || got != src {
		t.Errorf("copied source = %d,%v, want %d,true", got, ok, src)
	}
}

// TestCopyStatsGone copies nothing when the recipient has left play.
func TestCopyStatsGone(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(testCreature("src", 6), 0)
	gone := g.AddToBattleline(testCreature("gone", 2), 0)
	g.DestroyEach(0, []LocalID{gone})
	g.CopyStats(gone, src) // must not panic and must record nothing
	if g.State.Cards[gone].CopiedStatsSourcePlus != 0 {
		t.Error("copying onto a creature no longer in play should do nothing")
	}
}

// TestCopyPrintedStatsValidate rejects a missing recipient or source.
func TestCopyPrintedStatsValidate(t *testing.T) {
	if (CopyPrintedStats{Source: Target{Kind: TargetTheSameCreature}}).validate() == nil {
		t.Error("an unset target should be invalid")
	}
	if (CopyPrintedStats{Target: Target{Kind: TargetThisCreature}}).validate() == nil {
		t.Error("an unset source should be invalid")
	}
	if (CopyPrintedStats{
		Target: Target{Kind: TargetThisCreature},
		Source: Target{Kind: TargetTheSameCreature},
	}).validate() != nil {
		t.Error("a set target and source should be valid")
	}
}

// TestCopyPrintedStatsText renders the printed clause.
func TestCopyPrintedStatsText(t *testing.T) {
	e := CopyPrintedStats{
		Target: Target{Kind: TargetThisCreature},
		Source: Target{Kind: TargetTheSameCreature},
	}
	want := SelfName + " has power equal to the same creature's printed power" +
		" and gains its printed armor, keywords, and traits"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// TestCreatureCopiedStatsText names the copier and the card whose stats it copied.
func TestCreatureCopiedStatsText(t *testing.T) {
	e := CreatureCopiedStats{Creature: 4, Source: 7}
	if got := e.Text(stubNamer{}); got != "Card4 copies the printed stats of Card7" {
		t.Errorf("text = %q", got)
	}
}

// TestCopyPrintedStatsResolve overrides the recipient's power with the source's
// printed power and folds in the source's printed armor, keywords, and traits.
func TestCopyPrintedStatsResolve(t *testing.T) {
	g := started(t)
	g.SetRecording(true)
	source := g.AddToBattleline(
		testCreature(
			"source",
			6,
			WithArmor(3),
			WithTraits(Beast, Mutant),
			WithKeywords(Skirmish),
		),
		0,
	)
	// The recipient keeps its own armor (2) and its own trait (Mutant): armor and
	// traits are gained, so the copy adds to them rather than replacing them.
	recipient := g.AddToBattleline(
		testCreature("recipient", 2, WithArmor(2), WithTraits(Mutant)),
		0,
	)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     recipient,
		It:         source,
		HasIt:      true,
	}

	CopyPrintedStats{
		Target: Target{Kind: TargetThisCreature},
		Source: Target{Kind: TargetTheSameCreature},
	}.Resolve(ctx)

	if got := g.Power(recipient); got != 6 {
		t.Errorf("power = %d, want 6 (override to the source's printed power)", got)
	}
	if got := g.Armor(recipient); got != 5 {
		t.Errorf("armor = %d, want 5 (own 2 plus copied 3)", got)
	}
	if !g.hasKeyword(recipient, Skirmish) {
		t.Error("the recipient should gain the source's keyword")
	}
	if g.hasKeyword(recipient, Poison) {
		t.Error("the recipient should not gain a keyword the source lacks")
	}
	if !g.HasTrait(recipient, Beast) {
		t.Error("the recipient should gain the source's trait")
	}
	if !g.HasTrait(recipient, Mutant) {
		t.Error("the recipient should keep its own trait")
	}
	if g.HasTrait(recipient, Human) {
		t.Error("the recipient should not gain a trait the source lacks")
	}
	// Own {Mutant} plus copied {Beast}; the copied Mutant is not double-counted.
	if got := g.TraitCount(recipient); got != 2 {
		t.Errorf("trait count = %d, want 2 (Mutant plus copied Beast)", got)
	}
	if !hasLogLine(g, "recipient copies the printed stats of source") {
		t.Error("the copy should be logged")
	}
}

// TestCopyPrintedStatsResolveNoSource copies nothing when the source selects none.
func TestCopyPrintedStatsResolveNoSource(t *testing.T) {
	g := started(t)
	recipient := g.AddToBattleline(testCreature("recipient", 2), 0)
	CopyPrintedStats{
		Target: Target{Kind: TargetThisCreature},
		Source: Target{Kind: TargetTheSameCreature},
	}.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: recipient})
	if g.State.Cards[recipient].CopiedStatsSourcePlus != 0 {
		t.Error("a source that selects nothing should copy no stats")
	}
}

// TestTraitCountNoCopy counts only the printed traits of a creature not copying.
func TestTraitCountNoCopy(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 2, WithTraits(Beast, Mutant)), 0)
	if got := g.TraitCount(c); got != 2 {
		t.Errorf("trait count = %d, want 2", got)
	}
}
