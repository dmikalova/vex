package engine

import "testing"

// TestMayPlayOffHouseText covers the rendered text and validation of the
// off-house play grant across its Controlled, exclusion, type-filter, and
// play-or-use forms.
func TestMayPlayOffHouseText(t *testing.T) {
	cases := []struct {
		name string
		e    MayPlayOffHouse
		want string
	}{
		{
			"controlled",
			MayPlayOffHouse{Controlled: true, Grant: GrantPlay},
			"for the remainder of the turn, you may play cards from any house for which you have a card in play",
		},
		{
			"exclusion play one card",
			MayPlayOffHouse{Except: StarAlliance, Grant: GrantPlay, Count: 1},
			"you may play one non-Star Alliance card this turn",
		},
		{
			"play or use",
			MayPlayOffHouse{Except: StarAlliance, Grant: GrantPlay | GrantUse, Count: 1},
			"you may play or use one non-Star Alliance card this turn",
		},
		{
			"non-creature types",
			MayPlayOffHouse{Except: StarAlliance, NotType: Creature, Grant: GrantPlay, Count: 1},
			"you may play a non-Star Alliance artifact, upgrade, or Tactic this turn",
		},
		{
			"no excluded house",
			MayPlayOffHouse{Grant: GrantPlay, Count: 1},
			"you may play one card this turn",
		},
	}
	for _, c := range cases {
		if got := c.e.Text(); got != c.want {
			t.Errorf("%s: Text() = %q, want %q", c.name, got, c.want)
		}
	}

	if (MayPlayOffHouse{Count: 1}).validate() == nil {
		t.Error("a grant with no capability should be invalid")
	}
	if (MayPlayOffHouse{Grant: GrantPlay, Count: -1}).validate() == nil {
		t.Error("a negative count should be invalid")
	}
	if (MayPlayOffHouse{Grant: GrantPlay}).validate() != nil {
		t.Error("an unbounded play grant should be valid")
	}
}

// TestMayPlayOffHouseResolve confirms resolving the effect stores a permit whose
// Remaining tracks the count (unbounded stays permitUnlimited), records the log,
// and that the ready phase clears it.
func TestMayPlayOffHouseResolve(t *testing.T) {
	g := started(t)
	MayPlayOffHouse{Except: StarAlliance, Grant: GrantPlay, Count: 2}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if g.State.OffHousePermitCount[0] != 1 {
		t.Fatalf("permit count = %d, want 1", g.State.OffHousePermitCount[0])
	}
	if got := g.State.OffHousePermits[0][0].Remaining; got != 2 {
		t.Errorf("bounded Remaining = %d, want 2", got)
	}

	MayPlayOffHouse{Controlled: true, Grant: GrantPlay}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if got := g.State.OffHousePermits[0][1].Remaining; got != permitUnlimited {
		t.Errorf("unbounded Remaining = %d, want permitUnlimited", got)
	}

	g.EndPlayPhase(0)
	g.StartTurn(0)
	if g.State.OffHousePermitCount[0] != 0 {
		t.Error("ready phase should clear off-house permits")
	}
}

// TestOffHousePermitFrees exercises each reason a permit refuses a card.
func TestOffHousePermitFrees(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("brute", 3), 0) // a Brobnar card in play

	base := OffHousePermit{Grant: GrantPlay, Remaining: 1}
	if !base.frees(g, 0, Mars, Creature) {
		t.Error("an open permit should free any card")
	}
	spent := base
	spent.Remaining = 0
	if spent.frees(g, 0, Mars, Creature) {
		t.Error("a spent permit frees nothing")
	}
	excl := base
	excl.Except = Mars
	if excl.frees(g, 0, Mars, Creature) {
		t.Error("the excluded house is not freed")
	}
	if !excl.frees(g, 0, Dis, Creature) {
		t.Error("a house other than the excluded one is freed")
	}
	ctrl := base
	ctrl.Controlled = true
	if ctrl.frees(g, 0, Mars, Creature) {
		t.Error("Controlled frees only houses with a card in play")
	}
	if !ctrl.frees(g, 0, Brobnar, Creature) {
		t.Error("Controlled frees a house the player has in play")
	}
	notCreature := base
	notCreature.NotType = Creature
	if notCreature.frees(g, 0, Mars, Creature) {
		t.Error("the excluded type is not freed")
	}
	if !notCreature.frees(g, 0, Mars, Artifact) {
		t.Error("a type other than the excluded one is freed")
	}
}

// TestAddOffHousePermitCap confirms the store drops grants past its cap.
func TestAddOffHousePermitCap(t *testing.T) {
	g := started(t)
	for i := 0; i < maxOffHousePermits+2; i++ {
		g.addOffHousePermit(0, OffHousePermit{Grant: GrantPlay, Remaining: 1})
	}
	if int(g.State.OffHousePermitCount[0]) != maxOffHousePermits {
		t.Errorf("permit count = %d, want %d", g.State.OffHousePermitCount[0], maxOffHousePermits)
	}
}

// TestOffHousePlayViaPermit plays a non-active card off a stored permit, checks the
// permit's Remaining decrements, and that an unbounded permit is untouched.
func TestOffHousePlayViaPermit(t *testing.T) {
	t.Run("bounded permit decrements and runs out", func(t *testing.T) {
		g := started(t) // Brobnar active
		g.addOffHousePermit(0, OffHousePermit{Except: StarAlliance, Grant: GrantPlay, Remaining: 1})
		first := g.AddToHand(NewCard("mars a", Mars, Creature, Common, WithPower(3)), 0)
		second := g.AddToHand(NewCard("mars b", Mars, Creature, Common, WithPower(3)), 0)

		if err := g.CanPlay(0, first); err != nil {
			t.Fatalf("CanPlay off-house Mars = %v, want nil", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
			t.Fatalf("play off-house: %v", err)
		}
		if got := g.State.OffHousePermits[0][0].Remaining; got != 0 {
			t.Errorf("Remaining after one play = %d, want 0", got)
		}
		if err := g.CanPlay(0, second); err != ErrWrongHouse {
			t.Fatalf("second off-house play = %v, want ErrWrongHouse", err)
		}
	})

	t.Run("excluded house is not freed", func(t *testing.T) {
		g := started(t)
		g.addOffHousePermit(0, OffHousePermit{Except: StarAlliance, Grant: GrantPlay, Remaining: 1})
		sa := g.AddToHand(NewCard("sa card", StarAlliance, Creature, Common, WithPower(3)), 0)
		if err := g.CanPlay(0, sa); err != ErrWrongHouse {
			t.Fatalf("CanPlay excluded house = %v, want ErrWrongHouse", err)
		}
	})

	t.Run("unbounded permit is untouched", func(t *testing.T) {
		g := started(t)
		g.addOffHousePermit(0, OffHousePermit{Grant: GrantPlay, Remaining: permitUnlimited})
		a := g.AddToHand(NewCard("mars a", Mars, Creature, Common, WithPower(3)), 0)
		b := g.AddToHand(NewCard("mars b", Mars, Creature, Common, WithPower(3)), 0)
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, a), false); err != nil {
			t.Fatalf("first unbounded play: %v", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, b), false); err != nil {
			t.Fatalf("second unbounded play: %v", err)
		}
		if got := g.State.OffHousePermits[0][0].Remaining; got != permitUnlimited {
			t.Errorf("unbounded Remaining = %d, want permitUnlimited", got)
		}
	})
}

// TestOffHouseUseViaPermit spends a use permit reaping, fighting, and using an
// action ability, and confirms an active-house card spends nothing.
func TestOffHouseUseViaPermit(t *testing.T) {
	t.Run("reap spends one use", func(t *testing.T) {
		g := started(t) // Brobnar active
		g.addOffHousePermit(0, OffHousePermit{Except: StarAlliance, Grant: GrantUse, Remaining: 2})
		mars := g.AddToBattleline(NewCard("mars reaper", Mars, Creature, Common, WithPower(3)), 0)
		if err := g.Reap(0, mars); err != nil {
			t.Fatalf("off-house reap = %v", err)
		}
		if got := g.State.OffHousePermits[0][0].Remaining; got != 1 {
			t.Errorf("Remaining after reap = %d, want 1", got)
		}
	})

	t.Run("active-house use spends nothing", func(t *testing.T) {
		g := started(t)
		g.addOffHousePermit(0, OffHousePermit{Except: StarAlliance, Grant: GrantUse, Remaining: 1})
		brob := g.AddToBattleline(testCreature("brob reaper", 3), 0)
		if err := g.Reap(0, brob); err != nil {
			t.Fatalf("active-house reap = %v", err)
		}
		if got := g.State.OffHousePermits[0][0].Remaining; got != 1 {
			t.Errorf("active-house reap spent a use: Remaining = %d, want 1", got)
		}
	})

	t.Run("fight spends one use", func(t *testing.T) {
		g := started(t)
		g.addOffHousePermit(0, OffHousePermit{Except: StarAlliance, Grant: GrantUse, Remaining: 1})
		attacker := g.AddToBattleline(
			NewCard("mars fighter", Mars, Creature, Common, WithPower(4)),
			0,
		)
		defender := g.AddToBattleline(testCreature("enemy", 3), 1)
		if err := g.Fight(0, attacker, defender); err != nil {
			t.Fatalf("off-house fight = %v", err)
		}
		if got := g.State.OffHousePermits[0][0].Remaining; got != 0 {
			t.Errorf("Remaining after fight = %d, want 0", got)
		}
	})

	t.Run("action ability spends one use", func(t *testing.T) {
		g := started(t)
		g.addOffHousePermit(0, OffHousePermit{Except: StarAlliance, Grant: GrantUse, Remaining: 1})
		def := NewCard(
			"mars actor",
			Mars,
			Creature,
			Common,
			WithPower(3),
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1}),
		)
		actor := g.AddToBattleline(def, 0)
		if err := g.UseAction(0, actor); err != nil {
			t.Fatalf("off-house action = %v", err)
		}
		if got := g.State.OffHousePermits[0][0].Remaining; got != 0 {
			t.Errorf("Remaining after action = %d, want 0", got)
		}
	})
}

// TestNonActivePlayPermission covers Captain Val Jericho's live "any non-active
// house" grant: it applies only while centered, decrements a shared counter, and
// resets each turn.
func TestNonActivePlayPermission(t *testing.T) {
	jericho := NewCard(
		"Val",
		StarAlliance,
		Creature,
		Rare,
		WithPower(5),
		WithPlayPermission(
			PlayPermission{NonActive: true, Condition: SourceInCenterOfBattleline{}, Amount: 1},
		),
	)

	t.Run("frees one off-house play while centered", func(t *testing.T) {
		g := started(t)               // Brobnar active
		g.AddToBattleline(jericho, 0) // sole creature: centered
		first := g.AddToHand(NewCard("mars a", Mars, Creature, Common, WithPower(3)), 0)
		second := g.AddToHand(NewCard("mars b", Mars, Creature, Common, WithPower(3)), 0)

		if err := g.CanPlay(0, first); err != nil {
			t.Fatalf("CanPlay via Val = %v, want nil", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
			t.Fatalf("first non-active play: %v", err)
		}
		if got := g.State.NonActivePlaysUsedThisTurn[0]; got != 1 {
			t.Errorf("non-active plays used = %d, want 1", got)
		}
		if err := g.CanPlay(0, second); err != ErrWrongHouse {
			t.Fatalf("second non-active play = %v, want ErrWrongHouse", err)
		}

		g.EndPlayPhase(0)
		g.StartTurn(0)
		if got := g.State.NonActivePlaysUsedThisTurn[0]; got != 0 {
			t.Errorf("non-active counter not reset: %d", got)
		}
	})

	t.Run("no grant while off-center", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(testCreature("flank a", 3), 0)
		g.AddToBattleline(jericho, 0) // even line: no center
		mars := g.AddToHand(NewCard("mars c", Mars, Creature, Common, WithPower(3)), 0)
		if err := g.CanPlay(0, mars); err != ErrWrongHouse {
			t.Fatalf("off-center Val granted a play: %v", err)
		}
		if got := g.nonActivePlayLimit(0); got != 0 {
			t.Errorf("off-center limit = %d, want 0", got)
		}
	})
}

// TestNonActivePlayLimitUncentered covers the NonActive grant without the centered
// gate, and the remaining-count exhaustion path.
func TestNonActivePlayLimitUncentered(t *testing.T) {
	g := started(t)
	anywhere := NewCard(
		"Anywhere",
		StarAlliance,
		Creature,
		Rare,
		WithPower(3),
		WithPlayPermission(PlayPermission{NonActive: true, Amount: 1}),
	)
	g.AddToBattleline(anywhere, 0)
	if got := g.nonActivePlayLimit(0); got != 1 {
		t.Fatalf("uncentered NonActive limit = %d, want 1", got)
	}
	if got := g.nonActivePlayRemaining(0); got != 1 {
		t.Errorf("remaining = %d, want 1", got)
	}
	g.State.NonActivePlaysUsedThisTurn[0] = 1
	if got := g.nonActivePlayRemaining(0); got != 0 {
		t.Errorf("remaining after use = %d, want 0", got)
	}
}

// TestNonActivePlayPermissionText renders Val Jericho's centered grant and the
// uncentered NonActive form.
func TestNonActivePlayPermissionText(t *testing.T) {
	centered := playPermissionText(
		PlayPermission{NonActive: true, Condition: SourceInCenterOfBattleline{}, Amount: 1},
	)
	want := "During your turn, if " + SelfName + " is in the center of your battleline, you may play one card that is not of the active house."
	if centered != want {
		t.Errorf("centered text = %q, want %q", centered, want)
	}
	uncentered := playPermissionText(PlayPermission{NonActive: true, Amount: 2})
	wantU := "During your turn you may play 2 cards that is not of the active house."
	if uncentered != wantU {
		t.Errorf("uncentered text = %q, want %q", uncentered, wantU)
	}
}

// TestCannotUseThisTurnBar arms the this-turn use bar and confirms it stops the
// player using creatures until the ready phase lifts it.
func TestCannotUseThisTurnBar(t *testing.T) {
	g := started(t)
	reaper := g.AddToBattleline(testCreature("reaper", 3), 0)
	g.CannotUseThisTurn(0, reaper)
	if err := g.Reap(0, reaper); err == nil {
		t.Error("a barred player should not be able to reap")
	}
	g.EndPlayPhase(0)
	g.StartTurn(0)
	if g.State.CannotUse[0].Value {
		t.Error("the this-turn use bar should lift at the ready phase")
	}
}
