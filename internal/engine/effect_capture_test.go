package engine

import "testing"

func TestCaptureAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 2

	e := CaptureAember{Amount: 3, Target: Target{Kind: TargetThisCreature}, Source: Opponent}
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

func TestMoveAemberToCommonSupplyEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("aubade", 4), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := MoveAemberToCommonSupply{Amount: 1, Target: Target{Kind: TargetThisCreature}}
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
	big := MoveAemberToCommonSupply{Amount: 5, Target: Target{Kind: TargetThisCreature}}
	big.Resolve(ctx)
	if got := g.AmberOn(src); got != 0 {
		t.Errorf("over-discard = %d, want 0", got)
	}

	// A chosen target renders by its own noun rather than {self}.
	chosen := MoveAemberToCommonSupply{Amount: 2, Target: Target{Kind: TargetChosenEnemyCreature}}
	if chosen.Text() != "move 2 Æmber from an enemy creature to the common supply" {
		t.Errorf("chosen text = %q", chosen.Text())
	}

	// validate rejects an unset target and a non-positive amount, accepts a valid one.
	if err := (MoveAemberToCommonSupply{Amount: 1}).validate(); err == nil {
		t.Error("unset target should be rejected")
	}
	if err := (MoveAemberToCommonSupply{Target: Target{Kind: TargetThisCreature}}).validate(); err == nil {
		t.Error("non-positive amount should be rejected")
	}
	if err := e.validate(); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}
}

func TestCaptureAllAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("drumble", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 7

	e := CaptureAember{All: true, Target: Target{Kind: TargetThisCreature}, Source: Opponent}
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
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 5

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetThisCreature},
		Source: Opponent,
		Per:    InPlay{Player: Controller, Type: Creature, House: Mars},
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
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 9

	e := CaptureAember{By: AllBut(5), Target: Target{Kind: TargetThisCreature}, Source: Opponent}
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
			CaptureAember{Amount: 1, Target: this, Source: Opponent},
			"{self} captures 1 Æmber from your opponent",
		},
		{
			"friendly from own side",
			CaptureAember{Amount: 1, Target: this, Source: Controller},
			"{self} captures 1 Æmber from your own side",
		},
		{
			"enemy from their own side",
			CaptureAember{Amount: 1, Target: enemy, Source: Opponent},
			"an enemy creature captures 1 Æmber from their own side",
		},
		{
			"enemy from your pool",
			CaptureAember{Amount: 1, Target: enemy, Source: Controller},
			"an enemy creature captures 1 Æmber from you",
		},
		{
			"all from own pool",
			CaptureAember{All: true, Target: this, Source: Controller},
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
	if err := validateEffect(CaptureAember{Amount: 1, Source: Opponent}); err == nil {
		t.Error("unset Target should fail validation")
	}
	if err := validateEffect(CaptureAember{Amount: 1, Target: this}); err == nil {
		t.Error("unset Source should fail validation")
	}
	if err := validateEffect(CaptureAember{Amount: 1, Target: this, Source: Opponent}); err != nil {
		t.Errorf("valid: %v, want nil", err)
	}
	both := CaptureAember{Amount: 1, By: AllBut(5), Target: this, Source: Opponent}
	if err := validateEffect(both); err == nil {
		t.Error("setting both Amount and By should fail validation")
	}
	lone := CaptureAember{Amount: 1, Target: this, Source: Opponent, Distinct: true}
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

	spread.Resolve(&EffectContext{Resolver: g, Controller: 0})

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
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 3

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Times:  InPlay{Player: Controller, Type: Creature, House: Mars},
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
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 2

	// One friendly Mars creature but no enemy creatures: the loop finds nothing to
	// capture with and stops.
	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Times:  InPlay{Player: Controller, Type: Creature, House: Mars},
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
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[1] = 2

	e := CaptureAember{
		Amount: 1,
		Target: Target{Kind: TargetChosenEnemyCreature},
		Source: Opponent,
		Times:  InPlay{Player: Controller, Type: Creature, House: Mars},
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
	ctx := &EffectContext{Resolver: g, Source: mine, Controller: 0}
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
	ctx := &EffectContext{Resolver: g, Controller: 0, It: theirs}
	if got := ctx.PlayerFor(ItsOpponent); got != 0 {
		t.Errorf("PlayerFor(ItsOpponent) = %d, want 0", got)
	}
}
