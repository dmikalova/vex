package engine

import "testing"

func TestMoveWardText(t *testing.T) {
	e := MoveWard{
		From: Target{Kind: TargetChosenCreature},
		Onto: Target{Kind: TargetChosenOtherCreature},
	}
	if got := e.Text(); got != "move a ward from a creature to another creature" {
		t.Errorf("move ward text = %q", got)
	}
}

func TestMoveWardValidate(t *testing.T) {
	if (MoveWard{Onto: Target{Kind: TargetChosenOtherCreature}}).validate() == nil {
		t.Error("MoveWard with no source should fail validation")
	}
	if (MoveWard{From: Target{Kind: TargetChosenCreature}}).validate() == nil {
		t.Error("MoveWard with no destination should fail validation")
	}
	if (MoveWard{
		From: Target{Kind: TargetChosenCreature},
		Onto: Target{Kind: TargetChosenOtherCreature},
	}).validate() != nil {
		t.Error("MoveWard with both targets should validate")
	}
}

// TestMoveWard: a ward is taken off the one warded source and placed on the
// destination, and the move is narrated.
func TestMoveWard(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	from := g.AddToBattleline(testCreature("from", 3), 0)
	onto := g.AddToBattleline(testCreature("onto", 3), 1)
	g.SetWarded(from, true)

	e := MoveWard{
		From: Target{Kind: TargetChosenCreature},
		Onto: Target{Kind: TargetChosenOtherCreature},
	}
	// Only the warded creature is offered as the source; the destination is chosen.
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{onto}})
	entries := len(g.Log)
	e.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

	if g.Warded(from) {
		t.Error("the source should lose its ward")
	}
	if !g.Warded(onto) {
		t.Error("the destination should gain the ward")
	}
	if len(g.Log) == entries {
		t.Error("moving a ward should be narrated")
	}
}

// TestMoveWardNoWardedSource: with no warded creature to draw from, the effect
// does nothing.
func TestMoveWardNoWardedSource(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	onto := g.AddToBattleline(testCreature("onto", 3), 1)

	MoveWard{
		From: Target{Kind: TargetChosenCreature},
		Onto: Target{Kind: TargetChosenOtherCreature},
	}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

	if g.Warded(onto) {
		t.Error("with no warded source, no ward should move")
	}
}
