package engine

import (
	"strings"
	"testing"
)

func TestDrawModifierText(t *testing.T) {
	if got := drawModifierText(
		DrawModifier{Player: Controller, Amount: 1},
	); got != `During your "draw cards" step, refill your hand to 1 additional card.` {
		t.Errorf("controller text = %q", got)
	}
	if got := drawModifierText(
		DrawModifier{Player: Opponent, Amount: -1},
	); got != `During their "draw cards" step, your opponent refills their hand to 1 less card.` {
		t.Errorf("opponent text = %q", got)
	}
	if got := drawModifierText(
		DrawModifier{Player: EachPlayer, Amount: 2},
	); got != `During their "draw cards" step, each player refills their hand to 2 additional cards.` {
		t.Errorf("each-player text = %q", got)
	}
	if got := drawModifierText(DrawModifier{}); got != "" {
		t.Errorf("zero modifier text = %q, want empty", got)
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
		"refill your hand to 1 additional card",
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

	g.BeginTurn(0)
	g.EndTurn(0)

	if got := int(g.State.Hand[0].Count); got != HandSize+1 {
		t.Errorf("hand after draw = %d, want %d (one additional card)", got, HandSize+1)
	}
}
