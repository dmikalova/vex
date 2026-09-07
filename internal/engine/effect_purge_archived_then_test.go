package engine

import "testing"

// TestPurgeArchivedCardThenText covers the printed clause with a Stun follow-up.
func TestPurgeArchivedCardThenText(t *testing.T) {
	e := PurgeArchivedCardThen{Then: Stun{Target: Target{Kind: TargetChosenCreature}}}
	if got := e.Text(); got != "you may purge a card from your archives to stun a creature" {
		t.Errorf("text = %q", got)
	}
}

// TestPurgeArchivedCardThenValidate rejects an invalid Then.
func TestPurgeArchivedCardThenValidate(t *testing.T) {
	if err := (PurgeArchivedCardThen{Then: Stun{}}).validate(); err == nil {
		t.Error("an invalid Then should be rejected")
	}
	if err := (PurgeArchivedCardThen{
		Then: Stun{Target: Target{Kind: TargetChosenCreature}},
	}).validate(); err != nil {
		t.Errorf("a valid Then should pass: %v", err)
	}
}

// TestPurgeArchivedCardThenResolve purges a chosen archived card and resolves the
// follow-up (stunning the sole creature).
func TestPurgeArchivedCardThenResolve(t *testing.T) {
	g := started(t)
	archived := g.AddToArchives(testCreature("stored", 3), 0)
	victim := g.AddToBattleline(testCreature("victim", 3), 1)
	g.SetChooser(0, &declineAfterChooser{ids: []LocalID{archived}})

	PurgeArchivedCardThen{Then: Stun{Target: Target{Kind: TargetChosenCreature}}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if !g.State.Purge[0].contains(archived) {
		t.Error("the chosen card should have been purged")
	}
	if !g.State.Cards[victim].Stunned {
		t.Error("the follow-up should have stunned the creature")
	}
}

// TestPurgeArchivedCardThenDeclined does nothing when the controller declines.
func TestPurgeArchivedCardThenDeclined(t *testing.T) {
	g := started(t)
	archived := g.AddToArchives(testCreature("stored", 3), 0)
	victim := g.AddToBattleline(testCreature("victim", 3), 1)
	g.SetChooser(0, &declineAfterChooser{})

	PurgeArchivedCardThen{Then: Stun{Target: Target{Kind: TargetChosenCreature}}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.Purge[0].contains(archived) {
		t.Error("nothing should be purged when the controller declines")
	}
	if g.State.Cards[victim].Stunned {
		t.Error("the follow-up should not run when nothing is purged")
	}
}

// TestPurgeArchivedCardThenEmptyArchives does nothing with no archived cards.
func TestPurgeArchivedCardThenEmptyArchives(t *testing.T) {
	g := started(t)
	victim := g.AddToBattleline(testCreature("victim", 3), 1)

	PurgeArchivedCardThen{Then: Stun{Target: Target{Kind: TargetChosenCreature}}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.Cards[victim].Stunned {
		t.Error("the follow-up should not run with an empty archive")
	}
}
