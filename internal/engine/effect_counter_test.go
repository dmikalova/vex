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
	want := "for each damaged creature in play, give {self} 2 +1 power counters"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}

	e.Resolve(ctx)
	if g.Power(c) != 5 {
		t.Errorf("power = %d, want 5 (one damaged creature, two counters)", g.Power(c))
	}

	if got := (AddPowerCounter{Amount: -2}).counters(); got != "2 -1 power counters" {
		t.Errorf("negative counters = %q", got)
	}
}
