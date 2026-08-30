package engine

import (
	"strings"
	"testing"
)

func TestReveal(t *testing.T) {
	if got := (Reveal{Player: Controller, House: Mars}).Text(); got != "reveal any number of Mars cards from your hand" {
		t.Errorf("house text = %q", got)
	}
	if got := (Reveal{Player: Opponent}).Text(); got != "reveal your opponent's hand" {
		t.Errorf("whole-hand text = %q", got)
	}

	g := NewGame("Alice", "Bob", 1)
	g.AddToHand(NewCard("Marauder", Mars, Creature, Common, WithPower(1)), 0)
	g.AddToHand(NewCard("Missile", Mars, Tactic, Common), 0)
	g.AddToHand(NewCard("Brute", Brobnar, Creature, Common, WithPower(1)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Revealing your Mars cards counts and logs both; the Brobnar card is untouched.
	Reveal{Player: Controller, House: Mars}.Resolve(ctx)
	if ctx.Revealed != 2 {
		t.Errorf("revealed = %d, want 2", ctx.Revealed)
	}
	if cnt := (CardsRevealed{}); cnt.Value(ctx) != 2 || cnt.CountText() != "card revealed this way" {
		t.Errorf("CardsRevealed = %d / %q", cnt.Value(ctx), cnt.CountText())
	}
	line := g.Log[len(g.Log)-1]
	if !strings.Contains(line, "Alice reveals") || !strings.Contains(line, "Marauder") || !strings.Contains(line, "Missile") {
		t.Errorf("log = %q, want the revealed Mars cards", line)
	}
	if strings.Contains(line, "Brute") {
		t.Errorf("the unrevealed Brobnar card must not be logged: %q", line)
	}

	// A whole-hand reveal (no house filter) shows every card, of any house.
	g2 := NewGame("A", "B", 1)
	g2.AddToHand(NewCard("x", Mars, Tactic, Common), 1)
	g2.AddToHand(NewCard("y", Brobnar, Tactic, Common), 1)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	Reveal{Player: Opponent}.Resolve(ctx2)
	if ctx2.Revealed != 2 {
		t.Errorf("whole-hand revealed = %d, want 2", ctx2.Revealed)
	}

	// Revealing nothing counts zero and writes no log line.
	g3 := NewGame("A", "B", 1)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	before := len(g3.Log)
	Reveal{Player: Controller, House: Mars}.Resolve(ctx3)
	if ctx3.Revealed != 0 {
		t.Errorf("revealed = %d, want 0", ctx3.Revealed)
	}
	if len(g3.Log) != before {
		t.Error("revealing nothing should not log")
	}
}
