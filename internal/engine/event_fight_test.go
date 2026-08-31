package engine

import "testing"

func TestFightFiresLasting(t *testing.T) {
	g := started(t)
	g.AddLasting(EventFight, actGainAember, 0, 1)
	att := g.AddToBattleline(NewCard("att", Brobnar, Creature, Common, WithPower(4)), 0)
	def := g.AddToBattleline(testCreature("def", 2), 1)

	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.State.Aember[0] != 1 {
		t.Errorf("controller Æmber = %d, want 1 after a Warsong-style fight reaction", g.State.Aember[0])
	}
}

func TestEventFightIsReaction(t *testing.T) {
	if !EventFight.isReaction() {
		t.Error("EventFight should be a reaction point")
	}
	if EventFight.clause() != "each time a friendly creature fights" {
		t.Errorf("clause = %q", EventFight.clause())
	}
}
