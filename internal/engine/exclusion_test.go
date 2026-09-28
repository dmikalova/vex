package engine

import "testing"

// TestExclusionNamesItsCard pins which card each exclusion leaves out, which is
// the only thing that separates them: all three print the same "other". Source
// always names the card the ability is printed on, It names the card in context
// and drops nothing without one, and Focus falls back to the source so "another
// creature" can never land on the card it is defined against.
func TestExclusionNamesItsCard(t *testing.T) {
	g := NewGame("A", "B", 1)
	source := g.AddToBattleline(testCreature("source", 3), 0)
	other := g.AddToBattleline(testCreature("other", 3), 0)
	withIt := &EffectContext{Resolver: g, Controller: 0, Source: source, It: other, HasIt: true}
	withoutIt := &EffectContext{Resolver: g, Controller: 0, Source: source}

	cases := []struct {
		exclusion Exclusion
		ctx       *EffectContext
		want      LocalID
		drops     bool
	}{
		{excludeNone, withIt, 0, false},
		{ExcludeSource, withIt, source, true},
		{ExcludeSource, withoutIt, source, true},
		{ExcludeIt, withIt, other, true},
		{ExcludeIt, withoutIt, 0, false},
		{ExcludeFocus, withIt, other, true},
		{ExcludeFocus, withoutIt, source, true},
	}
	for _, c := range cases {
		got, drops := c.exclusion.excluded(c.ctx)
		if drops != c.drops || (drops && got != c.want) {
			t.Errorf("exclusion %d excluded = (%d, %v), want (%d, %v)",
				c.exclusion, got, drops, c.want, c.drops)
		}
		if admits := c.exclusion.admits(c.ctx, c.want); c.drops && admits {
			t.Errorf("exclusion %d admitted the card it excludes (%d)", c.exclusion, c.want)
		}
	}
	for _, exclusion := range allExclusions() {
		if !exclusion.filters() {
			t.Errorf("exclusion %d does not report that it filters", exclusion)
		}
		if got := exclusion.qualifyNoun("card"); got != "other card" {
			t.Errorf("exclusion %d qualifyNoun = %q, want %q", exclusion, got, "other card")
		}
	}
	if excludeNone.filters() {
		t.Error("the zero Exclusion reports that it filters")
	}
	if got := excludeNone.qualifyNoun("card"); got != "card" {
		t.Errorf("the zero Exclusion qualifyNoun = %q, want %q", got, "card")
	}
}
