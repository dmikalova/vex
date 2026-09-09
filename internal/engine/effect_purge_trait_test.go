package engine

import "testing"

func TestPurgeEachOfChosenTrait(t *testing.T) {
	g := NewGame("A", "B", 1)
	d1 := g.AddToBattleline(
		NewCard("d1", Dis, Creature, Common, WithPower(3), WithTraits(Demon)),
		0,
	)
	d2 := g.AddToBattleline(
		NewCard("d2", Dis, Creature, Common, WithPower(3), WithTraits(Demon)),
		0,
	)
	keep := g.AddToBattleline(
		NewCard("keep", Dis, Creature, Common, WithPower(3), WithTraits(Beast)),
		0,
	)
	foe := g.AddToBattleline(
		NewCard("foe", Untamed, Creature, Common, WithPower(3), WithTraits(Demon)),
		1,
	)

	ctx := &EffectContext{Resolver: g, Controller: 0}
	// The default chooser answers ChooseOption with index 0. Traits sort
	// alphabetically (Beast, Demon), so force the Demon pick with a scripted chooser.
	g.SetChooser(0, optionPicker{idx: 1})
	PurgeEachOfChosenTrait{}.Resolve(ctx)

	if g.inPlay(d1) || g.inPlay(d2) || g.inPlay(foe) {
		t.Errorf("Demon creatures should be purged")
	}
	if !g.inPlay(keep) {
		t.Errorf("the Beast creature should survive")
	}
	if g.State.Aember[0] != 2 {
		t.Errorf("controller Æmber = %d, want 2", g.State.Aember[0])
	}
	if g.State.Aember[1] != 1 {
		t.Errorf("opponent Æmber = %d, want 1", g.State.Aember[1])
	}
}

func TestPurgeEachOfChosenTraitTextAndValidate(t *testing.T) {
	e := PurgeEachOfChosenTrait{}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v", err)
	}
	want := "choose a trait, then purge each card with that trait. " +
		"Each player gains 1 Æmber for each card they controlled that was purged this way"
	if got := e.Text(); got != want {
		t.Errorf("text = %q", got)
	}

	// With no traited card in play there is nothing to choose, so Resolve is a
	// no-op (the len(present)==0 early return).
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e.Resolve(ctx)
	if g.State.Aember[0] != 0 || g.State.Aember[1] != 0 {
		t.Error("an empty board should pay no Æmber")
	}
}

// TestPurgeEachOfChosenTraitCaptures covers the payout being captured instead of
// gained when a continuous captor (Ether Spider) watches the paid pool.
func TestPurgeEachOfChosenTraitCaptures(t *testing.T) {
	g := NewGame("A", "B", 1)
	demon := g.AddToBattleline(
		NewCard("demon", Dis, Creature, Common, WithPower(3), WithTraits(Demon)), 0)
	captor := g.AddToBattleline(
		NewCard("spider", Logos, Creature, Common, WithPower(3),
			WithReplaces(Instead{Of: EventAemberAddedToPool, Player: Controller, With: Capture})),
		0)

	PurgeEachOfChosenTrait{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.inPlay(demon) {
		t.Error("the Demon should be purged")
	}
	if got := g.AmberOn(captor); got != 1 {
		t.Errorf("captor Æmber = %d, want 1 (payout captured, not gained)", got)
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("pool = %d, want 0 (captured, not gained)", g.State.Aember[0])
	}
}
