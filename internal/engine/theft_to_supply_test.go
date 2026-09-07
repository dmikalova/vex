package engine

import "testing"

// posPixies models Po's Pixies: while it is in play, Æmber stolen or captured from
// its controller's pool is drawn from the common supply instead, so the controller
// keeps their own Æmber while the thief still gains it.
func posPixies() CardDefinition {
	return NewCard("Po's Pixies", Untamed, Creature, Rare,
		WithPower(1), WithReplaces(Instead{
			Of: EventAemberTakenFromPool, Player: Controller, With: FromCommonSupply,
		}))
}

func TestTheftRedirectedToSupplySteal(t *testing.T) {
	g := NewGame("A", "B", 1)
	// P1 controls the Pixies and holds 3 Æmber; P0 steals 2 from P1.
	g.AddToBattleline(posPixies(), 1)
	src := g.AddToBattleline(testCreature("thief", 1), 0)
	g.State.Aember[1] = 3
	StealAember{Amount: 2}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
	// The thief gains 2, but the victim keeps their Æmber (the 2 came from supply).
	if g.State.Aember[0] != 2 || g.State.Aember[1] != 3 {
		t.Errorf("after redirected steal: you=%d opp=%d, want 2/3",
			g.State.Aember[0], g.State.Aember[1])
	}
}

func TestTheftRedirectedToSupplyCapture(t *testing.T) {
	g := NewGame("A", "B", 1)
	// P0 controls the capturer; P1 controls the Pixies and holds 3 Æmber.
	g.AddToBattleline(posPixies(), 1)
	captor := g.AddToBattleline(testCreature("captor", 1), 0)
	g.State.Aember[1] = 3
	CaptureAember{Amount: 2, Target: Target{Kind: TargetThisCreature}, Source: Opponent}.
		Resolve(&EffectContext{Resolver: g, Source: captor, Controller: 0})
	// The creature captures 2, but P1 keeps their pool (the 2 came from supply).
	if g.AmberOn(captor) != 2 || g.State.Aember[1] != 3 {
		t.Errorf("after redirected capture: onCreature=%d opp=%d, want 2/3",
			g.AmberOn(captor), g.State.Aember[1])
	}
}

func TestTheftRedirectedToSupplyReader(t *testing.T) {
	g := NewGame("A", "B", 1)
	if g.AemberTakenFromSupply(0) {
		t.Error("no Pixies in play, want not redirected")
	}
	g.AddToBattleline(posPixies(), 0)
	if !g.AemberTakenFromSupply(0) {
		t.Error("Pixies in play, want redirected")
	}
}

// A source replacement scoped to the card's Opponent redirects that opponent's
// pool, the mirror of how Ether Spider's Player scoping works on the destination
// side. No printed card does this today; the branch stays general with the
// destination half.
func TestAemberTakenFromSupplyOpponentScoped(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(
		NewCard("Denier", Untamed, Creature, Rare, WithPower(1), WithReplaces(Instead{
			Of: EventAemberTakenFromPool, Player: Opponent, With: FromCommonSupply,
		})),
		0,
	)
	if g.AemberTakenFromSupply(0) {
		t.Error("the card's own pool is not the watched pool when scoped to Opponent")
	}
	if !g.AemberTakenFromSupply(1) {
		t.Error("opponent-scoped source should redirect the opponent's pool")
	}
}
