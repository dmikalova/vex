package engine

import (
	"strings"
	"testing"
)

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
