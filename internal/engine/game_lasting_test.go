package engine

import "testing"

func TestEmitLastingOrders(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	// Two reactions fire on the same event; the controller orders them.
	g.AddLasting(EventCreaturePlayed, actGainAember, 0, 1)
	g.AddLasting(EventCreaturePlayed, actDealDamage, 0, 2)
	g.SetChooser(0, optionPicker{idx: 1}) // resolve "deal damage" first

	g.AddToHand(testCreature("minion", 4), 0)
	before := g.Aember(0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "minion"), false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if got := g.Aember(0) - before; got != 1 {
		t.Errorf("Æmber gained = %d, want 1", got)
	}
	if g.Damage(foe) != 2 {
		t.Errorf("foe damage = %d, want 2", g.Damage(foe))
	}
}

func TestLastingOnceExpiresAtEndOfTurn(t *testing.T) {
	g := started(t)
	g.AddLastingOnce(EventCreaturePlayed, actReadyPlayed, 0, 0, Mars)
	if g.State.LastingCount != 1 {
		t.Fatalf("setup: lasting count = %d, want 1", g.State.LastingCount)
	}

	g.EndTurn(0)

	if g.State.LastingCount != 0 {
		t.Errorf(
			"a one-shot reaction that never fired should be cleared at end of turn, count = %d",
			g.State.LastingCount,
		)
	}
}

func TestClearLastingKeepsOtherPlayer(t *testing.T) {
	g := started(t)
	g.AddLasting(EventCreaturePlayed, actGainAember, 0, 1)
	g.AddLasting(EventCreaturePlayed, actGainAember, 1, 1) // the opponent's reaction

	g.clearLasting(0)

	if g.State.LastingCount != 1 {
		t.Fatalf("lasting count = %d, want 1 (the opponent's kept)", g.State.LastingCount)
	}
	if g.State.Lasting[0].Controller != 1 {
		t.Errorf("kept controller = %d, want 1", g.State.Lasting[0].Controller)
	}
}

func TestAddLastingCap(t *testing.T) {
	g := started(t)
	for i := 0; i < maxLasting+3; i++ {
		g.AddLasting(EventCreaturePlayed, actGainAember, 0, 1)
	}
	if int(g.State.LastingCount) != maxLasting {
		t.Errorf("lasting count = %d, want %d (capped)", g.State.LastingCount, maxLasting)
	}
}

func TestLastingOnceReadiesMatchingHouseAndSelfRemoves(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(2)), 0)
	sanc := g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(2)), 0)
	g.State.Cards[mars].Exhausted = true
	g.State.Cards[sanc].Exhausted = true

	g.AddLastingOnce(EventCreaturePlayed, actReadyPlayed, 0, 0, Mars)

	// A non-Mars subject is filtered out: not readied, and the entry stays armed.
	g.emitLasting(EventCreaturePlayed, 0, sanc)
	if !g.State.Cards[sanc].Exhausted {
		t.Error("a non-Mars creature should not be readied")
	}
	if g.State.LastingCount != 1 {
		t.Errorf("entry should remain after a filtered subject, count = %d", g.State.LastingCount)
	}

	// A Mars subject is readied, and the one-shot entry removes itself.
	g.emitLasting(EventCreaturePlayed, 0, mars)
	if g.State.Cards[mars].Exhausted {
		t.Error("the next Mars creature should enter ready")
	}
	if g.State.LastingCount != 0 {
		t.Errorf("the one-shot entry should be removed, count = %d", g.State.LastingCount)
	}
}

func TestLastingOnceOrdersWithPersistentReaction(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(2)), 0)
	g.State.Cards[mars].Exhausted = true

	// The one-shot sits between two persistent reactions, so removing it scans past
	// the first and shifts the last down.
	g.AddLasting(EventCreaturePlayed, actGainAember, 0, 1)
	g.AddLastingOnce(EventCreaturePlayed, actReadyPlayed, 0, 0, Mars)
	g.AddLasting(EventCreaturePlayed, actGainAember, 0, 1)

	g.emitLasting(EventCreaturePlayed, 0, mars) // three reactions fire; ordering path runs

	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2", g.Aember(0))
	}
	if g.State.Cards[mars].Exhausted {
		t.Error("the Mars creature should be readied")
	}
	if g.State.LastingCount != 2 {
		t.Errorf("count = %d, want 2 (both persistent reactions remain)", g.State.LastingCount)
	}
}
