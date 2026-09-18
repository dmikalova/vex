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

// TestAddPowerCounterChosenBindsContext covers Animator's first step: a chosen
// target leaves the picked card in context (ctx.It) so a following effect can act
// on the same card the controller just chose.
func TestAddPowerCounterChosenBindsContext(t *testing.T) {
	g := NewGame("A", "B", 1)
	art := g.AddArtifact(testArtifact("art"), 0)
	g.SetChooser(0, FirstChooser{})
	ctx := &EffectContext{Resolver: g, Source: art, Controller: 0}

	AddPowerCounter{Target: Target{Kind: TargetChosenArtifact}, Amount: 3}.Resolve(ctx)

	if !ctx.HasIt || ctx.It != art {
		t.Fatalf("chosen counter target should be left in context: HasIt=%v It=%v",
			ctx.HasIt, ctx.It)
	}
	if got := int(g.State.Cards[art].PowerCounters); got != 3 {
		t.Fatalf("power counters = %d, want 3", got)
	}
}

// TestAddPowerCounterEqual covers Mimic Gel: the counter count is set equal to a
// live count (a chosen creature's power) rather than a fixed Amount.
func TestAddPowerCounterEqual(t *testing.T) {
	g := NewGame("A", "B", 1)
	mimic := g.AddToBattleline(testCreature("mimic", 0), 0)
	chosen := g.AddToBattleline(testCreature("chosen", 6), 0)
	ctx := &EffectContext{Resolver: g, Source: mimic, Controller: 0, It: chosen, HasIt: true}

	e := AddPowerCounter{Target: Target{Kind: TargetThisCreature}, Equal: PowerOfChosen{}}
	if got := e.Text(); got != "give {self} +1 power counters equal to its power" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(ctx)
	if g.Power(mimic) != 6 {
		t.Errorf(
			"power = %d, want 6 (counters equal to the chosen creature's power)",
			g.Power(mimic),
		)
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
		Per:    CardsInPlay{Player: EachPlayer, Type: Creature, Damaged: true},
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

// TestAddPowerCounterWalk covers Growth Surge: counters placed along an inward
// flank walk instead of on a named Target.
func TestAddPowerCounterWalk(t *testing.T) {
	e := AddPowerCounter{Walk: []int{3, 2, 1}}
	want := "choose a flank creature. Give it three +1 power counters, " +
		"its neighbor two +1 power counters, and the neighbor's other neighbor " +
		"a +1 power counter"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	// A walk chooses its own creatures, so it needs no Target to validate.
	if err := e.validate(); err != nil {
		t.Errorf("walk validate = %v, want nil", err)
	}

	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 4), 0)
	mid := g.AddToBattleline(testCreature("mid", 4), 0)
	right := g.AddToBattleline(testCreature("right", 4), 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Power(left) != 7 || g.Power(mid) != 6 || g.Power(right) != 5 {
		t.Errorf(
			"power = %d/%d/%d, want 7/6/5",
			g.Power(left), g.Power(mid), g.Power(right),
		)
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
