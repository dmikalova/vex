package web

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests cover the client's own state — the pieces no engine rule touches:
// the hover preview, the overlays, the sidebar, and the menu items that hand the
// player off to one of them.

func TestHoverPreview(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	id := c.hand()[0]

	c.g.hoverCard(c.ctx, id)
	if !c.g.hasHover || c.g.hoverID != id {
		t.Fatalf("hovering card %d left hasHover=%v id=%d", id, c.g.hasHover, c.g.hoverID)
	}
	if !c.g.hoverLive() {
		t.Error("a card in hand should read as a live hover")
	}

	c.g.hoverClear(c.ctx)
	if c.g.hasHover {
		t.Error("hoverClear left the preview up")
	}
}

// A card that leaves the zones the client draws vanishes from the DOM without
// firing a leave, so the preview has to notice on its own that it is stale.
func TestHoverOfAGoneCardIsNotLive(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	id := c.hand()[0]
	c.g.hoverCard(c.ctx, id)

	c.manual()
	c.g.g.ManualMove(id, engine.ManualPurge)

	if c.g.hoverLive() {
		t.Error("a purged card still reads as a live hover")
	}
}

func TestRemainingKeyColors(t *testing.T) {
	c := newClient(t)
	all := []engine.KeyColor{engine.KeyColorRed, engine.KeyColorBlue, engine.KeyColorYellow}

	if got := c.g.remainingKeyColors(0); len(got) != len(all) {
		t.Fatalf("a player who has forged nothing has %v remaining, want all three", got)
	}

	c.manual()
	c.g.g.ManualForgeKeyColor(0, engine.KeyColorBlue)
	got := c.g.remainingKeyColors(0)
	if len(got) != 2 {
		t.Fatalf("after forging blue, %v remain, want two", got)
	}
	for _, col := range got {
		if col == engine.KeyColorBlue {
			t.Error("blue is still offered after being forged")
		}
	}
}

// Restarting is guarded by a confirmation, so a misclick does not throw a match
// away.
func TestRestartConfirmation(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	before := c.g.seed

	c.do(c.g.askRestart)
	if !c.g.confirmRestart {
		t.Fatal("askRestart did not arm the confirmation")
	}
	if c.g.seed != before {
		t.Error("askRestart restarted the match instead of asking")
	}

	c.do(c.g.cancelRestart)
	if c.g.confirmRestart {
		t.Error("cancelRestart left the confirmation armed")
	}

	c.do(c.g.askRestart)
	c.do(c.g.restart)
	if c.g.confirmRestart {
		t.Error("restarting left the confirmation armed")
	}
	if c.g.phase != phaseHouse {
		t.Errorf("the restarted match is at phase %v, want phaseHouse", c.g.phase)
	}
	if len(c.g.undo) != 0 {
		t.Errorf("the restarted match kept %d undo steps", len(c.g.undo))
	}
}

// Every one of these guards exists so a prompt the player must answer is not
// walked away from by a stray click behind it.
func TestBusyAndPromptGuards(t *testing.T) {
	tests := []struct {
		name string
		arm  func(g *game)
	}{
		{"busy", func(g *game) { g.busy = true }},
		{"choosing", func(g *game) { g.choosing = true }},
		{"choosingOption", func(g *game) { g.choosingOption = true }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t)
			c.startTurn()
			tt.arm(c.g)

			c.do(c.g.askRestart)
			if c.g.confirmRestart {
				t.Error("askRestart went through")
			}
			wasManual := c.g.g.Manual()
			c.do(c.g.toggleManual)
			if c.g.g.Manual() != wasManual {
				t.Error("toggleManual went through")
			}
		})
	}
}

// restart is guarded separately: it is reachable from the confirmation, which a
// prompt can appear behind.
func TestRestartIsGuarded(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	seed := c.g.seed
	c.g.busy = true
	c.do(c.g.restart)
	if c.g.seed != seed {
		t.Error("restart went through while an action was resolving")
	}
}

func TestOverlayToggles(t *testing.T) {
	c := newClient(t)
	c.startTurn()

	c.do(c.g.toggleKeys)
	if !c.g.keysOpen {
		t.Error("toggleKeys did not open the shortcut sheet")
	}
	c.do(c.g.closeKeys)
	if c.g.keysOpen {
		t.Error("closeKeys did not shut the shortcut sheet")
	}

	c.do(c.g.toggleMenu)
	if !c.g.menuOpen {
		t.Error("toggleMenu did not open the menu")
	}
	c.do(c.g.closeMenu)
	if c.g.menuOpen {
		t.Error("closeMenu did not shut the menu")
	}

	collapsed := c.g.sidebarCollapsed
	c.do(c.g.toggleSidebar)
	if c.g.sidebarCollapsed == collapsed {
		t.Error("toggleSidebar did not move the sidebar")
	}
}

// A menu item that opens something else closes the menu behind it, so the panel
// does not hang over the thing it just opened. An item the player is likely to
// repeat leaves it open.
func TestMenuItemsCloseBehindThemselves(t *testing.T) {
	for _, name := range []string{"manual", "restart", "keys"} {
		t.Run(name, func(t *testing.T) {
			c := newClient(t)
			c.startTurn()
			c.g.menuOpen = true
			switch name {
			case "manual":
				c.do(c.g.manualMenu)
			case "restart":
				c.do(c.g.restartMenu)
			case "keys":
				c.do(c.g.keysMenu)
			}
			if c.g.menuOpen {
				t.Error("the menu stayed open behind the item it opened")
			}
		})
	}

	for _, name := range []string{"undo", "redo"} {
		t.Run(name+" stays open", func(t *testing.T) {
			c := newClient(t)
			c.startTurn()
			c.g.menuOpen = true
			if name == "undo" {
				c.do(c.g.undoMenu)
			} else {
				c.do(c.g.redoMenu)
			}
			if !c.g.menuOpen {
				t.Error("the menu closed on an item meant to be repeated")
			}
		})
	}
}
