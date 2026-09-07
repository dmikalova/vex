package engine

import "testing"

func TestAddPowerCounter(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: c, Controller: 0}

	e := AddPowerCounter{Target: Target{Kind: TargetThisCreature}, Amount: 1}
	if e.Text() != "give {self} a +1 power counter" {
		t.Errorf("text = %q", e.Text())
	}

	e.Resolve(ctx)
	if g.Power(c) != 4 {
		t.Errorf("power = %d, want 4 (+1 counter)", g.Power(c))
	}
	// Counters stack.
	e.Resolve(ctx)
	if g.Power(c) != 5 {
		t.Errorf("power = %d, want 5 (two +1 counters)", g.Power(c))
	}
}

// TestAddPowerCounterPer covers Martian Hounds: several counters at once, scaled
// by a board count, with the target chosen only once.
func TestAddPowerCounterPer(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	hurt := g.AddToBattleline(testCreature("hurt", 5), 1)
	g.DealDamage(0, []DamageTarget{{ID: hurt, Amount: 1}})
	ctx := &EffectContext{Resolver: g, Source: c, Controller: 0}

	e := AddPowerCounter{
		Target: Target{Kind: TargetThisCreature},
		Amount: 2,
		Per:    InPlay{Player: EachPlayer, Type: Creature, Damaged: true},
	}
	want := "for each damaged creature in play, give {self} two +1 power counters"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}

	e.Resolve(ctx)
	if g.Power(c) != 5 {
		t.Errorf("power = %d, want 5 (one damaged creature, two counters)", g.Power(c))
	}

	if got := (AddPowerCounter{Amount: -2}).counters(); got != "two -1 power counters" {
		t.Errorf("negative counters = %q", got)
	}
	// A count larger than KeyForge ever prints falls back to digits.
	if got := (AddPowerCounter{Amount: 11}).counters(); got != "11 +1 power counters" {
		t.Errorf("large counters = %q", got)
	}
}

// A -1 power counter that lowers a damaged creature's power to its damage
// destroys it at the resolution boundary, the same sweep a leaving buff triggers
// — CanUse's map order must never leave a lethal creature sitting in play.
func TestAddPowerCounterSettlesLethal(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 3), 1)
	g.DealDamage(0, []DamageTarget{{ID: c, Amount: 2}})
	if !g.inPlay(c) {
		t.Fatal("2 damage should not destroy a 3-power creature")
	}

	AddPowerCounter{Target: Target{Kind: TargetThisCreature}, Amount: -1}.
		Resolve(&EffectContext{Resolver: g, Source: c, Controller: 0})
	g.settleDestroyed(0) // the resolution boundary settles the counter (ADR 0029)

	if g.inPlay(c) {
		t.Errorf("a -1 counter dropping power to 2 with 2 damage should destroy it")
	}
}
