package engine

import "testing"

func TestCaptureReactionOnFight(t *testing.T) {
	if actCapture.describe() != "capture Æmber" {
		t.Errorf("describe = %q", actCapture.describe())
	}

	g := started(t)
	g.State.Aember[1] = 3

	// Register the lasting the way Take Hostages authors it, to exercise the
	// CaptureAember -> actCapture mapping and validation.
	ForRemainderOfTurn{On: EventFight, Do: CaptureAember{Amount: 1, Target: Target{Kind: TargetTriggeringCreature}, Source: Opponent}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	att := g.AddToBattleline(NewCard("att", Brobnar, Creature, Common, WithPower(4)), 0)
	def := g.AddToBattleline(testCreature("def", 2), 1)
	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.AmberOn(att) != 1 {
		t.Errorf("captured Æmber on fighter = %d, want 1", g.AmberOn(att))
	}
	if g.State.Aember[1] != 2 {
		t.Errorf("opponent pool = %d, want 2 after a capture", g.State.Aember[1])
	}
}
