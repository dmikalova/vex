package engine

import "testing"

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

	(TakeControl{Duration: UntilThisLeavesPlay}).Resolve(&EffectContext{Resolver: g, Source: host, Upgrade: collar, Controller: 0})

	if g.owner(host) != 1 {
		t.Fatalf("owner = %d, want unchanged P2", g.owner(host))
	}
	if g.controller(host) != 0 {
		t.Fatalf("controller = %d, want P1", g.controller(host))
	}
	if g.State.Cards[host].ControlSource != collar {
		t.Fatalf("control source = %d, want the collar %d", g.State.Cards[host].ControlSource, collar)
	}
	if !g.State.Battleline[0].contains(host) || g.State.Battleline[1].contains(host) {
		t.Fatalf("battlelines = %v/%v, want host only under P1", g.Battleline(0), g.Battleline(1))
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

func TestTakeControlArtifact(t *testing.T) {
	permArt := TakeControl{Target: Target{Kind: TargetChosenEnemyArtifact}, Duration: Permanent}
	if got := permArt.Text(); got != "take control of an enemy artifact" {
		t.Errorf("permanent artifact text = %q", got)
	}
	if (TakeControl{Duration: Permanent}).validate() != nil {
		t.Error("Permanent should be a valid duration")
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
		TakeControl{Target: Target{Kind: TargetChosenEnemyCreature}, Duration: UntilThisLeavesPlay}.Resolve(ctx)
		if g.controller(foe) != 0 {
			t.Errorf("controller = %d, want 0", g.controller(foe))
		}
		if !ctx.HasIt || ctx.It != foe {
			t.Error("the taken creature should be left in context as It")
		}
	})
}

func TestBelongToHouseIfOffIdentity(t *testing.T) {
	e := BelongToHouseIfOffIdentity{House: Shadows}
	if e.Text() != "if it does not belong to a house on your identity then it belongs to house Shadows" {
		t.Errorf("text = %q", e.Text())
	}
	if (BelongToHouseIfOffIdentity{}).validate() == nil {
		t.Error("a missing house should be invalid")
	}
	if e.validate() != nil {
		t.Error("a set house should validate")
	}

	t.Run("reassigns an off-identity card", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetPlayerHouses(0, []House{Shadows, Logos, Sanctum})
		art := g.AddArtifact(NewCard("Gizmo", Untamed, Artifact, Common), 0)
		e.Resolve(&EffectContext{Resolver: g, Controller: 0, It: art, HasIt: true})
		if g.House(art) != Shadows {
			t.Errorf("house = %v, want Shadows", g.House(art))
		}
	})

	t.Run("leaves an on-identity card unchanged", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetPlayerHouses(0, []House{Shadows, Mars, Sanctum})
		art := g.AddArtifact(NewCard("Gizmo", Mars, Artifact, Common), 0)
		e.Resolve(&EffectContext{Resolver: g, Controller: 0, It: art, HasIt: true})
		if g.House(art) != Mars {
			t.Errorf("house = %v, want Mars (unchanged)", g.House(art))
		}
	})

	t.Run("does nothing without a card in context", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	})
}

func TestTakeControlOfArtifactGuard(t *testing.T) {
	g := NewGame("A", "B", 1)
	creature := g.AddToBattleline(testCreature("c", 3), 1)
	g.takeControlOfArtifact(creature, 0) // not an artifact: no-op
	if g.controller(creature) != 1 {
		t.Error("a non-artifact should not be seized")
	}
}
