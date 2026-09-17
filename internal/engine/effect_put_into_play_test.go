package engine

import (
	"slices"
	"testing"
)

func TestPutIntoPlay(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		if got := (PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}, Control: ControlYours}).Text(); got != "put it into play under your control" {
			t.Errorf("text = %q", got)
		}
		if got := (PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}}).Text(); got != "put it into play" {
			t.Errorf("text = %q", got)
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (PutIntoPlay{}).validate() == nil {
			t.Error("unset target should be invalid")
		}
		if (PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}}).validate() != nil {
			t.Error("a set target should be valid")
		}
	})

	t.Run(
		"puts a discarded creature into play under your control, keeping ownership",
		func(t *testing.T) {
			g := NewGame("A", "B", 1)
			foe := g.AddToDiscard(testCreature("foe", 3), 1)
			ctx := &EffectContext{Resolver: g, Controller: 0, It: foe, HasIt: true}

			PutIntoPlay{
				Target:  Target{Kind: TargetTriggeringCreature},
				Control: ControlYours,
			}.Resolve(
				ctx,
			)

			if !slices.Contains(g.Battleline(0), foe) {
				t.Error("the creature should enter P0's battleline")
			}
			if g.controller(foe) != 0 {
				t.Errorf("controller = %d, want 0", g.controller(foe))
			}
			if g.owner(foe) != 1 {
				t.Errorf("owner = %d, want 1 (unchanged)", g.owner(foe))
			}
		},
	)

	t.Run("puts a card into play under its owner's control by default", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		c := g.AddToHand(testCreature("c", 3), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0, It: c, HasIt: true}

		PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}}.Resolve(ctx)

		if !slices.Contains(g.Battleline(1), c) {
			t.Error("the creature should enter its owner's battleline")
		}
		if g.controller(c) != 1 {
			t.Errorf("controller = %d, want 1", g.controller(c))
		}
	})

	t.Run("puts an artifact into play in the controller's artifact row", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		art := g.AddToDiscard(NewCard("relic", Mars, Artifact, Common), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0, It: art, HasIt: true}

		PutIntoPlay{
			Target:  Target{Kind: TargetTriggeringCreature},
			Control: ControlYours,
		}.Resolve(
			ctx,
		)

		if !slices.Contains(g.Artifacts(0), art) {
			t.Error("the artifact should enter P0's artifact row")
		}
	})

	t.Run("a card already in play is not moved", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		c := g.AddToBattleline(testCreature("c", 3), 1)
		g.putIntoPlay(c, 0)
		if !slices.Contains(g.Battleline(1), c) {
			t.Error("an in-play card should stay where it is")
		}
	})

	t.Run("a gigantic half is left where it came from, never put into play", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		half := g.AddToDiscard(
			NewCard("titan", Logos, Creature, Common, WithPower(9), WithGiganticRole(GiganticBase)),
			1,
		)
		g.putIntoPlay(half, 0)
		if g.inPlay(half) {
			t.Error("a lone gigantic half must not enter play")
		}
		if !g.State.Discard[1].contains(half) {
			t.Error("the half should stay in the discard pile it came from")
		}
	})

	t.Run("with no target, does nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		PutIntoPlay{
			Target: Target{Kind: TargetTriggeringCreature},
		}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	})
}

func TestEachPlayerPutsHandCreaturesIntoPlay(t *testing.T) {
	const readyText = "each player reveals their hand and puts each creature from their hand into play ready"
	const plainText = "each player reveals their hand and puts each creature from their hand into play"

	t.Run("text", func(t *testing.T) {
		if got := (EachPlayerPutsHandCreaturesIntoPlay{Ready: true}).Text(); got != readyText {
			t.Errorf("text = %q", got)
		}
		if got := (EachPlayerPutsHandCreaturesIntoPlay{}).Text(); got != plainText {
			t.Errorf("text = %q", got)
		}
	})

	t.Run(
		"both players put their hand creatures into play ready, leaving non-creatures",
		func(t *testing.T) {
			g := NewGame("A", "B", 1)
			mine := g.AddToHand(testCreature("mine", 3), 0)
			myAction := g.AddToHand(NewCard("tac", Mars, Tactic, Common), 0)
			theirs := g.AddToHand(testCreature("theirs", 4), 1)

			EachPlayerPutsHandCreaturesIntoPlay{Ready: true}.Resolve(
				&EffectContext{Resolver: g, Controller: 0},
			)

			if !slices.Contains(g.Battleline(0), mine) {
				t.Error("P0's hand creature should enter P0's battleline")
			}
			if !slices.Contains(g.Battleline(1), theirs) {
				t.Error("P1's hand creature should enter P1's battleline")
			}
			if g.Exhausted(mine) || g.Exhausted(theirs) {
				t.Error("creatures put into play ready should not be exhausted")
			}
			if !slices.Contains(g.Hand(0), myAction) {
				t.Error("a non-creature should stay in hand")
			}
		},
	)

	t.Run("without Ready, the creatures enter exhausted", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		mine := g.AddToHand(testCreature("mine", 3), 0)

		EachPlayerPutsHandCreaturesIntoPlay{}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)

		if !g.Exhausted(mine) {
			t.Error("a creature put into play without Ready should be exhausted")
		}
	})

	t.Run("the active player may have the opponent resolve first", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetChooser(0, optionPicker{idx: 1}) // opponent first
		mine := g.AddToHand(testCreature("mine", 3), 0)
		theirs := g.AddToHand(testCreature("theirs", 4), 1)

		EachPlayerPutsHandCreaturesIntoPlay{Ready: true}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)

		if !slices.Contains(g.Battleline(0), mine) {
			t.Error("P0's hand creature should still enter play")
		}
		if !slices.Contains(g.Battleline(1), theirs) {
			t.Error("P1's hand creature should still enter play")
		}
	})
}

// TestPutIntoPlayReady covers the Ready flag on PutIntoPlay: its text gains the
// "ready" clause, and Resolve leaves the entered creature ready rather than
// exhausted.
func TestPutIntoPlayReady(t *testing.T) {
	if got := (PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}}).Text(); got !=
		"put it into play" {
		t.Errorf("plain text = %q", got)
	}
	ready := PutIntoPlay{Target: Target{Kind: TargetTriggeringCreature}, Ready: true}
	if got := ready.Text(); got != "put it into play ready" {
		t.Errorf("ready text = %q", got)
	}

	g := NewGame("A", "B", 1)
	creature := g.Register(NewCard("saur", Saurian, Creature, Common, WithPower(2)), 0)
	g.State.Discard[0].add(creature)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: creature, HasIt: true}

	ready.Resolve(ctx)

	if !g.InPlay(creature) {
		t.Fatal("the creature should have been put into play")
	}
	if g.State.Cards[creature].Exhausted {
		t.Error("PutIntoPlay with Ready should leave the creature ready")
	}
}
