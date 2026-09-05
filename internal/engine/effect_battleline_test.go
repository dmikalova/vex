package engine

import (
	"slices"
	"testing"
)

func TestSwapBattlelinePositions(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 0)
	middle := g.AddToBattleline(testCreature("middle", 2), 0)
	right := g.AddToBattleline(testCreature("right", 2), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 2), 1)

	g.SwapBattlelinePositions(left, right)
	if got, want := g.Battleline(0), []LocalID{right, middle, left}; !slices.Equal(got, want) {
		t.Fatalf("battleline after swap = %v, want %v", got, want)
	}
	if got, want := g.Battleline(1), []LocalID{enemy}; !slices.Equal(got, want) {
		t.Fatalf("enemy battleline after friendly swap = %v, want %v", got, want)
	}

	g.SwapBattlelinePositions(left, enemy)
	if got, want := g.Battleline(0), []LocalID{right, middle, left}; !slices.Equal(got, want) {
		t.Fatalf("battleline after cross-battleline swap = %v, want %v", got, want)
	}
	if got, want := g.Battleline(1), []LocalID{enemy}; !slices.Equal(got, want) {
		t.Fatalf("enemy battleline after cross-battleline swap = %v, want %v", got, want)
	}
}

func TestSwap(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 2), 0)
	other := g.AddToBattleline(testCreature("other", 2), 0)
	right := g.AddToBattleline(testCreature("right", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: host, Controller: 0}
	e := Swap{With: Target{Kind: TargetChosenOtherFriendlyCreature}}

	if err := (Swap{}).validate(); err == nil {
		t.Fatal("an unset target should be rejected")
	}
	if err := validateEffect(e); err != nil {
		t.Fatalf("validate = %v", err)
	}
	if got, want := e.Text(), "swap this creature with another friendly creature in your battleline"; got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}

	e.Resolve(ctx)

	if got, want := g.Battleline(0), []LocalID{other, host, right}; !slices.Equal(got, want) {
		t.Fatalf("battleline after swap = %v, want %v", got, want)
	}
	if !ctx.HasIt || ctx.It != other {
		t.Fatalf("context card = %d (has %v), want swapped creature %d", ctx.It, ctx.HasIt, other)
	}
}

func TestMoveToFlankGameMethod(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 0)
	middle := g.AddToBattleline(testCreature("middle", 2), 0)
	right := g.AddToBattleline(testCreature("right", 2), 0)

	g.MoveToFlank(middle, true)
	if got, want := g.Battleline(0), []LocalID{left, right, middle}; !slices.Equal(got, want) {
		t.Fatalf("battleline after move right = %v, want %v", got, want)
	}

	g.MoveToFlank(middle, false)
	if got, want := g.Battleline(0), []LocalID{middle, left, right}; !slices.Equal(got, want) {
		t.Fatalf("battleline after move left = %v, want %v", got, want)
	}

	// A creature in no battleline leaves every line unchanged.
	g.MoveToFlank(LocalID(200), true)
	if got, want := g.Battleline(0), []LocalID{middle, left, right}; !slices.Equal(got, want) {
		t.Fatalf("battleline after moving an absent creature = %v, want %v", got, want)
	}
}

func TestMoveToFlank(t *testing.T) {
	if err := (MoveToFlank{}).validate(); err == nil {
		t.Fatal("an unset target should be rejected")
	}
	e := MoveToFlank{Target: Target{Kind: TargetTriggeringCreature}}
	if err := validateEffect(e); err != nil {
		t.Fatalf("validate = %v", err)
	}
	if got, want := e.Text(), "move it to either flank of its controller's battleline"; got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 2), 1)
	mover := g.AddToBattleline(testCreature("mover", 2), 1)
	c := g.AddToBattleline(testCreature("c", 2), 1)

	// Default chooser has no preference, so index 0 (the left flank) is taken.
	ctx := &EffectContext{Resolver: g, Controller: 0, It: mover, HasIt: true}
	e.Resolve(ctx)
	if got, want := g.Battleline(1), []LocalID{mover, a, c}; !slices.Equal(got, want) {
		t.Fatalf("battleline after move to left = %v, want %v", got, want)
	}

	// The effect's controller picks the right flank of the target's own line.
	g.SetChooser(0, optionPicker{idx: 1})
	e.Resolve(ctx)
	if got, want := g.Battleline(1), []LocalID{a, c, mover}; !slices.Equal(got, want) {
		t.Fatalf("battleline after move to right = %v, want %v", got, want)
	}

	// A target that selects nothing leaves the battleline unchanged.
	noTarget := MoveToFlank{Target: Target{Kind: TargetTriggeringCreature}}
	noTarget.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if got, want := g.Battleline(1), []LocalID{a, c, mover}; !slices.Equal(got, want) {
		t.Fatalf("battleline after empty target = %v, want %v", got, want)
	}
}
