package engine

import (
	"slices"
	"testing"
)

func TestSwapCards(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 0)
	middle := g.AddToBattleline(testCreature("middle", 2), 0)
	right := g.AddToBattleline(testCreature("right", 2), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 2), 1)

	g.SwapCards(left, right)
	if got, want := g.Battleline(0), []LocalID{right, middle, left}; !slices.Equal(got, want) {
		t.Fatalf("battleline after swap = %v, want %v", got, want)
	}
	if got, want := g.Battleline(1), []LocalID{enemy}; !slices.Equal(got, want) {
		t.Fatalf("enemy battleline after friendly swap = %v, want %v", got, want)
	}

	g.SwapCards(left, enemy)
	if got, want := g.Battleline(0), []LocalID{right, middle, left}; !slices.Equal(got, want) {
		t.Fatalf("battleline after cross-battleline swap = %v, want %v", got, want)
	}
	if got, want := g.Battleline(1), []LocalID{enemy}; !slices.Equal(got, want) {
		t.Fatalf("enemy battleline after cross-battleline swap = %v, want %v", got, want)
	}
}

// TestSwapCardsAcrossZones covers a creature on the board trading places with one
// resting in a discard pile: the resting creature enters play in the board
// creature's slot while the board creature leaves to that discard pile.
func TestSwapCardsAcrossZones(t *testing.T) {
	t.Run("resting creature enters play in the board creature's slot", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		left := g.AddToBattleline(testCreature("left", 2), 0)
		host := g.AddToBattleline(testCreature("host", 2), 0)
		right := g.AddToBattleline(testCreature("right", 2), 0)
		reborn := g.AddToDiscard(testCreature("reborn", 3), 0)

		g.SwapCards(host, reborn)

		if got, want := g.Battleline(0), []LocalID{left, reborn, right}; !slices.Equal(got, want) {
			t.Fatalf("battleline after cross-zone swap = %v, want %v", got, want)
		}
		if !g.State.Discard[0].contains(host) {
			t.Errorf("host should have left play to the discard pile")
		}
		if g.State.Discard[0].contains(reborn) {
			t.Errorf("reborn should have left the discard pile")
		}
		if !g.State.Cards[reborn].Exhausted {
			t.Errorf("reborn should enter play exhausted")
		}
	})

	t.Run("the board creature sheds its Æmber to its opponent on the way out", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		host := g.AddToBattleline(testCreature("host", 2), 0)
		reborn := g.AddToDiscard(testCreature("reborn", 3), 0)
		g.State.Cards[host].Amber = 2

		g.SwapCards(host, reborn)

		if got := g.State.Aember[1]; got != 2 {
			t.Errorf("opponent pool = %d, want 2 (Æmber on the departing creature)", got)
		}
	})

	t.Run("does nothing when neither card is in play", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToDiscard(testCreature("a", 2), 0)
		b := g.AddToDiscard(testCreature("b", 2), 0)

		g.SwapCards(a, b)

		if !g.State.Discard[0].contains(a) || !g.State.Discard[0].contains(b) {
			t.Errorf("two resting cards should be left in the discard pile")
		}
	})

	t.Run("resolves regardless of argument order", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		host := g.AddToBattleline(testCreature("host", 2), 0)
		reborn := g.AddToDiscard(testCreature("reborn", 3), 0)

		g.SwapCards(reborn, host) // resting card passed first

		if got, want := g.Battleline(0), []LocalID{reborn}; !slices.Equal(got, want) {
			t.Fatalf("battleline = %v, want %v", got, want)
		}
	})

	t.Run("does nothing when the resting card is not in a discard pile", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		host := g.AddToBattleline(testCreature("host", 2), 0)
		inHand := g.AddToHand(testCreature("hand", 3), 0)

		g.SwapCards(host, inHand)

		if got, want := g.Battleline(0), []LocalID{host}; !slices.Equal(got, want) {
			t.Fatalf("battleline = %v, want %v (no swap)", got, want)
		}
		if !g.State.Hand[0].contains(inHand) {
			t.Errorf("card in hand should stay in hand")
		}
	})

	t.Run("does nothing when the in-play card is an artifact", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		art := g.AddArtifact(testArtifact("relic"), 0)
		reborn := g.AddToDiscard(testCreature("reborn", 3), 0)

		g.SwapCards(art, reborn)

		if !g.State.Artifacts[0].contains(art) {
			t.Errorf("artifact should stay in play")
		}
		if !g.State.Discard[0].contains(reborn) {
			t.Errorf("resting creature should stay in the discard pile")
		}
	})
}

func TestSwap(t *testing.T) {
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 2), 0)
	other := g.AddToBattleline(testCreature("other", 2), 0)
	right := g.AddToBattleline(testCreature("right", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: host, Controller: 0}
	e := Swap{With: Target{Kind: TargetChosenOtherFriendlyCreature}}

	if err := (Swap{}).validate(); err == nil {
		t.Fatal("an unset target should be rejected")
	}
	if err := validateEffect(e); err != nil {
		t.Fatalf("validate = %v", err)
	}
	if got, want := e.Text(), "swap this creature with another friendly creature in your battleline"; got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}
	if got, want := (Swap{FromContext: true}).Text(), "swap it with "+SelfName; got != want {
		t.Fatalf("FromContext text = %q, want %q", got, want)
	}

	e.Resolve(ctx)

	if got, want := g.Battleline(0), []LocalID{other, host, right}; !slices.Equal(got, want) {
		t.Fatalf("battleline after swap = %v, want %v", got, want)
	}
	if !ctx.HasIt || ctx.It != other {
		t.Fatalf("context card = %d (has %v), want swapped creature %d", ctx.It, ctx.HasIt, other)
	}
}

func TestSwapChosen(t *testing.T) {
	if got := (SwapChosen{}).Text(); got != "swap the positions of two creatures in a battleline" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// No creatures: nothing to choose, no panic.
	(SwapChosen{}).Resolve(ctx)

	a := g.AddToBattleline(testCreature("a", 2), 0)
	// One creature: no second creature to pick, battleline unchanged.
	(SwapChosen{}).Resolve(ctx)
	if got, want := g.Battleline(0), []LocalID{a}; !slices.Equal(got, want) {
		t.Fatalf("single-creature swap changed battleline: %v", got)
	}

	b := g.AddToBattleline(testCreature("b", 2), 0)
	c := g.AddToBattleline(testCreature("c", 2), 0)
	// The default chooser picks the first creature (a) then the first of the
	// remainder (b), swapping a and b.
	(SwapChosen{}).Resolve(ctx)
	if got, want := g.Battleline(0), []LocalID{b, a, c}; !slices.Equal(got, want) {
		t.Fatalf("battleline after swap = %v, want %v", got, want)
	}
}

// rearrangeChooser scripts the optional first pick from a queue (declining when it
// runs dry) and the second pick from another queue, falling back to the first
// candidate and declining an empty pool.
type rearrangeChooser struct {
	FirstChooser
	firsts  []LocalID
	seconds []LocalID
}

func (c *rearrangeChooser) ChooseCardOrDecline(_, _ string, _ []LocalID) (LocalID, bool) {
	if len(c.firsts) == 0 {
		return 0, false
	}
	id := c.firsts[0]
	c.firsts = c.firsts[1:]
	return id, true
}

func (c *rearrangeChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(c.seconds) > 0 {
		id := c.seconds[0]
		c.seconds = c.seconds[1:]
		return id, true
	}
	if len(cands) == 0 {
		return 0, false
	}
	return cands[0], true
}

func TestRearrangeBattleline(t *testing.T) {
	if got := (RearrangeBattleline{}).Text(); got != "rearrange the creatures in a player's battleline" {
		t.Errorf("text = %q", got)
	}

	// An empty board is a no-op: the loop runs zero times.
	empty := NewGame("A", "B", 1)
	(RearrangeBattleline{}).Resolve(&EffectContext{Resolver: empty, Controller: 0})

	// Swap two creatures, then stop.
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 2), 0)
	b := g.AddToBattleline(testCreature("b", 2), 0)
	c := g.AddToBattleline(testCreature("c", 2), 0)
	g.SetChooser(0, &rearrangeChooser{firsts: []LocalID{a}, seconds: []LocalID{b}})
	(RearrangeBattleline{}).Resolve(&EffectContext{Resolver: g, Controller: 0})
	if got, want := g.Battleline(0), []LocalID{b, a, c}; !slices.Equal(got, want) {
		t.Fatalf("battleline after one swap = %v, want %v", got, want)
	}

	// A lone creature in a battleline has no partner to swap with: the second pick
	// declines and the line is unchanged.
	g2 := NewGame("A", "B", 1)
	x := g2.AddToBattleline(testCreature("x", 2), 0)
	g2.AddToBattleline(testCreature("y", 2), 1)
	g2.SetChooser(0, &rearrangeChooser{firsts: []LocalID{x}})
	(RearrangeBattleline{}).Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if got, want := g2.Battleline(0), []LocalID{x}; !slices.Equal(got, want) {
		t.Fatalf("lone-creature line changed: %v, want %v", got, want)
	}
}

func TestMoveToFlankGameMethod(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 0)
	middle := g.AddToBattleline(testCreature("middle", 2), 0)
	right := g.AddToBattleline(testCreature("right", 2), 0)

	g.MoveToFlank(middle, true)
	if got, want := g.Battleline(0), []LocalID{left, right, middle}; !slices.Equal(got, want) {
		t.Fatalf("battleline after move right = %v, want %v", got, want)
	}

	g.MoveToFlank(middle, false)
	if got, want := g.Battleline(0), []LocalID{middle, left, right}; !slices.Equal(got, want) {
		t.Fatalf("battleline after move left = %v, want %v", got, want)
	}

	// A creature in no battleline leaves every line unchanged.
	g.MoveToFlank(LocalID(200), true)
	if got, want := g.Battleline(0), []LocalID{middle, left, right}; !slices.Equal(got, want) {
		t.Fatalf("battleline after moving an absent creature = %v, want %v", got, want)
	}
}

func TestMoveToFlank(t *testing.T) {
	if err := (MoveToFlank{}).validate(); err == nil {
		t.Fatal("an unset target should be rejected")
	}
	e := MoveToFlank{Target: Target{Kind: TargetTriggeringCreature}}
	if err := validateEffect(e); err != nil {
		t.Fatalf("validate = %v", err)
	}
	if got, want := e.Text(), "move it to either flank of its controller's battleline"; got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 2), 1)
	mover := g.AddToBattleline(testCreature("mover", 2), 1)
	c := g.AddToBattleline(testCreature("c", 2), 1)

	// Default chooser has no preference, so index 0 (the left flank) is taken.
	ctx := &EffectContext{Resolver: g, Controller: 0, It: mover, HasIt: true}
	e.Resolve(ctx)
	if got, want := g.Battleline(1), []LocalID{mover, a, c}; !slices.Equal(got, want) {
		t.Fatalf("battleline after move to left = %v, want %v", got, want)
	}

	// The effect's controller picks the right flank of the target's own line.
	g.SetChooser(0, optionPicker{idx: 1})
	e.Resolve(ctx)
	if got, want := g.Battleline(1), []LocalID{a, c, mover}; !slices.Equal(got, want) {
		t.Fatalf("battleline after move to right = %v, want %v", got, want)
	}

	// A target that selects nothing leaves the battleline unchanged.
	noTarget := MoveToFlank{Target: Target{Kind: TargetTriggeringCreature}}
	noTarget.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if got, want := g.Battleline(1), []LocalID{a, c, mover}; !slices.Equal(got, want) {
		t.Fatalf("battleline after empty target = %v, want %v", got, want)
	}

	// A triggering creature the preceding effect destroyed still selects, but it
	// has left the battleline: it is skipped before a flank is asked, so no move
	// happens and the line is unchanged.
	gone := g.AddToBattleline(testCreature("gone", 2), 1)
	g.removeFromPlay(gone)
	before := slices.Clone(g.Battleline(1))
	MoveToFlank{Target: Target{Kind: TargetTriggeringCreature}}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, It: gone, HasIt: true})
	if got := g.Battleline(1); !slices.Equal(got, before) {
		t.Fatalf("moving a creature that left play changed the battleline: %v", got)
	}
}

func TestMoveWithinBattleline(t *testing.T) {
	if err := (MoveWithinBattleline{}).validate(); err == nil {
		t.Fatal("an unset target should be rejected")
	}
	e := MoveWithinBattleline{Target: Target{Kind: TargetChosenEnemyCreature}}
	if got, want := e.Text(),
		"move an enemy creature anywhere in its controller's battleline"; got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}
	if err := e.validate(); err != nil {
		t.Fatalf("a valid target should pass validation: %v", err)
	}

	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 2), 1)
	b := g.AddToBattleline(testCreature("b", 2), 1)
	mover := g.AddToBattleline(testCreature("mover", 2), 1)

	// Player 0 chooses to move an enemy (player 1) creature. The default option
	// chooser takes position 0 (the left flank), and picks the last candidate
	// creature by default, so mover slides to the front and is left in context.
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !e.resolveGate(ctx) {
		t.Fatal("moving a creature should report a choice was made")
	}
	if !ctx.HasIt {
		t.Fatal("the moved creature should be left in context")
	}
	moved := ctx.It
	got := g.Battleline(1)
	if got[0] != moved {
		t.Fatalf("battleline = %v, want the moved creature %d on the left flank", got, moved)
	}
	_ = a
	_ = b
	_ = mover

	// A target that selects nothing leaves the line unchanged and reports no choice.
	before := slices.Clone(g.Battleline(1))
	if (MoveWithinBattleline{Target: Target{Kind: TargetTriggeringCreature}}).
		resolveGate(&EffectContext{Resolver: g, Controller: 0}) {
		t.Fatal("an empty target should report no choice")
	}
	if now := g.Battleline(1); !slices.Equal(now, before) {
		t.Fatalf("empty target changed the battleline: %v", now)
	}

	// The public Resolve is the one-line wrapper over resolveGate.
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
}

func testArtifact(name string, opts ...CardOption) CardDefinition {
	return NewCard(name, Brobnar, Artifact, Common, opts...)
}

func TestTurnIntoCreature(t *testing.T) {
	if err := (TurnIntoCreature{}).validate(); err == nil {
		t.Fatal("an unset target should be rejected")
	}
	e := TurnIntoCreature{Target: Target{Kind: TargetThisCreature}}
	if err := validateEffect(e); err != nil {
		t.Fatalf("validate = %v", err)
	}
	if got, want := e.Text(),
		"move it to a flank of your battleline as a creature"; got != want {
		t.Fatalf("text = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	existing := g.AddToBattleline(testCreature("existing", 2), 0)
	art := g.AddArtifact(testArtifact("art", WithArmor(3)), 0)
	// A converted card's power is its power counters; without any it would be a
	// 0-power creature the next boundary sweeps (ADR 0029), so give it one.
	g.AddPowerCounter(art, 1)

	// Default chooser has no preference, so index 0 (the left flank) is taken.
	ctx := &EffectContext{Resolver: g, Source: art, Controller: 0}
	e.Resolve(ctx)
	if g.TypeOf(art) != Creature {
		t.Fatalf("converted card should read as a creature, got %v", g.TypeOf(art))
	}
	if got, want := g.Battleline(0), []LocalID{art, existing}; !slices.Equal(got, want) {
		t.Fatalf("battleline after left conversion = %v, want %v", got, want)
	}
	if g.Artifacts(0) != nil && len(g.Artifacts(0)) != 0 {
		t.Fatalf("converted card should leave the artifact row: %v", g.Artifacts(0))
	}
	if got := int(g.State.Cards[art].ArmorRemaining); got != 3 {
		t.Fatalf("converted creature armor = %d, want 3", got)
	}

	// A second card converts to the right flank when its controller picks index 1.
	art2 := g.AddArtifact(testArtifact("art2"), 0)
	g.AddPowerCounter(art2, 1)
	g.SetChooser(0, optionPicker{idx: 1})
	TurnIntoCreature{Target: Target{Kind: TargetThisCreature}}.
		Resolve(&EffectContext{Resolver: g, Source: art2, Controller: 0})
	if got, want := g.Battleline(0),
		[]LocalID{art, existing, art2}; !slices.Equal(got, want) {
		t.Fatalf("battleline after right conversion = %v, want %v", got, want)
	}

	// A repeated use finds the card already a creature in the battleline, so it is
	// pulled from there (not the artifact row) and repositioned to the chosen flank
	// rather than duplicated.
	g.SetChooser(0, optionPicker{idx: 1})
	TurnIntoCreature{Target: Target{Kind: TargetThisCreature}}.
		Resolve(&EffectContext{Resolver: g, Source: art, Controller: 0})
	if got, want := g.Battleline(0),
		[]LocalID{existing, art2, art}; !slices.Equal(got, want) {
		t.Fatalf("battleline after repositioning = %v, want %v", got, want)
	}

	// A source no longer in play is a safe no-op.
	gone := g.AddArtifact(testArtifact("gone"), 0)
	g.State.Artifacts[0].remove(gone)
	before := slices.Clone(g.Battleline(0))
	TurnIntoCreature{Target: Target{Kind: TargetThisCreature}}.
		Resolve(&EffectContext{Resolver: g, Source: gone, Controller: 0})
	if got := g.Battleline(0); !slices.Equal(got, before) {
		t.Fatalf("converting a card not in play changed the battleline: %v", got)
	}
}

// TestTurnIntoCreatureForRemainderOfTurn covers the turn-scoped conversion
// (Animator): a chosen card becomes a creature only until end of turn, then the
// ready phase reverts it to an artifact in its controller's row while its power
// counters persist and its creature-only combat state is dropped.
func TestTurnIntoCreatureForRemainderOfTurn(t *testing.T) {
	// A duration other than RemainderOfPlayerTurn is rejected.
	if err := (TurnIntoCreature{
		Target:   Target{Kind: TargetThisCreature},
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Fatal("a non-turn-scoped duration should be rejected")
	}
	turnScoped := TurnIntoCreature{
		Target:   Target{Kind: TargetTheChosenCreature},
		Duration: RemainderOfPlayerTurn,
	}
	if err := validateEffect(turnScoped); err != nil {
		t.Fatalf("validate = %v", err)
	}
	// A chosen (non-self) card moves to its controller's battleline; the turn-scoped
	// clause is appended.
	if got, want := turnScoped.Text(),
		"move it to a flank of its controller's battleline as a creature "+
			"for the remainder of the turn"; got != want {
		t.Fatalf("chosen text = %q, want %q", got, want)
	}
	// A card animating itself keeps "your battleline".
	if got, want := (TurnIntoCreature{
		Target:   Target{Kind: TargetThisCreature},
		Duration: RemainderOfPlayerTurn,
	}).Text(),
		"move it to a flank of your battleline as a creature "+
			"for the remainder of the turn"; got != want {
		t.Fatalf("self text = %q, want %q", got, want)
	}
	// Versatile inserts "with versatile" before the turn-scoped clause.
	if got, want := (TurnIntoCreature{
		Target:    Target{Kind: TargetTheChosenCreature},
		Duration:  RemainderOfPlayerTurn,
		Versatile: true,
	}).Text(),
		"move it to a flank of its controller's battleline as a creature "+
			"with versatile for the remainder of the turn"; got != want {
		t.Fatalf("versatile text = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	art := g.AddArtifact(testArtifact("art", WithArmor(2)), 0)
	g.AddPowerCounter(art, 3)
	TurnIntoCreature{
		Target:   Target{Kind: TargetThisCreature},
		Duration: RemainderOfPlayerTurn,
	}.Resolve(&EffectContext{Resolver: g, Source: art, Controller: 0})
	if g.TypeOf(art) != Creature {
		t.Fatalf("card should read as a creature, got %v", g.TypeOf(art))
	}
	if !g.State.Cards[art].CreatureUntilTurnEnd {
		t.Fatal("a turn-scoped conversion should mark the card for revert")
	}
	// Mark creature-only state that the revert must clear.
	g.State.Cards[art].Damage = 1
	g.State.Cards[art].Stunned = true
	g.State.Cards[art].Enraged = true
	g.State.Cards[art].Warded = true

	g.readyPhase(0)

	if g.TypeOf(art) != Artifact {
		t.Fatalf("card should revert to an artifact at end of turn, got %v", g.TypeOf(art))
	}
	if !slices.Contains(g.Artifacts(0), art) {
		t.Fatalf("reverted card should be back in the artifact row: %v", g.Artifacts(0))
	}
	if slices.Contains(g.Battleline(0), art) {
		t.Fatal("reverted card should leave the battleline")
	}
	if got := int(g.State.Cards[art].PowerCounters); got != 3 {
		t.Fatalf("power counters should persist, got %d", got)
	}
	c := g.State.Cards[art]
	if c.CreatureUntilTurnEnd || c.Damage != 0 || c.Stunned || c.Enraged || c.Warded {
		t.Fatalf("revert should drop creature-only state: %+v", c)
	}
}

// TestTurnIntoCreatureGrantsVersatile covers the versatile grant (Animator): the
// animated creature gains versatile for the remainder of the turn, so it can be
// used as if in the active house.
func TestTurnIntoCreatureGrantsVersatile(t *testing.T) {
	g := NewGame("A", "B", 1)
	art := g.AddArtifact(testArtifact("art"), 0)
	TurnIntoCreature{
		Target:    Target{Kind: TargetThisCreature},
		Duration:  RemainderOfPlayerTurn,
		Versatile: true,
	}.Resolve(&EffectContext{Resolver: g, Source: art, Controller: 0})
	if !g.HasKeyword(art, Versatile) {
		t.Fatal("an animated creature with versatile should have the keyword")
	}
	g.readyPhase(0)
	if g.HasKeyword(art, Versatile) {
		t.Fatal("versatile should lift when the turn-scoped conversion reverts")
	}
}

func TestTypeOf(t *testing.T) {
	g := NewGame("A", "B", 1)

	// A printed artifact reads as an artifact until it converts.
	art := g.AddArtifact(testArtifact("art"), 0)
	if got := g.TypeOf(art); got != Artifact {
		t.Fatalf("printed artifact TypeOf = %v, want Artifact", got)
	}

	// A card in an upgrade chain reads as an Upgrade whatever its printed type,
	// from its attachment rather than a stored type.
	host := g.AddToBattleline(testCreature("host", 3), 0)
	creatureUpgrade := g.AddToBattleline(testCreature("worn", 2), 0)
	g.State.Cards[creatureUpgrade].HostPlus = upgradePlus(host)
	if got := g.TypeOf(creatureUpgrade); got != Upgrade {
		t.Fatalf("attached creature TypeOf = %v, want Upgrade", got)
	}
}
