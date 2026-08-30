package engine

import "testing"

func TestPurge(t *testing.T) {
	// Text variants.
	if got := (Purge{Zone: Discard, Count: 2, UpTo: true}).Text(); got != "purge up to 2 cards from a discard pile" {
		t.Errorf("up-to text = %q", got)
	}
	if got := (Purge{Zone: Discard, Type: Creature}).Text(); got != "purge a creature from a discard pile" {
		t.Errorf("single text = %q", got)
	}
	if got := (Purge{Zone: Discard, Count: 2}).Text(); got != "purge 2 cards from a discard pile" {
		t.Errorf("count text = %q", got)
	}

	// Standalone up-to purge: one zone, purge 2 of 3 (default chooser takes 0).
	g := NewGame("A", "B", 1)
	a := g.Register(testCreature("a", 1), 0)
	b := g.Register(testCreature("b", 1), 0)
	c := g.Register(testCreature("c", 1), 0)
	g.State.Discard[0].add(a)
	g.State.Discard[0].add(b)
	g.State.Discard[0].add(c)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !(Purge{Zone: Discard, Count: 2, UpTo: true}).resolveGate(ctx) {
		t.Error("purging cards should report success")
	}
	if len(g.Purge(0)) != 2 || len(g.Discard(0)) != 1 {
		t.Errorf("purged=%d discard=%d, want 2/1", len(g.Purge(0)), len(g.Discard(0)))
	}

	// Two zones: the controller picks one (default: their own), which empties
	// before Count is reached.
	g2 := NewGame("A", "B", 1)
	m := g2.Register(testCreature("m", 1), 0)
	n := g2.Register(testCreature("n", 1), 1)
	g2.State.Discard[0].add(m)
	g2.State.Discard[1].add(n)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	Purge{Zone: Discard, Count: 2, UpTo: true}.Resolve(ctx2)
	if got := g2.Purge(0); len(got) != 1 || got[0] != m {
		t.Errorf("own zone purge = %v, want [m]", got)
	}
	if len(g2.Discard(1)) != 1 {
		t.Error("the opponent's zone should be untouched")
	}

	// Declining ("Done") purges nothing and reports failure.
	g3 := NewGame("A", "B", 1)
	x := g3.Register(testCreature("x", 1), 0)
	g3.State.Discard[0].add(x)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	g3.SetChooser(0, optionPicker{idx: 1}) // options [x, Done] -> idx 1 is Done
	if (Purge{Zone: Discard, Count: 2, UpTo: true}).resolveGate(ctx3) {
		t.Error("declining should report no purge")
	}
	if len(g3.Purge(0)) != 0 {
		t.Error("declining should purge nothing")
	}

	// Type filter: a creature is purged and a non-creature is left; reports success.
	g4 := NewGame("A", "B", 1)
	crea := g4.Register(testCreature("crea", 3), 1)
	act := g4.Register(NewCard("act", Dis, Tactic, Common), 1)
	g4.State.Discard[1].add(act)
	g4.State.Discard[1].add(crea)
	ctx4 := &EffectContext{Resolver: g4, Controller: 0}
	if !(Purge{Zone: Discard, Type: Creature}).resolveGate(ctx4) {
		t.Error("purging a creature should report success")
	}
	if got := g4.Purge(1); len(got) != 1 || got[0] != crea {
		t.Errorf("purge = %v, want [crea]", got)
	}
	if len(g4.Discard(1)) != 1 {
		t.Error("the non-creature should be left in the discard")
	}

	// No matching card: nothing is purged and it reports failure.
	g5 := NewGame("A", "B", 1)
	g5.State.Discard[0].add(g5.Register(NewCard("act2", Dis, Tactic, Common), 0))
	ctx5 := &EffectContext{Resolver: g5, Controller: 0}
	if (Purge{Zone: Discard, Type: Creature}).resolveGate(ctx5) {
		t.Error("no creature to purge should report failure")
	}
}

func TestPurgeFromHand(t *testing.T) {
	// validate rejects an unset player.
	if err := (PurgeFromHand{}).validate(); err == nil {
		t.Error("unset player should fail validation")
	}
	if err := (PurgeFromHand{Player: Opponent}).validate(); err != nil {
		t.Errorf("valid player should pass validation: %v", err)
	}

	// Text and noun variants.
	if got := (PurgeFromHand{Player: Opponent, House: Sanctum}).Text(); got != "you may purge a Sanctum card from your opponent's hand" {
		t.Errorf("house text = %q", got)
	}
	if got := (PurgeFromHand{Player: Controller}).Text(); got != "you may purge a card from your hand" {
		t.Errorf("any-card text = %q", got)
	}

	// Resolve: only the Sanctum card is eligible; the default chooser purges it.
	g := NewGame("A", "B", 1)
	sanctum := g.Register(NewCard("holy", Sanctum, Creature, Common, WithPower(3)), 1)
	other := g.Register(NewCard("dark", Shadows, Creature, Common, WithPower(3)), 1)
	g.State.Hand[1].add(sanctum)
	g.State.Hand[1].add(other)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	PurgeFromHand{Player: Opponent, House: Sanctum}.Resolve(ctx)
	if got := g.Purge(1); len(got) != 1 || got[0] != sanctum {
		t.Errorf("purge = %v, want [sanctum]", got)
	}
	if len(g.Hand(1)) != 1 || g.Hand(1)[0] != other {
		t.Errorf("hand = %v, want [other] (non-Sanctum left)", g.Hand(1))
	}

	// Declining ("Done") purges nothing.
	g2 := NewGame("A", "B", 1)
	holy := g2.Register(NewCard("holy", Sanctum, Creature, Common, WithPower(3)), 1)
	g2.State.Hand[1].add(holy)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	g2.SetChooser(0, optionPicker{idx: 1}) // options [holy, Done] -> idx 1 is Done
	PurgeFromHand{Player: Opponent, House: Sanctum}.Resolve(ctx2)
	if len(g2.Purge(1)) != 0 || len(g2.Hand(1)) != 1 {
		t.Error("declining should purge nothing")
	}

	// No matching card: nothing happens.
	g3 := NewGame("A", "B", 1)
	g3.State.Hand[1].add(g3.Register(NewCard("dark", Shadows, Creature, Common, WithPower(3)), 1))
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	PurgeFromHand{Player: Opponent, House: Sanctum}.Resolve(ctx3)
	if len(g3.Purge(1)) != 0 {
		t.Error("no matching card should purge nothing")
	}
}
