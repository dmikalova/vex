package engine

import "testing"

func TestCaptureAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}
	g.State.Aember[1] = 2

	e := CaptureAember{
		Amount: 3,
		Target: Target{Kind: TargetThisCreature},
		Source: Opponent,
	}
	if e.Text() != "{self} captures 3 Æmber from your opponent" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx) // opponent has only 2
	if g.State.Cards[src].Amber != 2 {
		t.Errorf("captured = %d, want 2", g.State.Cards[src].Amber)
	}
	if g.State.Aember[1] != 0 {
		t.Errorf("opponent aember = %d, want 0", g.State.Aember[1])
	}
}

func TestMoveAemberToSupplyEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("aubade", 4), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}

	e := MoveAemberToSupply{
		Amount: 1,
		Target: Target{Kind: TargetThisCreature},
	}
	if e.Text() != "move 1 Æmber from {self} to the common supply" {
		t.Errorf("text = %q", e.Text())
	}

	e.Resolve(ctx) // nothing captured: no-op
	if got := g.AmberOn(src); got != 0 {
		t.Errorf("empty discard = %d, want 0", got)
	}

	g.AddAmberOn(src, 3)
	e.Resolve(ctx)
	if got := g.AmberOn(src); got != 2 {
		t.Errorf("after discarding 1 of 3 = %d, want 2", got)
	}

	// Discarding more than is held empties the creature rather than going negative.
	big := MoveAemberToSupply{
		Amount: 5,
		Target: Target{Kind: TargetThisCreature},
	}
	big.Resolve(ctx)
	if got := g.AmberOn(src); got != 0 {
		t.Errorf("over-discard = %d, want 0", got)
	}

	// A chosen target renders by its own noun rather than {self}.
	chosen := MoveAemberToSupply{
		Amount: 2,
		Target: Target{Kind: TargetChosenEnemyCreature},
	}
	if chosen.Text() != "move 2 Æmber from an enemy creature to the common supply" {
		t.Errorf("chosen text = %q", chosen.Text())
	}

	// validate rejects an unset target and a non-positive amount, accepts a valid one.
	if err := (MoveAemberToSupply{Amount: 1}).validate(); err == nil {
		t.Error("unset target should be rejected")
	}
	if err := (MoveAemberToSupply{Target: Target{Kind: TargetThisCreature}}).validate(); err == nil {
		t.Error("non-positive amount should be rejected")
	}
	if err := e.validate(); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}

	// All mode renders "each Æmber" and rejects a combined Amount.
	all := MoveAemberToSupply{
		All:    true,
		Target: Target{Kind: TargetThisCreature},
	}
	if got := all.Text(); got != "move each Æmber on {self} to the common supply" {
		t.Errorf("all text = %q", got)
	}
	if err := all.validate(); err != nil {
		t.Errorf("valid All effect rejected: %v", err)
	}
	if err := (MoveAemberToSupply{
		Amount: 1,
		All:    true,
		Target: Target{Kind: TargetThisCreature},
	}).validate(); err == nil {
		t.Error("Amount together with All should be rejected")
	}
}

func TestMoveAemberToSupplyGate(t *testing.T) {
	// Bind leaves the moved-from creature in context and reports whether Æmber
	// moved, so a Then wards it only when it held Æmber.
	g := NewGame("A", "B", 1)
	held := g.AddToBattleline(testCreature("held", 4), 0)
	g.AddAmberOn(held, 2)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	e := MoveAemberToSupply{
		Amount: 1,
		Target: Target{Kind: TargetThisCreature},
		Bind:   true,
	}
	if moved := e.resolveGate(ctx); !moved {
		t.Error("moving 1 of 2 Æmber should report progress")
	}
	if !ctx.HasIt || ctx.It != held {
		t.Errorf("ctx.It = %v (HasIt %v), want %v bound", ctx.It, ctx.HasIt, held)
	}
	if got := g.AmberOn(held); got != 1 {
		t.Errorf("after gated move = %d, want 1", got)
	}

	// A creature holding no Æmber moves none: the gate reports no progress but still
	// binds the chosen creature.
	g2 := NewGame("A", "B", 1)
	bare := g2.AddToBattleline(testCreature("bare", 4), 0)
	ctx2 := &EffectContext{
		Resolver:   g2,
		Controller: 0,
	}
	if moved := (MoveAemberToSupply{
		Amount: 1,
		Target: Target{Kind: TargetThisCreature},
		Bind:   true,
	}).
		resolveGate(
			ctx2,
		); moved {
		t.Error("moving from an empty creature should report no progress")
	}
	if !ctx2.HasIt || ctx2.It != bare {
		t.Errorf("ctx.It = %v (HasIt %v), want %v bound", ctx2.It, ctx2.HasIt, bare)
	}
}

func TestCaptureAllAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("drumble", 2), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}
	g.State.Aember[1] = 7

	e := CaptureAember{
		All:    true,
		Target: Target{Kind: TargetThisCreature},
		Source: Opponent,
	}
	if e.Text() != "{self} captures all your opponent's Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.AmberOn(src) != 7 {
		t.Errorf("captured = %d, want 7 (all of the opponent's pool)", g.AmberOn(src))
	}
	if g.Aember(1) != 0 {
		t.Errorf("opponent aember = %d, want 0", g.Aember(1))
	}
}

// TestCaptureAemberScaled captures once onto one creature, with Per scaling how
// much that one capturer takes (Yxili Marauder: 1 per friendly ready Mars
// creature), as opposed to Times repeating the capture onto fresh creatures.
func TestCaptureAemberScaled(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(marsCreature("yxili", 2), 0)
	g.AddToBattleline(marsCreature("ally", 3), 0) // two friendly Mars creatures -> Per = 2
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}
	g.State.Aember[1] = 5

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetThisCreature},
		Source: Opponent,
		Per: CardsInPlay{
			Player: Controller,
			Type:   Creature,
			House:  namedHouse(Mars),
		},
	}
	e.Resolve(ctx)
	if got := g.AmberOn(src); got != 2 {
		t.Errorf("captured onto self = %d, want 2 (1 x 2 Mars)", got)
	}
	if got := g.Aember(1); got != 3 {
		t.Errorf("opponent aember = %d, want 3 (5 - 2)", got)
	}
}

func TestCaptureAemberBy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("gate", 5), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}
	g.State.Aember[1] = 9

	e := CaptureAember{
		By:     AllBut(5),
		Target: Target{Kind: TargetThisCreature},
		Source: Opponent,
	}
	if want := "{self} captures all but 5 Æmber from your opponent"; e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	e.Resolve(ctx)
	if g.AmberOn(src) != 4 {
		t.Errorf("captured = %d, want 4", g.AmberOn(src))
	}
	if g.Aember(1) != 5 {
		t.Errorf("opponent aember = %d, want 5", g.Aember(1))
	}
}

func TestCaptureAemberText(t *testing.T) {
	this := Target{Kind: TargetThisCreature}
	enemy := Target{Kind: TargetChosenEnemyCreature}
	cases := []struct {
		name string
		e    CaptureAember
		want string
	}{
		{
			"friendly from opponent",
			CaptureAember{
				Amount: 1,
				Target: this,
				Source: Opponent,
			},
			"{self} captures 1 Æmber from your opponent",
		},
		{
			"friendly from own side",
			CaptureAember{
				Amount: 1,
				Target: this,
				Source: Controller,
			},
			"{self} captures 1 Æmber from your own side",
		},
		{
			"enemy from their own side",
			CaptureAember{
				Amount: 1,
				Target: enemy,
				Source: Opponent,
			},
			"an enemy creature captures 1 Æmber from their own side",
		},
		{
			"enemy from your pool",
			CaptureAember{
				Amount: 1,
				Target: enemy,
				Source: Controller,
			},
			"an enemy creature captures 1 Æmber from you",
		},
		{
			"all from own pool",
			CaptureAember{
				All:    true,
				Target: this,
				Source: Controller,
			},
			"{self} captures all your Æmber",
		},
	}
	for _, tc := range cases {
		if got := tc.e.Text(); got != tc.want {
			t.Errorf("%s: text = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestCaptureAemberValidate(t *testing.T) {
	this := Target{Kind: TargetThisCreature}
	if err := validateEffect(CaptureAember{
		Amount: 1,
		Source: Opponent,
	}); err == nil {
		t.Error("unset Target should fail validation")
	}
	if err := validateEffect(CaptureAember{
		Amount: 1,
		Target: this,
	}); err == nil {
		t.Error("unset Source should fail validation")
	}
	if err := validateEffect(CaptureAember{
		Amount: 1,
		Target: this,
		Source: Opponent,
	}); err != nil {
		t.Errorf("valid: %v, want nil", err)
	}
	both := CaptureAember{
		Amount: 1,
		By:     AllBut(5),
		Target: this,
		Source: Opponent,
	}
	if err := validateEffect(both); err == nil {
		t.Error("setting both Amount and By should fail validation")
	}
	lone := CaptureAember{
		Amount:   1,
		Target:   this,
		Source:   Opponent,
		Distinct: true,
	}
	if err := validateEffect(lone); err == nil {
		t.Error("Distinct without a Times should fail validation")
	}
}

// TestCaptureAemberDistinct covers Unguarded Camp: a repeated capture spreads
// across different creatures, and stops early once every eligible creature has
// already captured.
func TestCaptureAemberDistinct(t *testing.T) {
	spread := CaptureAember{
		Amount:   1,
		Target:   Target{Kind: TargetChosenFriendlyCreature},
		Source:   Opponent,
		Times:    ExcessCreatures{Player: Controller},
		Distinct: true,
	}
	want := "for each creature you have in excess of your opponent, a friendly " +
		"creature captures 1 Æmber from your opponent. Each creature cannot " +
		"capture more than 1 Æmber this way"
	if got := spread.Text(); got != want {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	b := g.AddToBattleline(testCreature("b", 3), 0)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.SetAember(1, 5)
	// The chooser always takes the first candidate offered, so without Distinct
	// every capture would pile onto a.
	g.SetChooser(0, &idQueueChooser{})

	spread.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})

	for _, id := range []LocalID{a, b, c} {
		if got := g.AmberOn(id); got != 1 {
			t.Errorf("creature %d captured %d, want 1 each", id, got)
		}
	}
	if got := g.Aember(1); got != 2 {
		t.Errorf("opponent pool = %d, want 2", got)
	}
}
func marsCreature(name string, power int) CardDefinition {
	return NewCard(name, Mars, Creature, Common, WithPower(power), WithTraits(Martian))
}

func TestCaptureAemberByEnemy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(marsCreature("command", 1), 0)
	g.AddToBattleline(marsCreature("m1", 3), 0) // two friendly Mars creatures -> Per = 2
	foe := g.AddToBattleline(testCreature("foe", 4), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}
	g.State.Aember[1] = 3

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Times: CardsInPlay{
			Player: Controller,
			Type:   Creature,
			House:  namedHouse(Mars),
		},
	}
	want := "for each friendly Mars creature, an enemy creature captures 1 Æmber from their own side"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	e.Resolve(ctx)
	// Two friendly Mars creatures, so the enemy captures twice; the default chooser
	// picks foe both times.
	if g.AmberOn(foe) != 2 {
		t.Errorf("captured on foe = %d, want 2", g.AmberOn(foe))
	}
	if g.Aember(1) != 1 {
		t.Errorf("opponent pool = %d, want 1 (3 - 2 captured)", g.Aember(1))
	}
}

func TestCaptureAemberByEnemyNoEnemies(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(marsCreature("command", 1), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}
	g.State.Aember[1] = 2

	// One friendly Mars creature but no enemy creatures: the loop finds nothing to
	// capture with and stops.
	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Times: CardsInPlay{
			Player: Controller,
			Type:   Creature,
			House:  namedHouse(Mars),
		},
	}
	e.Resolve(ctx)
	if g.Aember(1) != 2 {
		t.Errorf("opponent pool = %d, want 2 (nothing captured)", g.Aember(1))
	}
}

func TestCaptureAemberByEnemyDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(marsCreature("command", 1), 0)
	g.AddToBattleline(testCreature("foe1", 4), 1)
	g.AddToBattleline(testCreature("foe2", 5), 1) // two candidates, so the chooser is consulted
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}
	g.State.Aember[1] = 2

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Times: CardsInPlay{
			Player: Controller,
			Type:   Creature,
			House:  namedHouse(Mars),
		},
	}
	e.Resolve(ctx)
	if g.Aember(1) != 2 {
		t.Errorf("opponent pool = %d, want 2 (choice declined)", g.Aember(1))
	}
}

func TestCaptureAemberFromItsOpponent(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.AddToBattleline(testCreature("mine", 3), 0)
	theirs := g.AddToBattleline(testCreature("theirs", 3), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     mine,
		Controller: 0,
	}
	g.State.Aember[0] = 4
	g.State.Aember[1] = 4

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetEachCreature},
		Source: ItsOpponent,
	}
	if want := "each creature captures 1 Æmber from its opponent"; e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	e.Resolve(ctx)

	// Each creature drains the pool across from it, so both pools drop.
	if g.AmberOn(mine) != 1 || g.AmberOn(theirs) != 1 {
		t.Errorf("captured = %d/%d, want 1/1", g.AmberOn(mine), g.AmberOn(theirs))
	}
	if g.Aember(0) != 3 || g.Aember(1) != 3 {
		t.Errorf("pools = %d/%d, want 3/3", g.Aember(0), g.Aember(1))
	}
}

func TestCaptureAemberFromItsOpponentRejectsAShare(t *testing.T) {
	e := CaptureAember{
		All:    true,
		Target: Target{Kind: TargetEachCreature},
		Source: ItsOpponent,
	}
	if e.validate() == nil {
		t.Error("validate accepted a share of a per-capturer pool")
	}
}

func TestPlayerForItsOpponent(t *testing.T) {
	g := NewGame("A", "B", 1)
	theirs := g.AddToBattleline(testCreature("theirs", 3), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		It:         theirs,
	}
	if got := ctx.PlayerFor(ItsOpponent); got != 0 {
		t.Errorf("PlayerFor(ItsOpponent) = %d, want 0", got)
	}
}

// queueOptionChooser answers ChooseOption from a fixed queue of indices, so a test
// can drive a per-unit split choice.
type queueOptionChooser struct {
	FirstChooser
	opts []int
	i    int
}

func (q *queueOptionChooser) ChooseOption(_, _ string, _ []string) int {
	v := q.opts[q.i]
	q.i++
	return v
}

// TestCaptureFromAnyPlayer covers Crassosaurus's capture: the controller splits
// the capture across both pools, it caps at what the pools hold, and an untouched
// pool contributes nothing.
func TestCaptureFromAnyPlayer(t *testing.T) {
	e := CaptureFromAnyPlayer{Amount: 10}
	if got := e.Text(); got != "{self} captures 10 Æmber from any combination of players" {
		t.Errorf("text = %q", got)
	}
	if err := (CaptureFromAnyPlayer{}).validate(); err == nil {
		t.Error("a non-positive Amount should be rejected")
	}
	if err := e.validate(); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}

	// Rich pools: the controller chooses the split, taking some from each pool.
	t.Run("splits across both pools", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("crass", 4), 0)
		g.SetAember(0, 4)
		g.SetAember(1, 4)
		// opp, own, opp, own, opp -> 3 from the opponent, 2 from the controller.
		g.SetChooser(0, &queueOptionChooser{opts: []int{1, 0, 1, 0, 1}})
		CaptureFromAnyPlayer{Amount: 5}.
			Resolve(&EffectContext{
				Resolver:   g,
				Source:     src,
				Controller: 0,
			})
		if got := g.AmberOn(src); got != 5 {
			t.Errorf("captured = %d, want 5", got)
		}
		if g.Aember(0) != 2 || g.Aember(1) != 1 {
			t.Errorf("pools = %d/%d, want 2/1", g.Aember(0), g.Aember(1))
		}
	})

	// Short pools: the capture stops once both pools are empty, taking fewer than
	// the full Amount. The default chooser takes the controller's pool first.
	t.Run("captures less when the pools are short", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("crass", 4), 0)
		g.SetAember(0, 1)
		g.SetAember(1, 2)
		g.SetChooser(0, optionPicker{idx: 0}) // "your pool" whenever both hold Æmber
		CaptureFromAnyPlayer{Amount: 10}.
			Resolve(&EffectContext{
				Resolver:   g,
				Source:     src,
				Controller: 0,
			})
		if got := g.AmberOn(src); got != 3 {
			t.Errorf("captured = %d, want 3", got)
		}
		if g.Aember(0) != 0 || g.Aember(1) != 0 {
			t.Errorf("pools = %d/%d, want 0/0", g.Aember(0), g.Aember(1))
		}
	})

	// Only the opponent has Æmber: the empty controller pool contributes nothing.
	t.Run("only the opponent has Æmber", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("crass", 4), 0)
		g.SetAember(1, 3)
		CaptureFromAnyPlayer{Amount: 2}.
			Resolve(&EffectContext{
				Resolver:   g,
				Source:     src,
				Controller: 0,
			})
		if got := g.AmberOn(src); got != 2 {
			t.Errorf("captured = %d, want 2", got)
		}
		if g.Aember(1) != 1 {
			t.Errorf("opponent pool = %d, want 1", g.Aember(1))
		}
	})

	// Only the controller has Æmber: the empty opponent pool contributes nothing.
	t.Run("only the controller has Æmber", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("crass", 4), 0)
		g.SetAember(0, 3)
		CaptureFromAnyPlayer{Amount: 2}.
			Resolve(&EffectContext{
				Resolver:   g,
				Source:     src,
				Controller: 0,
			})
		if got := g.AmberOn(src); got != 2 {
			t.Errorf("captured = %d, want 2", got)
		}
		if g.Aember(0) != 1 {
			t.Errorf("controller pool = %d, want 1", g.Aember(0))
		}
	})

	// Both pools empty: the capture resolves to nothing.
	t.Run("no Æmber anywhere", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("crass", 4), 0)
		CaptureFromAnyPlayer{Amount: 5}.
			Resolve(&EffectContext{
				Resolver:   g,
				Source:     src,
				Controller: 0,
			})
		if got := g.AmberOn(src); got != 0 {
			t.Errorf("captured = %d, want 0", got)
		}
	})
}

// TestDistributeCapture covers Bring Low's capture: the controller captures all
// but five of the opponent's Æmber one at a time, choosing a friendly creature
// each time, and the loop stops when no friendly creature remains.
func TestDistributeCapture(t *testing.T) {
	e := DistributeCapture{
		By:     AllBut(5),
		Source: Opponent,
	}
	const want = "capture all but 5 Æmber from your opponent, " +
		"distributed among any number of friendly creatures"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	const wantSelf = "capture all but 5 Æmber from you, " +
		"distributed among any number of friendly creatures"
	if got := (DistributeCapture{
		By:     AllBut(5),
		Source: Controller,
	}).Text(); got != wantSelf {
		t.Errorf("self-source text = %q, want %q", got, wantSelf)
	}
	if err := (DistributeCapture{Source: Opponent}).validate(); err == nil {
		t.Error("a missing By should be rejected")
	}
	if err := (DistributeCapture{By: AllBut(5)}).validate(); err == nil {
		t.Error("a missing Source should be rejected")
	}
	if err := (DistributeCapture{
		By:     AllBut(5),
		All:    true,
		Source: Opponent,
	}).validate(); err == nil {
		t.Error("setting both By and All should be rejected")
	}
	if err := e.validate(); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}

	const wantAll = "capture all your opponent's Æmber, " +
		"distributed among any number of friendly creatures"
	if got := (DistributeCapture{
		All:    true,
		Source: Opponent,
	}).Text(); got != wantAll {
		t.Errorf("all-mode text = %q, want %q", got, wantAll)
	}
	const wantAllSelf = "capture all your Æmber, " +
		"distributed among any number of friendly creatures"
	if got := (DistributeCapture{
		All:    true,
		Source: Controller,
	}).Text(); got != wantAllSelf {
		t.Errorf("all-mode self text = %q, want %q", got, wantAllSelf)
	}
	if err := (DistributeCapture{
		All:    true,
		Source: Opponent,
	}).validate(); err != nil {
		t.Errorf("valid all-mode effect rejected: %v", err)
	}

	t.Run("all mode captures the whole pool", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 4), 0)
		b := g.AddToBattleline(testCreature("b", 4), 0)
		g.SetAember(1, 3)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{a, b, a}})
		(DistributeCapture{
			All:    true,
			Source: Opponent,
		}).Resolve(
			&EffectContext{
				Resolver:   g,
				Source:     a,
				Controller: 0,
			},
		)
		if g.Aember(1) != 0 {
			t.Errorf("opponent pool = %d, want 0", g.Aember(1))
		}
		if g.AmberOn(a) != 2 || g.AmberOn(b) != 1 {
			t.Errorf("captured a/b = %d/%d, want 2/1", g.AmberOn(a), g.AmberOn(b))
		}
	})

	t.Run("distributes onto chosen creatures", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 4), 0)
		b := g.AddToBattleline(testCreature("b", 4), 0)
		g.SetAember(1, 8) // all but 5 -> 3 captured
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{a, a, b}})
		e.Resolve(&EffectContext{
			Resolver:   g,
			Source:     a,
			Controller: 0,
		})
		if g.Aember(1) != 5 {
			t.Errorf("opponent pool = %d, want 5", g.Aember(1))
		}
		if g.AmberOn(a) != 2 || g.AmberOn(b) != 1 {
			t.Errorf("captured a/b = %d/%d, want 2/1", g.AmberOn(a), g.AmberOn(b))
		}
	})

	t.Run("captures nothing without a friendly creature", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetAember(1, 8)
		e.Resolve(&EffectContext{
			Resolver:   g,
			Source:     0,
			Controller: 0,
		})
		if g.Aember(1) != 8 {
			t.Errorf("opponent pool = %d, want 8", g.Aember(1))
		}
	})

	t.Run("captures nothing when the pool is at the remainder", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 4), 0)
		g.SetAember(1, 5)
		e.Resolve(&EffectContext{
			Resolver:   g,
			Source:     a,
			Controller: 0,
		})
		if g.Aember(1) != 5 || g.AmberOn(a) != 0 {
			t.Errorf("pool/on = %d/%d, want 5/0", g.Aember(1), g.AmberOn(a))
		}
	})
}

// TestCrassosaurusSelfPurge covers the composed Play ability: after the capture,
// the creature purges itself only when it captured fewer than 10 Æmber.
func TestCrassosaurusSelfPurge(t *testing.T) {
	play := Sequence{Effects: []Effect{
		CaptureFromAnyPlayer{Amount: 10},
		Conditional{
			Cond: Not{Cond: CountIs{
				Count:  AemberOnThis{},
				Is:     AtLeast,
				Amount: 10,
			}},
			Then: PurgeCreature{Target: Target{Kind: TargetThisCreature}},
		},
	}}

	// Short of 10: the creature captures what it can, then purges itself.
	t.Run("purges when it captured fewer than 10", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("crass", 4), 0)
		g.SetAember(1, 6)
		play.Resolve(&EffectContext{
			Resolver:   g,
			Source:     src,
			Controller: 0,
		})
		if !g.State.Purge[0].contains(src) {
			t.Error("Crassosaurus should have purged itself after capturing only 6")
		}
	})

	// A full 10: the creature keeps its captured Æmber and stays in play.
	t.Run("stays when it captured 10", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("crass", 4), 0)
		g.SetAember(1, 10)
		play.Resolve(&EffectContext{
			Resolver:   g,
			Source:     src,
			Controller: 0,
		})
		if g.State.Purge[0].contains(src) {
			t.Error("Crassosaurus should stay in play after capturing 10")
		}
		if got := g.AmberOn(src); got != 10 {
			t.Errorf("captured = %d, want 10", got)
		}
	})
}

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
	StealAember{Amount: 2}.Resolve(&EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	})
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
	CaptureAember{
		Amount: 2,
		Target: Target{Kind: TargetThisCreature},
		Source: Opponent,
	}.
		Resolve(&EffectContext{
			Resolver:   g,
			Source:     captor,
			Controller: 0,
		})
	// The creature captures 2, but P1 keeps their pool (the 2 came from supply).
	if g.AmberOn(captor) != 2 || g.State.Aember[1] != 3 {
		t.Errorf("after redirected capture: onCreature=%d opp=%d, want 2/3",
			g.AmberOn(captor), g.State.Aember[1])
	}
}

func TestTheftRedirectedToSupplyReader(t *testing.T) {
	g := NewGame("A", "B", 1)
	if _, ok := g.AemberTakenFromSupply(0); ok {
		t.Error("no Pixies in play, want not redirected")
	}
	g.AddToBattleline(posPixies(), 0)
	if _, ok := g.AemberTakenFromSupply(0); !ok {
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
	if _, ok := g.AemberTakenFromSupply(0); ok {
		t.Error("the card's own pool is not the watched pool when scoped to Opponent")
	}
	if _, ok := g.AemberTakenFromSupply(1); !ok {
		t.Error("opponent-scoped source should redirect the opponent's pool")
	}
}
