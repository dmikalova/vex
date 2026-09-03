package engine

import "testing"

func TestLoseArmorValidatesItsTarget(t *testing.T) {
	if err := (LoseArmor{}).validate(); err == nil {
		t.Error("an untargeted LoseArmor should be rejected")
	}
	e := LoseArmor{Target: Target{Kind: TargetEachEnemyCreature}.WithArmor()}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if want := "each enemy creature with armor loses all of its armor"; e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
}

func TestLoseArmorStripsAndTallies(t *testing.T) {
	g := started(t)
	plated := g.AddToBattleline(
		NewCard("Plated", Dis, Creature, Common, WithPower(5), WithArmor(2)),
		1,
	)
	bare := g.AddToBattleline(testCreature("Bare", 5), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	target := Target{Kind: TargetEachEnemyCreature}.WithArmor()

	// Only the armored creature is a target at all, so only it is stripped.
	if got := target.Select(ctx); len(got) != 1 || got[0] != plated {
		t.Errorf("selected %v, want just the armored creature %d (not %d)", got, plated, bare)
	}
	LoseArmor{Target: target}.Resolve(ctx)

	if got := g.State.Cards[plated].ArmorRemaining; got != 0 {
		t.Errorf("armor remaining = %d, want 0", got)
	}
	if got := g.ArmorStripped(plated); got != 2 {
		t.Errorf("armor stripped = %d, want 2", got)
	}
	if got := ArmorLostThisWay.perTargetValue(ctx, plated); got != 2 {
		t.Errorf("per-target value = %d, want 2", got)
	}
	if want := "point of armor it lost this way"; ArmorLostThisWay.perTargetText() != want {
		t.Errorf("per-target text = %q, want %q", ArmorLostThisWay.perTargetText(), want)
	}

	// A second strip adds to the tally rather than replacing it.
	g.State.Cards[plated].ArmorRemaining = 1
	g.StripArmor(plated)
	if got := g.ArmorStripped(plated); got != 3 {
		t.Errorf("armor stripped = %d, want 3 after a second strip", got)
	}

	// The strip is turn-scoped: readying refreshes the armor and clears the tally.
	g.readyPhase(1)
	if got := g.State.Cards[plated].ArmorRemaining; got != 2 {
		t.Errorf("armor remaining = %d, want the full 2 back after readying", got)
	}
	if got := g.ArmorStripped(plated); got != 0 {
		t.Errorf("armor stripped = %d, want 0 after readying", got)
	}
}

func TestStripArmorSkipsACardOutOfPlay(t *testing.T) {
	g := started(t)
	gone := g.Register(NewCard("Gone", Dis, Creature, Common, WithPower(5), WithArmor(2)), 0)
	g.State.Discard[0].add(gone)

	g.StripArmor(gone) // must not write in-play state onto a discarded card

	if got := g.State.Cards[gone].ArmorStripped; got != 0 {
		t.Errorf("armor stripped = %d, want 0 for a card out of play", got)
	}
}
