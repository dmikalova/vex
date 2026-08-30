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
