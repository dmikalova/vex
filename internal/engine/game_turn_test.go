package engine

import (
	"strings"
	"testing"
)

// TestChooseHouseSkipsCardsRemovedMidWindow covers the AfterChooseHouse window
// when an earlier card's ability removes a later card from play: the snapshot is
// taken once, so the removed card must drop its own AfterChooseHouse trigger
// (ADR 0013) rather than fire it from the discard pile and strand in-play state
// (Æmber placed on a card no longer in play) — the invariant that caught it.
func TestChooseHouseSkipsCardsRemovedMidWindow(t *testing.T) {
	destroyer := NewCard("Purge Kin", Logos, Creature, Rare, WithPower(3),
		WithAbility(TriggerAfterChooseHouse, Conditional{
			Cond: ChoseHouse{House: Logos},
			Then: Destroy{Target: Target{Kind: TargetEachOtherFriendlyCreature}},
		}))
	hoarder := NewCard("Aember Sink", Logos, Creature, Rare, WithPower(3),
		WithAbility(TriggerAfterChooseHouse, Conditional{
			Cond: ChoseHouse{House: Logos},
			Then: PlaceAemberOnThis{Amount: 1},
		}))

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	// Order matters: the destroyer resolves first and removes the hoarder before
	// the loop reaches it.
	g.AddToBattleline(destroyer, 0)
	sink := g.AddToBattleline(hoarder, 0)

	if err := g.ChooseHouse(0, Logos); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if g.inPlay(sink) {
		t.Fatalf("hoarder should have been destroyed")
	}
	if core := g.State.Cards[sink]; core != (CardCore{}) {
		t.Errorf("destroyed card carries in-play state: %+v", core)
	}
	if err := g.InvariantError(); err != nil {
		t.Errorf("invariant violated: %v", err)
	}
}

// TestPlayerHasHouse covers the three outcomes: unknown deck houses treat every
// house as available, a declared house is available, and an undeclared one is not.
func TestPlayerHasHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	if !g.playerHasHouse(0, Mars) {
		t.Error("unset houses should treat every house as available")
	}
	g.SetPlayerHouses(0, []House{Mars, Logos, Untamed})
	if !g.playerHasHouse(0, Logos) {
		t.Error("a declared house should be available")
	}
	if g.playerHasHouse(0, Dis) {
		t.Error("an undeclared house should not be available")
	}
}

func TestDrawModifierText(t *testing.T) {
	if got := drawModifierText(
		DrawModifier{Player: Controller, Amount: 1},
	); got != `Your hand size is 1 more.` {
		t.Errorf("controller text = %q", got)
	}
	if got := drawModifierText(
		DrawModifier{Player: Opponent, Amount: -1},
	); got != `Your opponent's hand size is 1 less.` {
		t.Errorf("opponent text = %q", got)
	}
	if got := drawModifierText(
		DrawModifier{Player: EachPlayer, Amount: 2},
	); got != `Each player's hand size is 2 more.` {
		t.Errorf("each-player text = %q", got)
	}
	if got := drawModifierText(DrawModifier{}); got != "" {
		t.Errorf("zero modifier text = %q, want empty", got)
	}
	if got := drawModifierText(
		DrawModifier{Player: Opponent, Amount: -1, OnlyWhileOffFlank: true},
	); got != `While `+SelfName+` is not on a flank, your opponent's hand size is 1 less.` {
		t.Errorf("off-flank text = %q", got)
	}
	if got := drawModifierText(
		DrawModifier{Player: Controller, Amount: 2, OnlyWhileInCenter: true},
	); got != `While `+SelfName+` is in the center of the battleline, your hand size is 2 more.` {
		t.Errorf("in-center text = %q", got)
	}
	if got := drawModifierText(
		DrawModifier{
			Player: Controller,
			Amount: 1,
			Per:    InPlay{Player: Controller, Type: Creature, Trait: Sin},
		},
	); got != `For each friendly Sin creature your hand size is 1 more.` {
		t.Errorf("per text = %q", got)
	}
}

func TestDrawModifierAffects(t *testing.T) {
	self := DrawModifier{Player: Controller, Amount: 1}
	if !self.affects(0, 0) || self.affects(0, 1) {
		t.Error("Controller modifier should affect only its owner")
	}
	foe := DrawModifier{Player: Opponent, Amount: 1}
	if foe.affects(0, 0) || !foe.affects(0, 1) {
		t.Error("Opponent modifier should affect only the other player")
	}
	both := DrawModifier{Player: EachPlayer, Amount: 1}
	if !both.affects(0, 0) || !both.affects(0, 1) {
		t.Error("EachPlayer modifier should affect both players")
	}
}

func TestDrawModifierInCardRules(t *testing.T) {
	def := NewCard("Mother", Logos, Creature, Common, WithPower(5), WithDrawModifier(Controller, 1))
	if got := RenderCardRules(
		&def,
	); !strings.Contains(
		got,
		"Your hand size is 1 more.",
	) {
		t.Errorf("card rules = %q, want the draw-modifier line", got)
	}
}

func TestDrawStepModifier(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i < 10; i++ {
		g.AddToDeck(testCreature("d", 1), 0)
	}
	g.AddToBattleline(
		NewCard("mother", Logos, Creature, Common, WithPower(5), WithDrawModifier(Controller, 1)),
		0,
	)

	g.StartTurn(0)
	g.EndPlayPhase(0)

	if got := int(g.State.Hand[0].Count); got != HandSize+1 {
		t.Errorf("hand after draw = %d, want %d (one additional card)", got, HandSize+1)
	}
}

// TestDrawStepModifierPer covers a draw modifier that scales its amount by a count
// of the battleline (Greed: one extra card per friendly Sin creature).
func TestDrawStepModifierPer(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i < 12; i++ {
		g.AddToDeck(testCreature("d", 1), 0)
	}
	g.AddToBattleline(
		NewCard(
			"greed",
			Dis,
			Creature,
			Common,
			WithPower(4),
			WithTraits(Sin),
			WithDrawModifierPer(
				Controller,
				1,
				InPlay{Player: Controller, Type: Creature, Trait: Sin},
			),
		),
		0,
	)
	g.AddToBattleline(NewCard("sin", Dis, Creature, Common, WithPower(3), WithTraits(Sin)), 0)

	g.StartTurn(0)
	g.EndPlayPhase(0)

	// Two friendly Sin creatures raise the refill by two.
	if got := int(g.State.Hand[0].Count); got != HandSize+2 {
		t.Errorf("hand after draw = %d, want %d (two additional cards)", got, HandSize+2)
	}
}

// TestDrawModifierOffFlank covers the positional gate: the modifier applies only
// while its source is off a flank of its battleline (Streke).
func TestDrawModifierOffFlank(t *testing.T) {
	off := func() int {
		g := NewGame("A", "B", 1)
		streke := NewCard(
			"Streke", Dis, Creature, Common, WithPower(2),
			WithDrawModifierOffFlank(Opponent, -1),
		)
		g.AddToBattleline(testCreature("l", 3), 0)
		g.AddToBattleline(streke, 0) // buried in the middle: off a flank
		g.AddToBattleline(testCreature("r", 3), 0)
		return g.drawModifier(1)
	}
	if got := off(); got != -1 {
		t.Errorf("off-flank drawModifier = %d, want -1", got)
	}

	on := func() int {
		g := NewGame("A", "B", 1)
		streke := NewCard(
			"Streke", Dis, Creature, Common, WithPower(2),
			WithDrawModifierOffFlank(Opponent, -1),
		)
		g.AddToBattleline(streke, 0) // alone: on a flank
		return g.drawModifier(1)
	}
	if got := on(); got != 0 {
		t.Errorf("on-flank drawModifier = %d, want 0", got)
	}
}

// TestDrawModifierInCenter covers the positional gate: the modifier applies only
// while its source sits in the center of its battleline (Zenzizenzizenzic).
func TestDrawModifierInCenter(t *testing.T) {
	center := func() int {
		g := NewGame("A", "B", 1)
		zzz := NewCard(
			"Zenzizenzizenzic", Logos, Creature, Rare, WithPower(4),
			WithDrawModifierInCenter(Controller, 2),
		)
		g.AddToBattleline(testCreature("l", 3), 0)
		g.AddToBattleline(zzz, 0) // middle of an odd line: centered
		g.AddToBattleline(testCreature("r", 3), 0)
		return g.drawModifier(0)
	}
	if got := center(); got != 2 {
		t.Errorf("in-center drawModifier = %d, want 2", got)
	}

	offCenter := func() int {
		g := NewGame("A", "B", 1)
		zzz := NewCard(
			"Zenzizenzizenzic", Logos, Creature, Rare, WithPower(4),
			WithDrawModifierInCenter(Controller, 2),
		)
		g.AddToBattleline(zzz, 0) // on a flank with another creature: not centered
		g.AddToBattleline(testCreature("r", 3), 0)
		return g.drawModifier(0)
	}
	if got := offCenter(); got != 0 {
		t.Errorf("off-center drawModifier = %d, want 0", got)
	}
}

// Forgemaster Og reacts to any player forging a key — its own controller or the
// opponent — and empties the forging player's pool, so "that player" always means
// whoever just forged.
func TestForgemasterOgDrainsForger(t *testing.T) {
	og := NewCard("Forgemaster Og", Brobnar, Creature, Rare, WithPower(4),
		WithAbility(TriggerAfterPlayerForgesKey,
			LoseAember{Player: ThatPlayer, By: AllAember}))

	if got := RenderCardRules(&og); !strings.Contains(got,
		"After a player forges a key, that player loses all their Æmber.") {
		t.Fatalf("Og rules = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(og, 0)
	g.SetAember(0, 5)
	g.SetAember(1, 7)

	// The opponent (player 1) forges: Og, on player 0's side, drains player 1.
	g.forgeKeyFree(1)
	if got := g.State.Aember[1]; got != 0 {
		t.Errorf("forger's pool = %d, want 0 (Og drained it)", got)
	}
	if got := g.State.Aember[0]; got != 5 {
		t.Errorf("non-forger's pool = %d, want 5 (untouched)", got)
	}

	// Og's own controller forging is drained the same way.
	g.forgeKeyFree(0)
	if got := g.State.Aember[0]; got != 0 {
		t.Errorf("controller's pool after their forge = %d, want 0", got)
	}
}

// An ability that reacts only to the opponent forging (Forge Compiler) fires when
// the opponent forges and not when its own controller does.
func TestAfterOpponentForgesKeyFiresOnlyOnOpponentForge(t *testing.T) {
	watcher := NewCard("Watcher", Logos, Creature, Common, WithPower(3),
		WithAbility(TriggerAfterOpponentForgesKey, GainAember{Player: Controller, Amount: 2}))
	if got := RenderCardRules(&watcher); !strings.Contains(got,
		"After your opponent forges a key, gain 2 Æmber.") {
		t.Fatalf("Watcher rules = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(watcher, 0)

	// The opponent (player 1) forges: the ability on player 0's side fires.
	g.forgeKeyFree(1)
	if got := g.State.Aember[0]; got != 2 {
		t.Errorf("after opponent forge, controller pool = %d, want 2", got)
	}
	// The controller forging does not fire it.
	g.forgeKeyFree(0)
	if got := g.State.Aember[0]; got != 2 {
		t.Errorf("after own forge, controller pool = %d, want 2 (unchanged)", got)
	}
}

// TestForgeKeyOrdersMultipleReactions proves the "after you forge a key" reactions
// on two different cards are gathered into one window the forger orders, not fired
// one at a time in board order. A reversing chooser resolves the later card's
// reaction first, which it could not do if each fired in its own window.
func TestForgeKeyOrdersMultipleReactions(t *testing.T) {
	g := started(t)
	var log []string
	g.SetChooser(0, &countingReverseChooser{})
	g.AddToBattleline(
		testCreature("first", 3, WithAbility(TriggerAfterForgeKey, orderMark{&log, "first"})),
		0,
	)
	g.AddToBattleline(
		testCreature("second", 3, WithAbility(TriggerAfterForgeKey, orderMark{&log, "second"})),
		0,
	)

	g.forgeKeyFree(0)

	if len(log) != 2 || log[0] != "second" || log[1] != "first" {
		t.Errorf(
			"forge reaction order = %v, want [second first] (the forger orders one window)",
			log,
		)
	}
}
