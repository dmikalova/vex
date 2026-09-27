package web

import (
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/cards"
	"github.com/dmikalova/vex/internal/engine"
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

// A log mention opens a printed-card preview and marks it as a log preview so it
// is placed to clear the line it was read from. Off-browser there is no window to
// measure, so the placement falls back to drawing over the sidebar.
func TestLogMentionPreview(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	var def *engine.CardDefinition
	for _, d := range c.g.defByName {
		def = d
		break
	}
	if def == nil {
		t.Fatal("no card definitions to preview")
	}

	c.g.setLogPreview(c.ctx, def)
	if !c.g.hoverInLog || c.g.hoverDef != def {
		t.Fatalf("setLogPreview left hoverInLog=%v hoverDef=%v", c.g.hoverInLog, c.g.hoverDef)
	}
	if !c.g.previewUp() {
		t.Error("a log-mention preview should read as up")
	}
	if !c.g.hoverOverSidebar {
		t.Error("with no window to measure, the preview should default over the sidebar")
	}
}

// On a touchscreen a tap raises a synthetic mouseenter just before its click, so
// onLogCardHover must ignore it once a touch has been seen — otherwise the enter
// opens the preview and the same tap's click toggles it shut, taking two taps to
// stick. With isTouch set the enter is a no-op, leaving the single tap
// (onLogCardTap) to open it.
func TestLogMentionHoverIgnoredOnTouch(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.isTouch = true

	c.g.onLogCardHover(c.ctx, app.Event{})
	if c.g.hoverDef != nil || c.g.previewUp() {
		t.Errorf("a touch mouseenter opened the preview: hoverDef=%v up=%v",
			c.g.hoverDef, c.g.previewUp())
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
	c.g.eng().ManualMove(id, engine.ManualPurge)

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

	up := c.g.eng().Register(
		engine.NewCard("Test Upgrade", testHouse, engine.Upgrade, engine.Common), c.g.active())
	c.g.eng().AttachUpgrade(host, up)

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

	buried := c.g.eng().Register(
		engine.NewCard("Buried", testHouse, engine.Creature, engine.Common), c.g.active())
	c.g.eng().AttachUnder(host, buried, false)

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

	buried := c.g.eng().Register(
		engine.NewCard("Buried", testHouse, engine.Creature, engine.Common), controller)
	c.g.eng().AttachUnder(host, buried, true)

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

// The opponent, who may not peek a facedown Under-card, still previews something
// when they hover it: a plain card back, not nothing and not the hidden face.
func TestOpponentPreviewsFacedownUnderAsCardBack(t *testing.T) {
	c := newClient(t)
	c.g.onCardBackHover(c.ctx, app.Event{})
	if !c.g.hoverBack || !c.g.previewUp() {
		t.Fatalf("hovering a hidden facedown under-card should preview a card back: "+
			"hoverBack=%v previewUp=%v", c.g.hoverBack, c.g.previewUp())
	}
	if c.g.hoverLive() {
		t.Error("a card-back preview must not read as a live hover of the hidden face")
	}
	c.g.hoverClear(c.ctx)
	if c.g.hoverBack || c.g.previewUp() {
		t.Error("leaving the tab should clear the card-back preview")
	}
}

func TestRemainingKeyColors(t *testing.T) {
	c := newClient(t)
	all := []engine.KeyColor{engine.KeyColorRed, engine.KeyColorBlue, engine.KeyColorYellow}

	if got := c.g.remainingKeyColors(0); len(got) != len(all) {
		t.Fatalf("a player who has forged nothing has %v remaining, want all three", got)
	}

	c.manual()
	c.g.eng().ManualForgeKeyColor(0, engine.KeyColorBlue)
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
	if len(c.g.rootMarks) != 0 {
		t.Errorf("the new match kept %d undo steps", len(c.g.rootMarks))
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

// Starting a new game while a prompt is still up — here the opening mulligan —
// re-deals cleanly: the abandoned match's session is dropped whole, so the new
// match settles at its own house choice rather than being clobbered by the old
// one finishing behind it.
func TestNewGameDuringMulliganPrompt(t *testing.T) {
	c := newBlankClient(t)
	c.g.dealMatch(testSeed)
	c.await("the first mulligan prompt", c.g.choosingOption)

	c.do(c.g.openSetup)
	if !c.g.awaitingSetup {
		t.Fatal("New game did not open the picker during the mulligan prompt")
	}
	set := cards.DeckSetNames()[0]
	c.do(func(ctx app.Context, _ app.Event) { c.g.pickSet(ctx, set) })
	c.do(func(ctx app.Context, _ app.Event) { c.g.pickSet(ctx, set) })

	// The new match runs its own opening mulligans and rests at the house choice.
	c.keepOpeningHands()
	if c.g.phase != phaseHouse {
		t.Errorf("the re-dealt match is at phase %v, want phaseHouse", c.g.phase)
	}
	if c.g.awaitingSetup {
		t.Error("the picker stayed open after the new match was dealt")
	}
}

// Opening the set picker is allowed from any state — including behind a prompt —
// because it only offers a new game and discards nothing until sets are confirmed;
// confirming then re-deals, which abandons the session the prompt belonged to.
// Manual mode stays toggleable mid-prompt too, so the player can reach the Cancel
// that backs out of an answerless prompt.
func TestOpenSetupIsAlwaysAllowed(t *testing.T) {
	tests := []struct {
		name string
		arm  func(c *client)
	}{
		{"a card prompt is up", func(c *client) {
			board := c.board()
			c.ask(board[0], "Choose a creature", false, board)
			c.await("the card prompt", c.g.choosing)
		}},
		{"an option prompt is up", func(c *client) {
			c.option(c.hand()[0], "Take them?", []string{"Yes", "No"})
			c.await("the option prompt", c.g.choosingOption)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t)
			c.manualTurn(testHouse)
			c.playFromHand(c.deal(testCreature))
			c.playFromHand(c.deal(testCreature))
			seed := c.g.seed
			tt.arm(c)

			c.do(c.g.openSetup)
			if !c.g.awaitingSetup {
				t.Error("openSetup did not open the picker")
			}
			if c.g.seed != seed {
				t.Error("opening the picker dealt a new match before sets were confirmed")
			}
			wasManual := c.g.eng().Manual()
			c.do(c.g.toggleManual)
			if c.g.eng().Manual() == wasManual {
				t.Error("manual mode could not be toggled while a prompt was up")
			}
		})
	}
}

// A tap on a card while an option prompt is up inspects it (a read-only lift) and
// leaves the prompt unanswered: the prompt is answered only on its own buttons.
// This checks the opening mulligan prompt — a plain option prompt — lifts a tapped
// hand card without settling the mulligan.
func TestTapDuringOptionPromptInspects(t *testing.T) {
	c := newBlankClient(t)
	c.g.dealMatch(testSeed)
	c.await("the first mulligan prompt", func() bool { return c.g.choosingOption() })

	hand := c.g.eng().Hand(c.g.active())
	if len(hand) == 0 {
		t.Fatal("no cards in hand during the mulligan prompt")
	}
	id := hand[0]
	c.do(func(ctx app.Context, _ app.Event) { c.g.selectHandID(ctx, id) })

	if !c.g.choosingOption() {
		t.Error("the tap answered the mulligan prompt; it should only inspect")
	}
	if !c.g.inspecting || !c.g.hasSel || c.g.sel != id {
		t.Errorf("the tap did not raise an inspect lift: inspecting=%v hasSel=%v sel=%v want %v",
			c.g.inspecting, c.g.hasSel, c.g.sel, id)
	}
}

// The n key turns a prompt down. At the opening mulligan the declining answer is
// labeled "Mulligan" rather than "No" — it sheds the hand — so n must answer that
// too, not fall through as a no-op.
func TestDenyKeyAnswersTheMulligan(t *testing.T) {
	c := newBlankClient(t)
	c.g.dealMatch(testSeed)
	c.await("the first mulligan prompt", func() bool { return c.g.choosingOption() })

	player := c.g.active()
	before := len(c.g.eng().Hand(player))
	c.press("n")
	c.await("the mulligan to be taken up", func() bool {
		return len(c.g.eng().Hand(player)) != before
	})

	if got := len(c.g.eng().Hand(player)); got != before-1 {
		t.Errorf("hand after n = %d, want %d (a mulligan draws one fewer)", got, before-1)
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
