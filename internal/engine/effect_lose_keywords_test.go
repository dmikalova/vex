package engine

import "testing"

// TestLoseKeywordsText covers the printed clause for one and two keywords.
func TestLoseKeywordsText(t *testing.T) {
	one := LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Elusive},
	}
	if got := one.Text(); got != "for the remainder of the turn, it loses elusive" {
		t.Errorf("text = %q", got)
	}
	two := LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Taunt, Elusive},
	}
	if got := two.Text(); got != "for the remainder of the turn, it loses taunt and elusive" {
		t.Errorf("text = %q", got)
	}
}

// TestLoseKeywordsValidate rejects a missing target, an empty keyword list, and
// an unset keyword.
func TestLoseKeywordsValidate(t *testing.T) {
	if err := (LoseKeywords{Keywords: []Keyword{Elusive}}).validate(); err == nil {
		t.Error("a missing target should be rejected")
	}
	if err := (LoseKeywords{Target: Target{Kind: TargetTriggeringCreature}}).validate(); err == nil {
		t.Error("an empty keyword list should be rejected")
	}
	if err := (LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{0},
	}).validate(); err == nil {
		t.Error("an unset keyword should be rejected")
	}
	if err := (LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Taunt, Elusive},
	}).validate(); err != nil {
		t.Errorf("a valid strip should pass: %v", err)
	}
}

// TestLoseKeywordsResolve strips a creature of taunt and elusive for the turn,
// and the ready phase restores them.
func TestLoseKeywordsResolve(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("warded", 3, WithKeywords(Elusive, Taunt)), 0)
	if !g.hasKeyword(id, Elusive) || !g.hasKeyword(id, Taunt) {
		t.Fatal("the creature should start with elusive and taunt")
	}

	LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Taunt, Elusive},
	}.Resolve(&EffectContext{Resolver: g, Controller: 0, It: id, HasIt: true})

	if g.hasKeyword(id, Elusive) {
		t.Error("elusive should be lost for the turn")
	}
	if g.hasKeyword(id, Taunt) {
		t.Error("taunt should be lost for the turn")
	}

	// A second loss of a keyword already lost is a no-op.
	before := len(g.Log)
	g.SetRecording(true)
	g.LoseKeywordFrom(id, Taunt)
	g.SetRecording(false)
	if len(g.Log) != before {
		t.Error("re-losing a keyword already lost should log nothing new")
	}

	g.EndPlayPhase(0)
	if !g.hasKeyword(id, Elusive) || !g.hasKeyword(id, Taunt) {
		t.Error("the lost keywords should return when the turn ends")
	}
}
