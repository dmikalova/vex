package engine

import "testing"

// TestMoveAemberAll checks the All mode empties every source instead of moving a
// fixed amount, and that it cannot be combined with one.
func TestMoveAemberAll(t *testing.T) {
	e := MoveAember{
		All:  true,
		From: Target{Kind: TargetEachEnemyCreature},
		To:   Controller,
	}
	want := "move all Æmber from each enemy creature to your pool"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if (MoveAember{All: true, Amount: 1, From: e.From, To: Controller}).validate() == nil {
		t.Error("All combined with Amount should not validate")
	}

	g := NewGame("A", "B", 1)
	rich := g.AddToBattleline(testCreature("rich", 4), 1)
	poor := g.AddToBattleline(testCreature("poor", 4), 1)
	g.AddAmberOn(rich, 3)
	g.AddAmberOn(poor, 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.AmberOn(rich) != 0 || g.AmberOn(poor) != 0 {
		t.Errorf("Æmber left on the creatures = %d/%d, want 0/0",
			g.AmberOn(rich), g.AmberOn(poor))
	}
	if got := g.Aember(0); got != 4 {
		t.Errorf("pool = %d, want 4", got)
	}
}

// TestMoveAemberFraction checks the Fraction mode moves a share of the source's
// Æmber (Patronage moves half, rounding up), renders the fractional phrase, and
// cannot be combined with a fixed Amount or with All.
func TestMoveAemberFraction(t *testing.T) {
	friendly := Target{Kind: TargetChosenFriendlyCreatureOrArtifact}
	e := MoveAember{From: friendly, Fraction: HalfRoundedUp, To: Controller}

	want := "move half the Æmber from a friendly creature or artifact to your pool, rounding up"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if (MoveAember{From: friendly, To: Controller, Fraction: HalfRoundedUp, Amount: 1}).
		validate() == nil {
		t.Error("Fraction combined with Amount should not validate")
	}
	if (MoveAember{From: friendly, To: Controller, Fraction: HalfRoundedUp, All: true}).
		validate() == nil {
		t.Error("Fraction combined with All should not validate")
	}

	// Odd pool: half of 3 rounds up to 2.
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.AddAmberOn(c, 3)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.AmberOn(c) != 1 || g.Aember(0) != 2 {
		t.Errorf("after half move: card=%d pool=%d, want 1/2", g.AmberOn(c), g.Aember(0))
	}
}

func TestMoveAember(t *testing.T) {
	friendly := Target{Kind: TargetChosenFriendlyCreatureOrArtifact}

	// Text: pool destination and card destination.
	toPool := MoveAember{Amount: 1, From: friendly, To: Controller}
	if got := toPool.Text(); got != "move 1 \u00c6mber from a friendly creature or artifact to your pool" {
		t.Errorf("pool text = %q", got)
	}
	toOpp := MoveAember{Amount: 2, From: friendly, To: Opponent}
	if got := toOpp.Text(); got != "move 2 \u00c6mber from a friendly creature or artifact to your opponent's pool" {
		t.Errorf("opponent-pool text = %q", got)
	}
	toCard := MoveAember{Amount: 1, From: friendly, Onto: Target{Kind: TargetChosenEnemyCreature}}
	if got := toCard.Text(); got != "move 1 \u00c6mber from a friendly creature or artifact to an enemy creature" {
		t.Errorf("card text = %q", got)
	}

	// validate: source and exactly one destination.
	if (MoveAember{Amount: 1, To: Controller}).validate() == nil {
		t.Error("unset source should be invalid")
	}
	if (MoveAember{Amount: 1, From: friendly}).validate() == nil {
		t.Error("no destination should be invalid")
	}
	if (MoveAember{Amount: 1, From: friendly, To: Controller, Onto: friendly}).validate() == nil {
		t.Error("two destinations should be invalid")
	}
	if (MoveAember{From: friendly, To: Controller}).validate() == nil {
		t.Error("an unset amount should be invalid without All or Fraction")
	}
	if toPool.validate() != nil || toCard.validate() != nil {
		t.Error("one destination should be valid")
	}

	// Resolve into a pool, capping the move at what the source holds.
	g := NewGame("A", "B", 1)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.AddAmberOn(c, 2)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	MoveAember{Amount: 3, From: friendly, To: Controller}.Resolve(ctx)
	if g.AmberOn(c) != 0 || g.Aember(0) != 2 {
		t.Errorf("after capped move: card=%d pool=%d, want 0/2", g.AmberOn(c), g.Aember(0))
	}

	// Resolve onto another card.
	g2 := NewGame("A", "B", 1)
	src := g2.AddToBattleline(testCreature("src", 3), 0)
	dst := g2.AddToBattleline(testCreature("dst", 3), 1)
	g2.AddAmberOn(src, 2)
	MoveAember{Amount: 1, From: friendly, Onto: Target{Kind: TargetChosenEnemyCreature}}.
		Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.AmberOn(src) != 1 || g2.AmberOn(dst) != 1 {
		t.Errorf("after card move: src=%d dst=%d, want 1/1", g2.AmberOn(src), g2.AmberOn(dst))
	}

	// No source carries Æmber: nothing happens.
	g3 := NewGame("A", "B", 1)
	g3.AddToBattleline(testCreature("bare", 3), 0)
	MoveAember{
		Amount: 1,
		From:   friendly,
		To:     Controller,
	}.Resolve(
		&EffectContext{Resolver: g3, Controller: 0},
	)
	if g3.Aember(0) != 0 {
		t.Error("moving with no Æmber-bearing card should do nothing")
	}

	// A card destination with no candidate card: nothing moves.
	g4 := NewGame("A", "B", 1)
	only := g4.AddToBattleline(testCreature("only", 3), 0)
	g4.AddAmberOn(only, 1)
	MoveAember{Amount: 1, From: friendly, Onto: Target{Kind: TargetChosenEnemyCreature}}.
		Resolve(&EffectContext{Resolver: g4, Controller: 0})
	if g4.AmberOn(only) != 1 {
		t.Error("a move with no destination card should move nothing")
	}

	// "to another creature" is another than the one the Æmber leaves: the source
	// goes into focus while the destination is chosen, so it cannot pick itself
	// (Consul Primus). Without the focus the destination widens to every creature
	// and the Æmber moves off a card and straight back onto it.
	g5 := NewGame("A", "B", 1)
	giver := g5.AddToBattleline(testCreature("giver", 3), 0)
	taker := g5.AddToBattleline(testCreature("taker", 3), 0)
	g5.AddAmberOn(giver, 1)
	MoveAember{Amount: 1, From: friendly, Onto: Target{Kind: TargetChosenOtherCreature}}.
		Resolve(&EffectContext{Resolver: g5, Controller: 0})
	if g5.AmberOn(giver) != 0 || g5.AmberOn(taker) != 1 {
		t.Errorf(
			"after move onto another creature: giver=%d taker=%d, want 0/1",
			g5.AmberOn(giver), g5.AmberOn(taker),
		)
	}
}

// TestMoveAemberBind checks that Bind selects through a Refinement on the full
// candidate set (not only the Æmber-bearers), moves all the Æmber, and leaves the
// moved-from creature in context.
func TestMoveAemberBind(t *testing.T) {
	mostPowerful := Target{Kind: TargetEachCreature}.Refine(MostPowerful)

	// The most powerful creature carries Æmber: it is emptied into the pool and
	// bound as ctx.It.
	g := NewGame("A", "B", 1)
	big := g.AddToBattleline(testCreature("big", 6), 0)
	small := g.AddToBattleline(testCreature("small", 3), 0)
	g.AddAmberOn(big, 3)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	MoveAember{All: true, From: mostPowerful, To: Controller, Bind: true}.Resolve(ctx)
	if g.AmberOn(big) != 0 || g.Aember(0) != 3 {
		t.Errorf("after bind move: big=%d pool=%d, want 0/3", g.AmberOn(big), g.Aember(0))
	}
	if g.AmberOn(small) != 0 {
		t.Errorf("weaker creature Æmber = %d, want 0", g.AmberOn(small))
	}
	if !ctx.HasIt || ctx.It != big {
		t.Errorf("ctx.It = %v (HasIt %v), want %v bound", ctx.It, ctx.HasIt, big)
	}

	// The overall most powerful creature has no Æmber while a weaker one does: Bind
	// still selects and binds the most powerful (moving nothing), never the
	// Æmber-holder.
	g2 := NewGame("A", "B", 1)
	strong := g2.AddToBattleline(testCreature("strong", 6), 0)
	weak := g2.AddToBattleline(testCreature("weak", 3), 0)
	g2.AddAmberOn(weak, 2)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	MoveAember{All: true, From: mostPowerful, To: Controller, Bind: true}.Resolve(ctx2)
	if g2.AmberOn(strong) != 0 || g2.AmberOn(weak) != 2 || g2.Aember(0) != 0 {
		t.Errorf("after bind of Æmber-less top: strong=%d weak=%d pool=%d, want 0/2/0",
			g2.AmberOn(strong), g2.AmberOn(weak), g2.Aember(0))
	}
	if !ctx2.HasIt || ctx2.It != strong {
		t.Errorf("ctx.It = %v (HasIt %v), want %v bound", ctx2.It, ctx2.HasIt, strong)
	}

	// No creatures at all: nothing is bound.
	g3 := NewGame("A", "B", 1)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	MoveAember{All: true, From: mostPowerful, To: Controller, Bind: true}.Resolve(ctx3)
	if ctx3.HasIt {
		t.Error("binding with no creatures should leave ctx.It unset")
	}
}

func TestMoveAemberDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	b := g.AddToBattleline(testCreature("b", 3), 0)
	g.AddAmberOn(a, 1)
	g.AddAmberOn(b, 1)
	g.SetChooser(0, orderRejectChooser{})
	MoveAember{
		Amount: 1,
		From:   Target{Kind: TargetChosenFriendlyCreatureOrArtifact},
		To:     Controller,
	}.
		Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	if g.Aember(0) != 0 {
		t.Error("a declined move should move nothing")
	}
}
