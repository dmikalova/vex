package engine

import "testing"

func TestArchiveSource(t *testing.T) {
	if got := (ArchiveSource{}).Text(); got != "archive "+SelfName {
		t.Errorf("text = %q", got)
	}

	t.Run("archives an in-play source from play", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 3), 0)

		ArchiveSource{}.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})

		if !g.State.Archives[0].contains(src) {
			t.Error("an in-play source should be archived from play")
		}
	})

	t.Run("archives a resolving action instead of discarding it", func(t *testing.T) {
		g := started(t)
		idx := int(g.State.Hand[0].Count)
		id := g.AddToHand(
			NewCard(
				"Self Archive", Brobnar, Tactic, Common,
				WithAbility(TriggerAfterPlay, ArchiveSource{}),
			),
			0,
		)

		if err := g.PlayAction(0, idx); err != nil {
			t.Fatalf("PlayAction: %v", err)
		}

		if !g.State.Archives[0].contains(id) {
			t.Error("a self-archiving action should go to the archives")
		}
		if g.State.Discard[0].contains(id) {
			t.Error("a self-archiving action should not go to the discard pile")
		}
	})
}
