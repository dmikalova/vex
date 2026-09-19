package engine

import (
	"errors"
	"testing"
)

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
	notCreature.Types = CardTypesOf(Artifact, Upgrade, Tactic)
	if notCreature.frees(g, 0, Mars, Creature) {
		t.Error("a type outside the permit's set is not freed")
	}
	if !notCreature.frees(g, 0, Mars, Artifact) {
		t.Error("a type in the permit's set is freed")
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
		if err := g.CanPlay(0, second); !errors.Is(err, ErrWrongHouse) {
			t.Fatalf("second off-house play = %v, want ErrWrongHouse", err)
		}
	})

	t.Run("excluded house is not freed", func(t *testing.T) {
		g := started(t)
		g.addOffHousePermit(0, OffHousePermit{Except: StarAlliance, Grant: GrantPlay, Remaining: 1})
		sa := g.AddToHand(NewCard("sa card", StarAlliance, Creature, Common, WithPower(3)), 0)
		if err := g.CanPlay(0, sa); !errors.Is(err, ErrWrongHouse) {
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
		if err := g.CanPlay(0, second); !errors.Is(err, ErrWrongHouse) {
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
		if err := g.CanPlay(0, mars); !errors.Is(err, ErrWrongHouse) {
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

// TestTypeUnlimitedPermissionText renders Matter Maker's type-scoped waiver, and
// checks the plural noun for a single type and for a set of types.
func TestTypeUnlimitedPermissionText(t *testing.T) {
	got := playPermissionText(PlayPermission{Types: CardTypesOf(Upgrade)})
	want := "You may play upgrades as if they were in the active house."
	if got != want {
		t.Errorf("type waiver text = %q, want %q", got, want)
	}
	if got := CardTypesOf(Artifact, Upgrade).playablePlural(); got != "artifacts or upgrades" {
		t.Errorf("two-type plural = %q, want %q", got, "artifacts or upgrades")
	}
	if got := CardTypesOf(Creature, Artifact, Upgrade).playablePlural(); got !=
		"creatures, artifacts, or upgrades" {
		t.Errorf("three-type plural = %q", got)
	}
	if got := CardTypes(0).playablePlural(); got != "cards" {
		t.Errorf("empty plural = %q, want %q", got, "cards")
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
