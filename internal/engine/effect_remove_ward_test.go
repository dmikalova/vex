package engine

import "testing"

func TestRemoveWardText(t *testing.T) {
	e := RemoveWard{Target: Target{Kind: TargetChosenCreature}}
	if got := e.Text(); got != "remove a ward from a creature" {
		t.Errorf("remove ward text = %q", got)
	}
}

func TestRemoveWardValidate(t *testing.T) {
	if (RemoveWard{}).validate() == nil {
		t.Error("RemoveWard with no target should fail validation")
	}
	if (RemoveWard{Target: Target{Kind: TargetChosenCreature}}).validate() != nil {
		t.Error("RemoveWard with a target should validate")
	}
}

// TestRemoveWard: the ward is taken off the chosen creature, and the removal is
// narrated.
func TestRemoveWard(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	warded := g.AddToBattleline(testCreature("warded", 3), 1)
	g.SetWarded(warded, true)

	g.SetChooser(0, &idQueueChooser{ids: []LocalID{warded}})
	entries := len(g.Log)
	RemoveWard{Target: Target{Kind: TargetChosenCreature}}.
		Resolve(&EffectContext{
			Resolver:   g,
			Source:     src,
			Controller: 0,
		})

	if g.Warded(warded) {
		t.Error("the chosen creature should lose its ward")
	}
	if len(g.Log) == entries {
		t.Error("removing a ward should be narrated")
	}
}

// TestRemoveWardUnwardedTarget: choosing a creature that carries no ward leaves it
// unchanged and narrates that there was no ward to remove.
func TestRemoveWardUnwardedTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	plain := g.AddToBattleline(testCreature("plain", 3), 1)

	g.SetChooser(0, &idQueueChooser{ids: []LocalID{plain}})
	entries := len(g.Log)
	RemoveWard{Target: Target{Kind: TargetChosenCreature}}.
		Resolve(&EffectContext{
			Resolver:   g,
			Source:     src,
			Controller: 0,
		})

	if g.Warded(plain) {
		t.Error("an unwarded creature should stay unwarded")
	}
	if len(g.Log) == entries {
		t.Error("choosing an unwarded creature should still be narrated")
	}
}
