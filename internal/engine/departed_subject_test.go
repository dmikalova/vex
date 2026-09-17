package engine

import "testing"

// TestDepartedSubjectPowerAfterDestroy is the Power of Fire regression: an effect
// that destroys a creature and then reads its power (LoseAember equal to its
// power) must read the power the creature had the instant before it left play —
// including power counters — not the zeroed core it leaves behind.
func TestDepartedSubjectPowerAfterDestroy(t *testing.T) {
	g := started(t)
	victim := g.AddToBattleline(testCreature("victim", 4), 0)
	g.AddPowerCounter(victim, 2) // power 6
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// While the creature is still in play its live power is authoritative, even
	// after a capture has been recorded.
	captureDepartingSubject(ctx, victim)
	if got := ctx.powerOf(victim); got != 6 {
		t.Fatalf("in-play powerOf = %d, want 6 (live)", got)
	}

	Destroy{Target: Target{Kind: TargetChosenCreature}}.destroy(ctx, []LocalID{victim})
	if resolverInPlay(ctx, victim) {
		t.Fatal("victim should have left play")
	}
	if got := ctx.powerOf(victim); got != 6 {
		t.Errorf("departed powerOf = %d, want 6 (last-known incl. counter)", got)
	}
	if got := (PowerOfChosen{}).Value(ctx); got != 6 {
		t.Errorf("PowerOfChosen after destroy = %d, want 6 (was zeroed core)", got)
	}
}

// TestDepartedSubjectAmberAndDamageAfterDestroy covers the Æmber-on-card and
// damage readers: a source that destroys itself and then reads its own on-card
// Æmber or damage sees the last-known amounts, not the zeroed core.
func TestDepartedSubjectAmberAndDamageAfterDestroy(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.AddAmberOn(src, 2)
	g.State.Cards[src].Damage = 1
	ctx := &EffectContext{Resolver: g, Controller: 0, Source: src}

	Destroy{Target: Target{Kind: TargetChosenCreature}}.destroy(ctx, []LocalID{src})
	if resolverInPlay(ctx, src) {
		t.Fatal("src should have left play")
	}
	if got := ctx.amberOn(src); got != 2 {
		t.Errorf("departed amberOn = %d, want 2", got)
	}
	if got := ctx.damageOn(src); got != 1 {
		t.Errorf("departed damageOn = %d, want 1", got)
	}
	if got := (AemberOnThis{}).Value(ctx); got != 2 {
		t.Errorf("AemberOnThis after destroy = %d, want 2", got)
	}
	if got := (DamageOnThis{}).Value(ctx); got != 1 {
		t.Errorf("DamageOnThis after destroy = %d, want 1", got)
	}
}
