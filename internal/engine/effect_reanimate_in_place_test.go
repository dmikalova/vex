package engine

import (
	"slices"
	"testing"
)

func TestReanimateTopOfDeckInPlaceText(t *testing.T) {
	want := "discard the top card of your deck. If it is a creature, after " +
		SelfName + " leaves play, put that creature into play in " + SelfName +
		"'s position in the battleline"
	if got := (ReanimateTopOfDeckInPlace{}).Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

func TestReanimateTopOfDeckInPlaceResolve(t *testing.T) {
	t.Run("arms a reanimation when the top card is a creature", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 5), 0)
		top := g.AddToDeck(testCreature("reborn", 3), 0)

		ReanimateTopOfDeckInPlace{}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

		if !slices.Contains(g.Discard(0), top) {
			t.Error("the top card should be discarded")
		}
		if g.State.ReanimationsCount != 1 {
			t.Fatalf("ReanimationsCount = %d, want 1", g.State.ReanimationsCount)
		}
		r := g.State.Reanimations[0]
		if r != (ReanimateInPlace{Source: src, Creature: top, Controller: 0, Position: 0}) {
			t.Errorf("armed record = %+v", r)
		}
	})

	t.Run("captures the source's battleline position", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToBattleline(testCreature("left", 4), 0)
		src := g.AddToBattleline(testCreature("src", 5), 0) // index 1
		g.AddToDeck(testCreature("reborn", 3), 0)

		ReanimateTopOfDeckInPlace{}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

		if g.State.Reanimations[0].Position != 1 {
			t.Errorf("position = %d, want 1", g.State.Reanimations[0].Position)
		}
	})

	t.Run("does not arm when the top card is not a creature", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 5), 0)
		top := g.AddToDeck(NewCard("relic", Mars, Artifact, Common), 0)

		ReanimateTopOfDeckInPlace{}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

		if g.State.ReanimationsCount != 0 {
			t.Errorf("ReanimationsCount = %d, want 0", g.State.ReanimationsCount)
		}
		if !slices.Contains(g.Discard(0), top) {
			t.Error("a non-creature top card is still discarded")
		}
	})

	t.Run("does nothing with an empty deck", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 5), 0)

		ReanimateTopOfDeckInPlace{}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

		if g.State.ReanimationsCount != 0 {
			t.Errorf("ReanimationsCount = %d, want 0", g.State.ReanimationsCount)
		}
	})
}

func TestFireReanimations(t *testing.T) {
	t.Run("puts the creature into play at the source's former slot", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		left := g.AddToBattleline(testCreature("left", 4), 0)
		src := g.AddToBattleline(testCreature("src", 5), 0)     // index 1
		right := g.AddToBattleline(testCreature("right", 4), 0) // index 2
		reborn := g.AddToDiscard(testCreature("reborn", 3), 0)
		g.armReanimateInPlace(
			ReanimateInPlace{Source: src, Creature: reborn, Controller: 0, Position: 1},
		)

		g.State.Battleline[0].remove(src) // the source has left play

		g.fireReanimations()

		if want := []LocalID{left, reborn, right}; !slices.Equal(g.Battleline(0), want) {
			t.Errorf("battleline = %v, want %v", g.Battleline(0), want)
		}
		if !g.State.Cards[reborn].Exhausted {
			t.Error("a reanimated creature enters play exhausted")
		}
		if g.State.ReanimationsCount != 0 {
			t.Errorf(
				"ReanimationsCount = %d, want 0 (fired record dropped)",
				g.State.ReanimationsCount,
			)
		}
	})

	t.Run("keeps a record whose source is still in play", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 5), 0)
		reborn := g.AddToDiscard(testCreature("reborn", 3), 0)
		g.armReanimateInPlace(
			ReanimateInPlace{Source: src, Creature: reborn, Controller: 0, Position: 0},
		)

		g.fireReanimations()

		if g.State.ReanimationsCount != 1 {
			t.Errorf("ReanimationsCount = %d, want 1 (record kept)", g.State.ReanimationsCount)
		}
		if slices.Contains(g.Battleline(0), reborn) {
			t.Error("nothing should enter play while the source is in play")
		}
	})
}

func TestPutIntoPlayAtClampsPosition(t *testing.T) {
	g := NewGame("A", "B", 1)
	other := g.AddToBattleline(testCreature("other", 4), 0)
	reborn := g.AddToDiscard(testCreature("reborn", 3), 0)

	g.putIntoPlayAt(reborn, 0, 3) // position 3 past a 1-creature line clamps to the right flank

	if want := []LocalID{other, reborn}; !slices.Equal(g.Battleline(0), want) {
		t.Errorf("battleline = %v, want %v", g.Battleline(0), want)
	}
}

func TestArmReanimateInPlaceDropsWhenFull(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i < maxReanimations; i++ {
		g.armReanimateInPlace(ReanimateInPlace{Source: LocalID(i + 1)})
	}
	g.armReanimateInPlace(ReanimateInPlace{Source: 99})

	if g.State.ReanimationsCount != maxReanimations {
		t.Errorf(
			"ReanimationsCount = %d, want %d (extra arm dropped)",
			g.State.ReanimationsCount,
			maxReanimations,
		)
	}
}
