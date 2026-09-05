package web

import (
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/cards"
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

// An attached Upgrade has no zone of its own — its host's Battleline slot is the
// only thing keeping it on the board — so hoverLive has to recognize it by
// walking the host's attachments rather than by finding it in a zone.
func TestHoverLiveForAnAttachedUpgrade(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)

	up := c.g.g.Register(
		engine.NewCard("Test Upgrade", testHouse, engine.Upgrade, engine.Common), c.g.active())
	c.g.g.AttachUpgrade(host, up)

	c.g.hoverCard(c.ctx, up)
	if !c.g.hoverLive() {
		t.Error("an attached upgrade should read as a live hover")
	}
}

// A faceup Under-card (Graft's rule) is visible to both players, so it hovers
// live regardless of whose turn is active.
func TestHoverLiveForAFaceupUnderCard(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)

	buried := c.g.g.Register(
		engine.NewCard("Buried", testHouse, engine.Creature, engine.Common), c.g.active())
	c.g.g.AttachUnder(host, buried, false)

	c.g.hoverCard(c.ctx, buried)
	if !c.g.hoverLive() {
		t.Error("a faceup under-card should read as a live hover")
	}
}

// A facedown Under-card (Masterplan, Jargogle) is only visible to the host's own
// controller — the master rulebook's Peekable rule — so it hovers live only on
// that player's own turn.
func TestHoverLiveForAFacedownUnderCardDependsOnPeek(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)
	controller := c.g.active()

	buried := c.g.g.Register(
		engine.NewCard("Buried", testHouse, engine.Creature, engine.Common), controller)
	c.g.g.AttachUnder(host, buried, true)

	c.g.hoverCard(c.ctx, buried)
	if !c.g.hoverLive() {
		t.Error("the controller should be able to hover their own facedown under-card")
	}

	c.pass()
	c.manualTurn(testHouse)
	if c.g.active() == controller {
		t.Fatal("test setup: the opponent should be active now")
	}
	if c.g.hoverLive() {
		t.Error("the opponent should not be able to hover a facedown under-card they cannot peek")
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

// New game opens the set picker in the action bar rather than dealing at once, so
// each player chooses a set; Cancel backs out and leaves the match untouched.
func TestNewGameSetPicker(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	before := c.g.seed

	c.do(c.g.openSetup)
	if !c.g.awaitingSetup {
		t.Fatal("New game did not open the set picker")
	}
	if c.g.seed != before {
		t.Error("opening the picker dealt a new match instead of asking")
	}

	c.do(c.g.cancelSetup)
	if c.g.awaitingSetup {
		t.Error("Cancel left the picker open")
	}
	if c.g.seed != before {
		t.Error("cancelling the picker dealt a new match")
	}

	// Choosing sets for both players deals the new match.
	c.do(c.g.openSetup)
	set := cards.DeckSetNames()[0]
	c.do(func(ctx app.Context, _ app.Event) { c.g.pickSet(ctx, set) })
	if c.g.awaitingSetup != true {
		t.Error("the picker closed before player 2 chose a set")
	}
	c.do(func(ctx app.Context, _ app.Event) { c.g.pickSet(ctx, set) })
	if c.g.awaitingSetup {
		t.Fatal("the picker stayed open after both players chose")
	}
	if c.g.phase != phaseHouse {
		t.Errorf("the new match is at phase %v, want phaseHouse", c.g.phase)
	}
	if len(c.g.undo) != 0 {
		t.Errorf("the new match kept %d undo steps", len(c.g.undo))
	}
}

// The same-sets shortcut deals a rematch straight from the previous game's sets.
func TestNewGameSameSets(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.setNames = [2]string{cards.DeckSetNames()[0], cards.DeckSetNames()[0]}

	c.do(c.g.openSetup)
	if !c.g.hasPrevSets() {
		t.Fatal("the picker did not carry the previous game's sets")
	}
	c.do(func(ctx app.Context, _ app.Event) { c.g.continueSameSets(ctx) })
	if c.g.awaitingSetup {
		t.Error("the same-sets shortcut left the picker open")
	}
	if c.g.phase != phaseHouse {
		t.Errorf("the rematch is at phase %v, want phaseHouse", c.g.phase)
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

			c.do(c.g.openSetup)
			if c.g.awaitingSetup {
				t.Error("openSetup went through")
			}
			wasManual := c.g.g.Manual()
			c.do(c.g.toggleManual)
			if c.g.g.Manual() != wasManual {
				t.Error("toggleManual went through")
			}
		})
	}
}

// Opening the set picker is guarded: it is refused while an action is resolving,
// so a match is never thrown away behind a prompt.
func TestNewGameIsGuarded(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	seed := c.g.seed
	c.g.busy = true
	c.do(c.g.openSetup)
	if c.g.awaitingSetup {
		t.Error("openSetup went through while an action was resolving")
	}
	if c.g.seed != seed {
		t.Error("a new match was dealt while an action was resolving")
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
