package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// TestDeckListShowsRoster checks the deck list draws from the retained roster: a
// deck icon per player, a three-column popover, and a real card in every slot.
func TestDeckListShowsRoster(t *testing.T) {
	c := newClient(t)
	for p := range 2 {
		if c.g.deckTip(p) == nil {
			t.Fatalf("player %d has no deck list icon", p)
		}
		if c.g.rosters[p].Empty() {
			t.Fatalf("player %d roster is empty", p)
		}
	}
	c.wants("the player bar", "deck-list.svg")

	pop := app.HTMLString(c.g.deckListPopover(0))
	for _, frag := range []string{"deck-list-cols", "deck-list-col", "deck-list-row"} {
		if !strings.Contains(pop, frag) {
			t.Errorf("deck list popover does not show %q", frag)
		}
	}
	first := c.g.rosters[0].Houses[0]
	if !strings.Contains(pop, first.House.String()) {
		t.Errorf("deck list popover omits house %s", first.House)
	}
	if name := first.Cards[0].Def.Name; !strings.Contains(pop, name) {
		t.Errorf("deck list popover omits card %q", name)
	}

	// A dealt roster names a real card in each of its 36 slots.
	slots := 0
	for _, hr := range c.g.rosters[0].Houses {
		for _, card := range hr.Cards {
			if card.Def.Name == "" {
				t.Error("roster slot has no card")
			}
			slots++
		}
	}
	if slots != 36 {
		t.Errorf("roster has %d slots, want 36", slots)
	}
}

// TestDeckListToggle checks a tap pins the deck list open and a second tap closes
// it, the touch path for a screen with no hover.
func TestDeckListToggle(t *testing.T) {
	c := newClient(t)
	if c.g.deckOpen != -1 {
		t.Fatalf("deck list starts open (%d)", c.g.deckOpen)
	}
	c.do(c.g.onDeckToggle(0))
	if c.g.deckOpen != 0 {
		t.Errorf("tap did not pin player 0's deck list open, got %d", c.g.deckOpen)
	}
	c.do(c.g.onDeckToggle(0))
	if c.g.deckOpen != -1 {
		t.Errorf("second tap did not close the deck list, got %d", c.g.deckOpen)
	}
}

// TestDeckListTypeOrder checks cards sort creatures first, then artifacts and
// upgrades, then one-shot Tactics last (ADR 0025).
func TestDeckListTypeOrder(t *testing.T) {
	ranks := []int{
		typeRank(engine.Creature),
		typeRank(engine.Artifact),
		typeRank(engine.Upgrade),
		typeRank(engine.Tactic),
	}
	for i := 1; i < len(ranks); i++ {
		if ranks[i-1] >= ranks[i] {
			t.Errorf("type order is not strictly increasing: %v", ranks)
		}
	}
}

// TestDeckListVisibleOpenFormat checks both lists are readable in the only format
// that exists today.
func TestDeckListVisibleOpenFormat(t *testing.T) {
	c := newClient(t)
	if !c.g.deckListVisible(0, 1) || !c.g.deckListVisible(1, 0) {
		t.Error("open-format deck lists are not visible to both players")
	}
}
