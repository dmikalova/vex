package engine

import "testing"

func TestPurgeArchivesText(t *testing.T) {
	if got := (PurgeArchives{}).Text(); got != "purge any number of cards from your archives" {
		t.Errorf("Text() = %q", got)
	}
	if (PurgeArchives{}).validate() != nil {
		t.Error("PurgeArchives should validate")
	}
}

func TestPurgeArchivesResolve(t *testing.T) {
	g := started(t)
	one := g.AddToArchives(NewCard("Archived One", Logos, Creature, Common), 0)
	two := g.AddToArchives(NewCard("Archived Two", Logos, Creature, Common), 0)
	g.SetChooser(0, &declineAfterChooser{ids: []LocalID{one, two}})

	ctx := &EffectContext{Resolver: g, Controller: 0}
	PurgeArchives{}.Resolve(ctx)

	if got := (CardsPurged{}).Value(ctx); got != 2 {
		t.Errorf("purged tally = %d, want 2", got)
	}
	if g.State.Archives[0].contains(one) || g.State.Archives[0].contains(two) {
		t.Error("both archived cards should have been purged")
	}
	if !g.State.Purge[0].contains(one) || !g.State.Purge[0].contains(two) {
		t.Error("both purged cards should be in the purge pile")
	}
}

func TestPurgeArchivesPurgingNone(t *testing.T) {
	g := started(t)
	g.AddToArchives(NewCard("Archived", Logos, Creature, Common), 0)
	g.SetChooser(0, &declineAfterChooser{}) // decline immediately

	ctx := &EffectContext{Resolver: g, Controller: 0}
	PurgeArchives{}.Resolve(ctx)

	if got := (CardsPurged{}).Value(ctx); got != 0 {
		t.Errorf("purged tally = %d, want 0 when nothing is purged", got)
	}
}

func TestPurgeArchivesEmptyArchives(t *testing.T) {
	g := started(t)

	ctx := &EffectContext{Resolver: g, Controller: 0}
	PurgeArchives{}.Resolve(ctx)

	if got := (CardsPurged{}).Value(ctx); got != 0 {
		t.Errorf("purged tally = %d, want 0 with empty archives", got)
	}
}
