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
