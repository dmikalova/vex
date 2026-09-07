package engine

import "testing"

func TestTakeControlTargeted(t *testing.T) {
	if got := (TakeControl{Duration: UntilThisLeavesPlay}).Text(); got != "take control of this creature until "+UpgradeName+" leaves play" {
		t.Errorf("host text = %q", got)
	}
	tgt := TakeControl{
		Target:   Target{Kind: TargetChosenEnemyCreature}.OnFlank(),
		Duration: UntilThisLeavesPlay,
	}
	if got := tgt.Text(); got != "take control of an enemy flank creature until "+SelfName+" leaves play" {
		t.Errorf("targeted text = %q", got)
	}
	if (TakeControl{Target: Target{Kind: TargetChosenEnemyCreature}, Duration: EndOfTurn}).validate() == nil {
		t.Error("only UntilThisLeavesPlay should be valid")
	}

	g := NewGame("A", "B", 1)
	harland := g.AddToBattleline(testCreature("harland", 1), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	tgt.Resolve(&EffectContext{Resolver: g, Source: harland, Controller: 0})
	if g.controller(foe) != 0 {
		t.Errorf("controller of the seized creature = %d, want 0", g.controller(foe))
	}
}
