package engine

import "testing"

func TestEnrageEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	foe1 := g.AddToBattleline(testCreature("foe1", 3), 1)
	foe2 := g.AddToBattleline(testCreature("foe2", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := Enrage{Target: Target{Kind: TargetEachEnemyCreature}}
	if e.Text() != "enrage each enemy creature" {
		t.Errorf("enrage text = %q", e.Text())
	}
	e.Resolve(ctx)
	if !g.Enraged(foe1) || !g.Enraged(foe2) {
		t.Error("enrage should enrage each enemy creature")
	}
	if g.Enraged(src) {
		t.Error("enraging enemy creatures should not touch a friendly creature")
	}

	// An enrage that finds its target already enraged still logs the choice, just
	// without a state change.
	entries := len(g.Log)
	e.Resolve(ctx)
	if len(g.Log) == entries {
		t.Error("re-enraging an already-enraged creature should still log the choice")
	}
}

func TestEnrageValidate(t *testing.T) {
	if (Enrage{}).validate() == nil {
		t.Error("Enrage with no target should fail validation")
	}
	if (Enrage{Target: Target{Kind: TargetEachEnemyCreature}}).validate() != nil {
		t.Error("Enrage with a target should validate")
	}
}
