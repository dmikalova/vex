package engine

import (
	"strings"
	"testing"
)

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
