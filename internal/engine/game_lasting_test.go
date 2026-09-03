package engine

import "testing"

// TestLastingActionDescriptions pins the short labels the ordering prompt shows
// when several reactions fire at once (ADR 0013).
func TestLastingActionDescriptions(t *testing.T) {
	for act, want := range map[lastingAction]string{
		actGainAember:  "gain Æmber",
		actDealDamage:  "deal damage",
		actReadyPlayed: "ready the creature",
		actDraw:        "draw a card",
	} {
		if got := act.describe(); got != want {
			t.Errorf("describe(%d) = %q, want %q", act, got, want)
		}
	}
}

func TestEmitLastingOrders(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	// Two reactions fire on the same event; the controller orders them.
	g.AddLasting(
		LastingEffect{On: EventCreaturePlayed, Do: actGainAember, Controller: 0, Amount: 1},
	)
	g.AddLasting(
		LastingEffect{On: EventCreaturePlayed, Do: actDealDamage, Controller: 0, Amount: 2},
	)
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
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actReadyPlayed,
			Controller: 0,
			Amount:     0,
			House:      Mars,
			Type:       Creature,
			Once:       true,
		},
	)
	if g.State.LastingCount != 1 {
		t.Fatalf("setup: lasting count = %d, want 1", g.State.LastingCount)
	}

	g.EndPlayPhase(0)

	if g.State.LastingCount != 0 {
		t.Errorf(
			"a one-shot reaction that never fired should be cleared at end of turn, count = %d",
			g.State.LastingCount,
		)
	}
}

func TestClearLastingKeepsOtherPlayer(t *testing.T) {
	g := started(t)
	g.AddLasting(
		LastingEffect{On: EventCreaturePlayed, Do: actGainAember, Controller: 0, Amount: 1},
	)
	g.AddLasting(
		LastingEffect{On: EventCreaturePlayed, Do: actGainAember, Controller: 1, Amount: 1},
	) // the opponent's reaction

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
		g.AddLasting(
			LastingEffect{On: EventCreaturePlayed, Do: actGainAember, Controller: 0, Amount: 1},
		)
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

	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actReadyPlayed,
			Controller: 0,
			Amount:     0,
			House:      Mars,
			Type:       Creature,
			Once:       true,
		},
	)

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

// TestLastingOnceFiltersByCardType covers the type filter a "next creature or
// artifact" entry carries (Soft Landing): an upgrade play is passed over, while a
// creature and an artifact both satisfy it.
func TestLastingOnceFiltersByCardType(t *testing.T) {
	g := NewGame("A", "B", 1)
	up := g.Register(NewCard("u", Mars, Upgrade, Common), 0)
	artifact := g.Register(NewCard("a", Mars, Artifact, Common), 0)
	g.State.Artifacts[0].add(artifact)
	g.State.Cards[artifact].Exhausted = true

	g.AddLasting(
		LastingEffect{
			On:         EventCardEntersPlay,
			Do:         actReadyPlayed,
			Controller: 0,
			Amount:     0,
			House:      HouseNone,
			Type:       AnyType,
			Once:       true,
		},
	)

	g.emitLasting(EventCardEntersPlay, 0, up)
	if g.State.LastingCount != 1 {
		t.Errorf("an upgrade should not satisfy the entry, count = %d", g.State.LastingCount)
	}

	g.emitLasting(EventCardEntersPlay, 0, artifact)
	if g.State.Cards[artifact].Exhausted {
		t.Error("the artifact should enter ready")
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
	g.AddLasting(
		LastingEffect{On: EventCreaturePlayed, Do: actGainAember, Controller: 0, Amount: 1},
	)
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actReadyPlayed,
			Controller: 0,
			Amount:     0,
			House:      Mars,
			Type:       Creature,
			Once:       true,
		},
	)
	g.AddLasting(
		LastingEffect{On: EventCreaturePlayed, Do: actGainAember, Controller: 0, Amount: 1},
	)

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
