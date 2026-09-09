package engine

import (
	"slices"
	"testing"
)

func TestTakeControl(t *testing.T) {
	if got := (TakeControl{Duration: UntilThisLeavesPlay}).Text(); got != "take control of this creature until {upgrade} leaves play" {
		t.Fatalf("TakeControl text = %q", got)
	}
	if err := (TakeControl{}).validate(); err == nil {
		t.Fatal("TakeControl without a duration should be rejected")
	}
	if err := (TakeControl{Duration: UntilThisLeavesPlay}).validate(); err != nil {
		t.Fatalf("valid TakeControl = %v", err)
	}

	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 1)
	collar := g.Register(NewCard("collar", Dis, Upgrade, Rare), 0)

	(TakeControl{Duration: UntilThisLeavesPlay}).Resolve(
		&EffectContext{Resolver: g, Source: host, Upgrade: collar, Controller: 0},
	)

	if g.owner(host) != 1 {
		t.Fatalf("owner = %d, want unchanged P2", g.owner(host))
	}
	if g.controller(host) != 0 {
		t.Fatalf("controller = %d, want P1", g.controller(host))
	}
	if src := controlSourceOf(g, host); src != collar {
		t.Fatalf("control source = %d, want the collar %d", src, collar)
	}
	if !g.State.Battleline[0].contains(host) || g.State.Battleline[1].contains(host) {
		t.Fatalf("battlelines = %v/%v, want host only under P1", g.Battleline(0), g.Battleline(1))
	}
}

// TestTakeControlAndExhaust covers Forever control that also exhausts the seized
// creature in one effect — Lord Invidius's "take control ... and exhaust it".
func TestTakeControlAndExhaust(t *testing.T) {
	take := TakeControl{
		Target:     Target{Kind: TargetChosenEnemyCreature},
		Duration:   Forever,
		AndExhaust: true,
	}
	if got := take.Text(); got != "take control of an enemy creature and exhaust it" {
		t.Fatalf("AndExhaust text = %q", got)
	}

	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	take.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

	if g.controller(foe) != 0 {
		t.Fatalf("controller = %d, want P1", g.controller(foe))
	}
	if !g.Exhausted(foe) {
		t.Fatal("seized creature should be exhausted")
	}
}

// TestTakeControlPlacesSeizedCreatureOnChosenFlank pins that the player gaining
// control places the seized creature on the flank they choose, and that an empty
// taker line — the creature's only home — is not worth a prompt.
func TestTakeControlPlacesSeizedCreatureOnChosenFlank(t *testing.T) {
	take := TakeControl{
		Target:   Target{Kind: TargetChosenEnemyCreature},
		Duration: UntilThisLeavesPlay,
	}

	// A creature already in the taker's line: the flank is asked and honored.
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	mine := g.AddToBattleline(testCreature("mine", 3), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	g.SetChooser(0, optionPicker{idx: 1}) // right flank
	take.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
	if got, want := g.Battleline(0), []LocalID{src, mine, foe}; !slices.Equal(got, want) {
		t.Fatalf("right-flank placement = %v, want %v", got, want)
	}

	// The left flank places the seized creature at the head of the line.
	g2 := NewGame("A", "B", 1)
	s2 := g2.AddToBattleline(testCreature("s2", 3), 0)
	f2 := g2.AddToBattleline(testCreature("f2", 3), 1)
	g2.SetChooser(0, optionPicker{idx: 0})
	take.Resolve(&EffectContext{Resolver: g2, Source: s2, Controller: 0})
	if got, want := g2.Battleline(0), []LocalID{f2, s2}; !slices.Equal(got, want) {
		t.Fatalf("left-flank placement = %v, want %v", got, want)
	}

	// An empty taker line: the seized creature has one home, so no flank is asked.
	g3 := NewGame("A", "B", 1)
	box := g3.AddArtifact(NewCard("box", Logos, Artifact, Rare), 0)
	f3 := g3.AddToBattleline(testCreature("f3", 3), 1)
	take.Resolve(&EffectContext{Resolver: g3, Source: box, Controller: 0})
	if got, want := g3.Battleline(0), []LocalID{f3}; !slices.Equal(got, want) {
		t.Fatalf("single placement = %v, want %v", got, want)
	}
}

func TestControlRevertsWhenTakingUpgradeLeaves(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 1)
	collar := g.Register(NewCard("collar", Dis, Upgrade, Rare), 0)
	g.AttachUpgrade(host, collar)
	g.takeControl(host, 0, collar)

	g.destroyAttachedUpgrade(collar)

	if g.controller(host) != 1 {
		t.Fatalf("controller = %d, want owner P2", g.controller(host))
	}
	if !g.State.Battleline[1].contains(host) || g.State.Battleline[0].contains(host) {
		t.Fatalf("battlelines = %v/%v, want host only under P2", g.Battleline(0), g.Battleline(1))
	}
	if !g.State.Discard[0].contains(collar) {
		t.Fatalf("collar discard = %v, want collar in owner discard", g.Discard(0))
	}
}

func TestControlledCreatureLeavesForOwner(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 1)
	collar := g.Register(NewCard("collar", Dis, Upgrade, Rare), 0)
	g.AttachUpgrade(host, collar)
	g.takeControl(host, 0, collar)

	g.DealDamage(0, []DamageTarget{{ID: host, Amount: 3}})

	if !g.State.Discard[1].contains(host) {
		t.Fatalf("host discard = %v/%v, want owner P2 discard", g.Discard(0), g.Discard(1))
	}
	if !g.State.Discard[0].contains(collar) {
		t.Fatalf("collar discard = %v, want owner P1 discard", g.Discard(0))
	}
	if g.State.Cards[host].ControlPlus != 0 {
		t.Fatalf("ControlPlus after leaving play = %d, want 0", g.State.Cards[host].ControlPlus)
	}
}

// A creature that captured Æmber and was then taken by the opponent gives that
// Æmber to the opponent of whoever controls it at death, not the owner's
// opponent: taken into enemy control, it returns its Æmber to its owner.
func TestCapturedAemberGoesToControllersOpponentOnDeath(t *testing.T) {
	g := started(t)
	mine := g.AddToBattleline(testCreature("mine", 3), 0)
	g.State.Cards[mine].Amber = 2 // captured Æmber sitting on the creature
	seizer := g.Register(NewCard("seizer", Dis, Artifact, Rare), 1)
	g.takeControl(mine, 1, seizer) // P2 takes control of P1's creature

	g.DealDamage(0, []DamageTarget{{ID: mine, Amount: 3}}) // it dies under P2

	if g.Aember(0) != 2 {
		t.Errorf("P1 pool = %d, want 2 (Æmber returns to the controller's opponent)", g.Aember(0))
	}
	if g.Aember(1) != 0 {
		t.Errorf("P2 pool = %d, want 0", g.Aember(1))
	}
}

func TestTakeControlArtifact(t *testing.T) {
	permArt := TakeControl{Target: Target{Kind: TargetChosenEnemyArtifact}, Duration: Forever}
	if got := permArt.Text(); got != "take control of an enemy artifact" {
		t.Errorf("permanent artifact text = %q", got)
	}
	if (TakeControl{Duration: Forever}).validate() != nil {
		t.Error("Forever should be a valid duration")
	}
	if got := (TakeControl{Target: Target{Kind: TargetChosenEnemyCreature}, Duration: UntilThisLeavesPlay}).Text(); got != "take control of an enemy creature until "+SelfName+" leaves play" {
		t.Errorf("targeted reverting text = %q", got)
	}

	t.Run("takes permanent control of a chosen enemy artifact and records it", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		art := g.AddArtifact(NewCard("Gizmo", Mars, Artifact, Common), 1)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		permArt.Resolve(ctx)
		if g.controller(art) != 0 {
			t.Errorf("controller = %d, want 0", g.controller(art))
		}
		if !ctx.HasIt || ctx.It != art {
			t.Error("the taken artifact should be left in context as It")
		}
	})

	t.Run("targeted creature control anchors to the source and records it", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 3), 0)
		foe := g.AddToBattleline(testCreature("foe", 3), 1)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
		TakeControl{
			Target:   Target{Kind: TargetChosenEnemyCreature},
			Duration: UntilThisLeavesPlay,
		}.Resolve(
			ctx,
		)
		if g.controller(foe) != 0 {
			t.Errorf("controller = %d, want 0", g.controller(foe))
		}
		if !ctx.HasIt || ctx.It != foe {
			t.Error("the taken creature should be left in context as It")
		}
	})

	t.Run("ToOpponent gives control away instead of taking it", func(t *testing.T) {
		giveAway := TakeControl{
			Target:     Target{Kind: TargetThisCreature},
			Duration:   Forever,
			ToOpponent: true,
		}
		if got := giveAway.Text(); got != "your opponent gains control of "+SelfName {
			t.Errorf("ToOpponent text = %q", got)
		}

		g := NewGame("A", "B", 1)
		art := g.AddArtifact(NewCard("Spangler Box", Logos, Artifact, Rare), 0)
		ctx := &EffectContext{Resolver: g, Source: art, Controller: 0}
		TakeControl{
			Target:     Target{Kind: TargetThisCreature},
			Duration:   Forever,
			ToOpponent: true,
		}.Resolve(ctx)
		if g.controller(art) != 1 {
			t.Errorf("controller = %d, want 1 (the opponent)", g.controller(art))
		}
	})

	t.Run("gives a chosen friendly artifact away and reports progress", func(t *testing.T) {
		give := TakeControl{
			Target:     Target{Kind: TargetChosenFriendlyArtifact},
			Duration:   Forever,
			ToOpponent: true,
		}
		if got := give.Text(); got != "your opponent gains control of a friendly artifact" {
			t.Errorf("friendly give text = %q", got)
		}

		g := NewGame("A", "B", 1)
		art := g.AddArtifact(NewCard("Gizmo", Sanctum, Artifact, Common), 0)
		ctx := &EffectContext{Resolver: g, Controller: 0}
		if !give.resolveGate(ctx) {
			t.Error("giving an artifact away should report progress")
		}
		if g.controller(art) != 1 {
			t.Errorf("controller = %d, want 1 (the opponent)", g.controller(art))
		}

		empty := NewGame("A", "B", 1)
		if give.resolveGate(&EffectContext{Resolver: empty, Controller: 0}) {
			t.Error("no friendly artifact should report no progress")
		}
	})
}

func TestItIsOffIdentity(t *testing.T) {
	c := ItIsOffIdentity{}
	if c.CondText() != "if it does not belong to a house on your identity" {
		t.Errorf("text = %q", c.CondText())
	}

	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Shadows, Logos, Sanctum})
	off := g.AddArtifact(NewCard("Gizmo", Untamed, Artifact, Common), 0)
	on := g.AddArtifact(NewCard("Relic", Logos, Artifact, Common), 0)

	if !c.Met(&EffectContext{Resolver: g, Controller: 0, It: off, HasIt: true}) {
		t.Error("an off-identity card should meet the condition")
	}
	if c.Met(&EffectContext{Resolver: g, Controller: 0, It: on, HasIt: true}) {
		t.Error("an on-identity card should not meet the condition")
	}
	if c.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("no card in context should not meet the condition")
	}
}

// controlSourceOf returns the source of the topmost control entry over a card, or
// 0 when the card carries none.
func controlSourceOf(g *Game, id LocalID) LocalID {
	for i := int(g.State.ControlCount) - 1; i >= 0; i-- {
		if g.State.Controls[i].Card == id {
			return g.State.Controls[i].Source
		}
	}
	return 0
}

func TestTakeControlOfArtifactGuard(t *testing.T) {
	g := NewGame("A", "B", 1)
	shelved := g.Register(NewCard("relic", Dis, Artifact, Rare), 1) // registered, not in play
	g.takeControl(shelved, 0, shelved)                              // not in play: no-op
	if g.controller(shelved) != 1 {
		t.Error("a card that is not in play should not be seized")
	}
	if g.State.ControlCount != 0 {
		t.Errorf("control stack = %d entries, want 0", g.State.ControlCount)
	}
}

// TestControlStackIsLIFO pins that stacked control effects revert in LIFO order:
// removing the newer effect falls back to the one beneath, and only removing the
// last returns the card to its owner.
func TestControlStackIsLIFO(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	srcA := g.AddToBattleline(testCreature("srcA", 3), 0)
	srcB := g.AddToBattleline(testCreature("srcB", 3), 0)

	g.takeControl(foe, 0, srcA) // A seizes foe for P0
	g.takeControl(foe, 0, srcB) // B stacks on top, also for P0

	// The newer effect lapsing falls back to A's control, not to the owner.
	g.releaseControlHeldBy(srcB)
	if g.controller(foe) != 0 {
		t.Fatalf("controller after B lapses = %d, want P0 (falls back to A)", g.controller(foe))
	}
	if src := controlSourceOf(g, foe); src != srcA {
		t.Fatalf("control source after B lapses = %d, want A %d", src, srcA)
	}

	// The last effect lapsing returns the creature to its owner.
	g.releaseControlHeldBy(srcA)
	if g.controller(foe) != 1 {
		t.Fatalf("controller after A lapses = %d, want owner P1", g.controller(foe))
	}
	if g.State.ControlCount != 0 {
		t.Fatalf("control stack = %d entries, want 0", g.State.ControlCount)
	}
}

// TestControlTableFull pins that overflowing the control stack panics rather than
// silently dropping a control effect.
func TestControlTableFull(t *testing.T) {
	g := NewGame("A", "B", 1)
	for range maxControlEntries {
		g.pushControl(LocalID(1), 0, LocalID(1))
	}
	defer func() {
		if recover() == nil {
			t.Fatal("an entry past capacity should panic")
		}
	}()
	g.pushControl(LocalID(1), 0, LocalID(1))
}

// TestPermanentControlDoesNotAccumulate pins that swapping a card's control every
// turn (Whirlpool) does not grow the stack: a permanent take supersedes the whole
// stack, so far more takes than the table's capacity never panic and leave one
// entry.
func TestPermanentControlDoesNotAccumulate(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	for i := range maxControlEntries * 3 {
		g.takeControl(foe, i%2, foe) // permanent (self-sourced), swapping sides
	}
	if g.State.ControlCount != 1 {
		t.Fatalf(
			"control stack = %d entries, want 1 after repeated permanent takes",
			g.State.ControlCount,
		)
	}
	if g.controller(foe) != (maxControlEntries*3-1)%2 {
		t.Fatalf("controller = %d, want the last taker", g.controller(foe))
	}
}

// TestRevertibleControlDedupesBySource pins that the same source re-taking a card
// replaces its own entry rather than stacking a duplicate.
func TestRevertibleControlDedupesBySource(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.takeControl(foe, 0, src)
	g.takeControl(foe, 0, src)
	if g.State.ControlCount != 1 {
		t.Fatalf("control stack = %d entries, want 1 (same source dedupes)", g.State.ControlCount)
	}
}

// A creature kept alive only by its own side's +power aura is destroyed at the
// resolution boundary once it changes sides and loses that aura.
func TestTakeControlDestroysNewlyLethalCreature(t *testing.T) {
	g := started(t)
	g.AddArtifact(NewCard("Banner", Brobnar, Artifact, Rare, WithConstantAbility(
		ConstantAbility{PowerBonus: 2, Target: Target{Kind: TargetEachFriendlyCreature}},
	)), 1)
	ape := g.AddToBattleline(testCreature("ape", 3), 1)
	g.State.Cards[ape].Damage = 3
	if got := g.Power(ape); got != 5 {
		t.Fatalf("power under its own side = %d, want 5 (3 + banner)", got)
	}

	g.takeControl(ape, 0, 0)
	g.settleDestroyed(0) // the resolution boundary settles the swap (ADR 0029)

	if g.inPlay(ape) {
		t.Fatalf("ape is still in play at %d power with 3 damage", g.Power(ape))
	}
	if !containsID(g.Discard(1), ape) {
		t.Errorf("discard = %v, want the ape in its owner's pile", g.Discard(1))
	}
}

// TestTakeControlIsNotLeavingPlay pins that a change of control is not an exit: it
// does not fire the creature's Leaves Play abilities and does not shed its
// counters — the creature stays in play, only on the other side.
func TestTakeControlIsNotLeavingPlay(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(
		testCreature(
			"foe",
			3,
			WithAbility(TriggerLeavesPlay, GainAember{Player: Controller, Amount: 1}),
		),
		1,
	)
	g.PlaceCounter(foe, CounterDoom, 2)
	before := g.Aember(1)

	g.takeControl(foe, 0, 0)

	if g.Aember(1) != before {
		t.Errorf(
			"a control change fired a Leaves Play ability: Æmber %d -> %d",
			before,
			g.Aember(1),
		)
	}
	if g.CountersOn(foe, CounterDoom) != 2 {
		t.Errorf(
			"counters = %d, want 2: a control change must not shed counters",
			g.CountersOn(foe, CounterDoom),
		)
	}
	if g.controller(foe) != 0 || !g.State.Battleline[0].contains(foe) {
		t.Errorf("controller = %d, want P0 with the creature in play", g.controller(foe))
	}
}
