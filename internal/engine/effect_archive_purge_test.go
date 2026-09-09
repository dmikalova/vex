package engine

import "testing"

func TestArchivePurgedCardText(t *testing.T) {
	if got := (ArchivePurgedCard{}).Text(); got != "archive a purged card you own" {
		t.Errorf("text = %q", got)
	}
}

func TestArchivePurgedCard(t *testing.T) {
	g := NewGame("A", "B", 1)
	purged := g.Register(testCreature("p", 1), 0)
	g.State.Purge[0].add(purged)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// A sole purged card is auto-chosen and moved into archives.
	(ArchivePurgedCard{}).Resolve(ctx)
	if len(g.Purge(0)) != 0 {
		t.Errorf("purge pile = %v, want empty", g.Purge(0))
	}
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != purged {
		t.Errorf("archives = %v, want [%d]", g.State.Archives[0].slice(), purged)
	}
}

func TestArchivePurgedCardEmpty(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// An empty purge pile archives nothing.
	(ArchivePurgedCard{}).Resolve(ctx)
	if g.State.Archives[0].Count != 0 {
		t.Errorf("archives = %v, want empty", g.State.Archives[0].slice())
	}
}
