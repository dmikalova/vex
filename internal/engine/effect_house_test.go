package engine

import "testing"

// TestHouseFromUpgradeOverride covers an Upgrade that overrides its host's house
// (Academy Training makes its creature a Logos creature).
func TestHouseFromUpgradeOverride(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(NewCard("Host", Brobnar, Creature, Common, WithPower(3)), 0)
	attachUpgrade(g, host,
		NewCard("Academy Training", Logos, Upgrade, Common,
			WithStatic(StaticModifier{HouseOverride: Logos})))
	if g.House(host) != Logos {
		t.Errorf("house = %s, want Logos (upgrade override)", g.House(host))
	}
}

func TestBelongToHouseValidate(t *testing.T) {
	if err := (BelongToHouse{House: Mars, Duration: RemainderOfPlayerTurn}).validate(); err == nil {
		t.Error("an unset target should be rejected")
	}
	if err := (BelongToHouse{Target: Target{Kind: TargetThisCreature}, Duration: RemainderOfPlayerTurn}).validate(); err == nil {
		t.Error("an unset house should be rejected")
	}
	if err := (BelongToHouse{Target: Target{Kind: TargetThisCreature}, House: Mars}).validate(); err == nil {
		t.Error("an unset duration should be rejected")
	}
	if err := (BelongToHouse{Target: Target{Kind: TargetThisCreature}, House: Mars, Duration: RemainderOfPlayerTurn}).validate(); err != nil {
		t.Errorf("valid BelongToHouse = %v", err)
	}
}

func TestBelongToHouseEndOfTurn(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(NewCard("Host", Brobnar, Creature, Common, WithPower(3)), 0)
	e := BelongToHouse{
		Target:   Target{Kind: TargetThisCreature},
		House:    Mars,
		Duration: RemainderOfPlayerTurn,
	}
	if got := e.Text(); got != "for the remainder of the turn, {self} belongs to house Mars" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(&EffectContext{Resolver: g, Source: host, Controller: 0})
	if g.House(host) != Mars {
		t.Fatalf("house = %s, want Mars", g.House(host))
	}
	g.EndPlayPhase(0)
	if g.House(host) != Brobnar {
		t.Fatalf("house after end of turn = %s, want Brobnar (lasts only this turn)", g.House(host))
	}
}

func TestBelongToHouseUntilLeavesPlay(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(NewCard("Host", Brobnar, Creature, Common, WithPower(3)), 0)
	e := BelongToHouse{
		Target:   Target{Kind: TargetThisCreature},
		House:    Mars,
		Duration: UntilThisLeavesPlay,
	}
	if got := e.Text(); got != "{self} belongs to house Mars until it leaves play" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(&EffectContext{Resolver: g, Source: host, Controller: 0})
	if g.House(host) != Mars {
		t.Fatalf("house = %s, want Mars", g.House(host))
	}
	g.EndPlayPhase(0)
	if g.House(host) != Mars {
		t.Fatalf(
			"house after end of turn = %s, want still Mars (lasts until it leaves play)",
			g.House(host),
		)
	}
	g.State.Battleline[0].remove(host)
	g.resetCore(host)
	if g.House(host) != Brobnar {
		t.Fatalf("house out of play = %s, want printed Brobnar", g.House(host))
	}
}
