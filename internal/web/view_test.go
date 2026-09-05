package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests draw the client and read the markup back. A view function that
// panics or quietly draws nothing is otherwise invisible to a test — the client
// is a screen, and what it puts on the screen is the thing worth pinning down.
//
// The tree is built without mounting it: go-app renders nested components all the
// way down whether or not there is a page behind them, so app.HTMLString over a
// freshly built tree is the whole screen.

// html draws the client and hands back the markup.
func (c *client) html() string {
	c.t.Helper()
	return app.HTMLString(c.g.Render())
}

// wants asserts that the drawn client shows each of the given fragments.
func (c *client) wants(what string, fragments ...string) {
	c.t.Helper()
	h := c.html()
	for _, f := range fragments {
		if !strings.Contains(h, f) {
			c.t.Errorf("%s does not show %q", what, f)
		}
	}
}

// lacks asserts that the drawn client shows none of the given fragments.
func (c *client) lacks(what string, fragments ...string) {
	c.t.Helper()
	h := c.html()
	for _, f := range fragments {
		if strings.Contains(h, f) {
			c.t.Errorf("%s still shows %q", what, f)
		}
	}
}

// Before the deal there is no match to draw, so the client puts up a placeholder
// rather than reaching into a game that is not there.
func TestDrawingBeforeTheDeal(t *testing.T) {
	c := newBlankClient(t)
	if h := c.html(); h != "<div></div>" {
		t.Errorf("the client before the deal draws %q", h)
	}
}

// The opening screen asks for a house and offers nothing else: End turn is
// withheld until a house is chosen.
func TestDrawingTheHousePrompt(t *testing.T) {
	c := newClient(t)
	c.wants("the house prompt", "control-dock", "house-pick")
	for _, h := range c.g.pickableHouses() {
		c.wants("the house prompt", h.String())
	}
	c.lacks("the house prompt", "End turn")
}

// Ordinary play draws the board, both players' bars, the log, and the way out of
// the turn.
func TestDrawingAPlayedTurn(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	c.wants("a played turn",
		"board-area", "score-pill", "log-list", "End turn", testCreature)
}

// A selected card is lifted off the board as a copy of itself carrying exactly
// the verbs that card has, while the dock keeps End turn.
func TestDrawingTheLiftedCard(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)

	c.g.selectHandID(c.ctx, id)
	c.wants("a card selected in hand", "card-focus", "Play", "Discard", "End turn")

	c.playFromHand(id)
	c.ownNextTurn(testHouse)
	c.g.selectBoardID(c.ctx, id)
	c.wants("a ready creature selected", "card-focus", "Reap", "Fight")
}

// Nothing selected lifts nothing.
func TestDrawingNoLiftedCard(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.lacks("no selection", "card-focus")
}

// Picking a flank is a question about the card being placed, so it is asked on
// the lifted copy of that card rather than in the dock.
func TestTheLiftAsksWhichFlank(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	c.g.selectHandID(c.ctx, c.deal(testCreature))
	c.do(c.g.play)
	c.wants("the flank prompt", "card-focus", "Left flank", "Right flank", "Cancel")
	c.lacks("the flank prompt", "End turn")
}

// Selecting another card while a flank is pending is a change of mind about which
// card to play, so the question goes away rather than being answered with the new
// card.
func TestSelectingAnotherCardTakesBackTheFlankQuestion(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	c.g.selectHandID(c.ctx, c.deal(testCreature))
	c.do(c.g.play)
	if c.g.phase != phaseFlank {
		t.Fatalf("the phase is %v, want phaseFlank", c.g.phase)
	}

	other := c.deal(testArtifact)
	c.g.selectHandID(c.ctx, other)
	if c.g.phase != phaseMain {
		t.Errorf("selecting another card left the phase at %v, want phaseMain", c.g.phase)
	}
	c.wants("the new selection", "Play", "Discard")
	c.lacks("the new selection", "Left flank")
}

// A card that cannot act is still lifted, and says why instead of offering verbs.
func TestALiftedCardSaysWhyItCannotAct(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	// The creature entered play exhausted, so this turn it can do nothing.
	c.g.selectBoardID(c.ctx, id)
	c.wants("an exhausted creature", "card-focus", "Cannot act")
	// The lift lies over the hand, whose cards may print a "Reap:" ability of
	// their own, so the no-verb check reads the lifted copy alone.
	h := c.html()
	from := strings.Index(h, `class="card-focus`)
	if from < 0 {
		t.Fatal("the client drew no lifted card")
	}
	if strings.Contains(h[from:], "Reap") {
		t.Error("an exhausted creature still offers a Reap verb")
	}
}

// A prompt takes the controls over, and an optional one draws the way to pass.
func TestDrawingAPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	answer := c.ask("Choose a creature to destroy", true, c.board())
	c.await("the prompt to go up", func() bool { return c.g.choosing })
	c.wants("a prompt", "Choose a creature to destroy", "Done", "prompt")
	c.lacks("a prompt", "End turn")

	c.do(c.g.declineChooser)
	<-answer
	c.await("the prompt to come down", func() bool { return !c.g.choosing })
	c.lacks("the answered prompt", "Choose a creature to destroy")
}

func TestDrawingAnOptionPrompt(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	out := make(chan int, 1)
	go func() { out <- c.g.chooser.ChooseOption("A Card", "Take them?", []string{"Yes", "No"}) }()
	c.await("the prompt to go up", func() bool { return c.g.choosingOption })

	c.wants("an option prompt", "Take them?", ">Yes<", ">No<")
	c.do(c.g.chooseOptionIdx(1))
	<-out
}

// A reap/fight/action prompt another card raised (Inspiration's "use a friendly
// creature") draws the standard use buttons — capitalised labels in their own
// colours — rather than the plain lowercase option list a generic prompt would.
func TestAUseVerbPromptDrawsTheStandardButtons(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	out := make(chan int, 1)
	labels := []string{"reap", "fight", "use its action"}
	go func() { out <- c.g.chooser.ChooseOption("A Card", "Choose how to use it", labels) }()
	c.await("the prompt to go up", func() bool { return c.g.choosingOption })

	c.wants("a use-verb prompt", ">Reap<", ">Fight<", ">Action<", "btn-warning", "btn-danger")
	c.lacks("a use-verb prompt", ">reap<", ">fight<")
	c.do(c.g.chooseOptionIdx(0))
	<-out
}

// An Upgrade attached to a creature draws as a peeking tab on the board rather
// than only as a rules-text line, so it reads as an attached card at a glance.
func TestDrawingAnUpgradeTab(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)

	up := c.g.g.Register(
		engine.NewCard("Test Upgrade", testHouse, engine.Upgrade, engine.Common), c.g.active())
	c.g.g.AttachUpgrade(host, up)

	c.wants("a creature with an attached upgrade", "card-host", "card-tabs--right", "card-tab")
	c.lacks("a creature with an attached upgrade", "card-tab--back")
}

// A faceup Under-card (Graft's rule) draws its own house colour, since it is
// visible to both players.
func TestDrawingARevealedUnderTab(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)

	buried := c.g.g.Register(
		engine.NewCard("Buried", testHouse, engine.Creature, engine.Common), c.g.active())
	c.g.g.AttachUnder(host, buried, false)

	c.wants("a creature with a faceup under-card", "card-tabs--left", "card-tab")
	c.lacks("a creature with a faceup under-card", "card-tab--back")
}

// A facedown Under-card the opponent may not peek draws as a plain card back,
// never the buried card's own face.
func TestDrawingAHiddenUnderTab(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)

	buried := c.g.g.Register(
		engine.NewCard("Buried", testHouse, engine.Creature, engine.Common), c.g.active())
	c.g.g.AttachUnder(host, buried, true)

	c.pass()
	c.manualTurn(testHouse)

	c.wants("the opponent's view of a facedown under-card", "card-tabs--left", "card-tab--back")
	c.lacks("the opponent's view of a facedown under-card", "Buried")
}

func TestDrawingTheFlankPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	id := c.deal(testCreature)
	c.g.selectHandID(c.ctx, id)
	c.do(c.g.play)
	if c.g.phase != phaseFlank {
		t.Fatalf("the phase is %v, want phaseFlank", c.g.phase)
	}
	c.wants("the flank prompt", "Left", "Right")
	c.lacks("the flank prompt", "End turn")
}

func TestDrawingTheEndTurnConfirmation(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.confirmEndTurn = true
	c.wants("an armed end turn", "Confirm end turn")
}

// The end-turn confirm wobbles the cards the player could still act with, so the
// warning points at exactly what it means. A playable card in hand jiggles only
// once the confirm is armed.
func TestTheEndTurnConfirmJigglesUsableCards(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.deal(testCreature)
	c.lacks("a playable hand card before the confirm", "card--jiggle")
	c.g.confirmEndTurn = true
	c.wants("a playable hand card once the confirm is armed", "card--jiggle")
}

// A face-up pile names its cards in the hover tip, so a zone reads like the
// upgrade title bars: just the names. A hidden zone stays a bare label.
func TestZoneTipNamesTheCardsInAFaceUpPile(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	def, ok := c.g.defByName[testCreature]
	if !ok {
		t.Fatalf("no card named %q", testCreature)
	}
	names := c.g.zoneNames(
		c.g.active(),
		"Discard",
		[]engine.LocalID{c.g.g.AddToDiscard(*def, c.g.active())},
	)
	if len(names) != 1 || names[0] != testCreature {
		t.Fatalf("discard tip listed %v, want [%s]", names, testCreature)
	}
	if got := c.g.zoneNames(c.g.active(), "Deck", []engine.LocalID{1}); got != nil {
		t.Errorf("a hidden deck leaked its names: %v", got)
	}
	if got := c.g.zoneNames(1-c.g.active(), "Hand", []engine.LocalID{1}); got != nil {
		t.Errorf("an opponent's hand leaked its names: %v", got)
	}
}

func TestDrawingTheSetPicker(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.do(c.g.restartMenu)
	c.wants("the set picker", "New game", "choose a set", "Cancel")
}

// A rejected click is reported in the dock, next to the control that rejected it.
func TestDrawingAStatusBanner(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.status = "Not enough Æmber"
	c.wants("a rejected action", "Not enough Æmber", "status")
}

func TestDrawingTheManualPanel(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.wants("manual mode", "Manual", "Add card")

	id := c.deal(testCreature)
	c.g.selectHandID(c.ctx, id)
	c.wants("a card selected in manual mode",
		"Move "+testCreature+" to:", "Deck bottom", "Archives", "Purge")

	c.playFromHand(id)
	c.g.selectBoardID(c.ctx, id)
	c.wants("a creature that entered exhausted", "Ready "+testCreature)
	c.do(c.g.manualReady)
	c.wants("a readied creature in manual mode", "Exhaust "+testCreature)
}

// The picker filters the whole pool by name, so a search narrows to what was
// typed.
func TestDrawingTheCardPicker(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.do(c.g.openPicker)
	c.wants("the open picker", pickerInputID)

	c.g.pickerQuery = testCreature
	c.wants("a searched picker", testCreature)
	c.lacks("a searched picker", "Bumpsy")
}

func TestDrawingTheForgePicker(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.do(c.g.manualForgeKey(c.g.active()))
	c.wants("the forge picker", "Red", "Blue", "Yellow")
}

func TestDrawingTheZoneViewer(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.zonesPlayer = 0
	c.wants("the zone viewer",
		"zones-panel", "Deck (", "Discard (", "Archives (", "Purge (")
}

func TestDrawingTheShortcutSheet(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.do(c.g.keysMenu)
	c.wants("the shortcut sheet", "Keyboard shortcuts", "keys-grid")
}

// Hovering a card puts its whole face up beside the board, which is the only way
// to read a card the board draws small.
func TestDrawingTheHoverPreview(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	c.g.hoverID, c.g.hasHover = id, true
	c.wants("a hovered card", "card-preview")
}

// With the sidebar away the log goes with it, but the controls float over the
// board instead, so the game stays playable.
func TestDrawingWithTheSidebarAway(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.do(c.g.toggleSidebar)
	c.wants("a collapsed sidebar",
		"app--sidebar-collapsed", "control-dock--floating", "sidebar-reveal")
	c.lacks("a collapsed sidebar", "log-list")
}

// A finished game replaces every control with the result rather than a modal
// over the position that ended it.
func TestDrawingAFinishedGame(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.g.State.Winner = 0
	c.g.settlePhase()
	c.wants("a finished game", "over-panel", "wins!", "New game")
	c.lacks("a finished game", "End turn")
}

// The log rules a line at each turn and phase, and marks the newest bubble so
// the eye lands on what just happened.
func TestDrawingTheLog(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.pass()

	c.wants("the log", "log-rule", "log-group", "log-group--new", "log-line")
}

// A restriction is named above the HUD, so the rule is read off the card that
// imposed it rather than discovered as a rejected click.
func TestDrawingRestrictionNotes(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.lacks("an unrestricted turn", "restriction log-card")
}

// The bars carry the counters a player watches, and manual mode adds the steppers
// that edit them.
func TestDrawingThePlayerBars(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.wants("a player bar", "aember.svg")
	c.lacks("an ordinary match", "amber-btn")

	c.manual()
	c.wants("manual mode", "amber-btn-plus", "amber-btn-minus", "chains.svg")
}

// Forged keys are drawn in the bar, in the colour they were forged.
func TestDrawingForgedKeys(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	me := c.g.active()
	c.do(c.g.manualForgeKey(me))
	c.do(c.g.pickForgeColor(engine.KeyColorRed))
	c.wants("a forged key", "score-keys", "key-red.svg", "key-btn")
}

// A card on the board draws under its own id, with its power on its face.
func TestDrawingACreatureOnTheBoard(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	c.wants("a creature on the board", boardCardID(id), testCreature, "power.svg")
}

// A card in hand draws under its own id and is a drag source, since dragging it
// onto the board is one of the two ways to play it.
func TestDrawingTheHand(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	if !c.g.playableFromHand(id) {
		t.Fatalf("%s is not playable from hand on its own house's turn", testCreature)
	}
	c.wants("a card in hand", handCardID(id), "draggable")
}

// A damaged creature shows the damage it is carrying; an undamaged one does not
// carry a zero.
func TestDrawingDamage(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)
	c.lacks("an undamaged creature", "damage.svg")

	c.g.g.State.Cards[id].Damage = 1
	c.wants("a damaged creature", "damage.svg")
}

// The menu is drawn behind its button, and undo is offered as unavailable rather
// than missing so the row does not move under the pointer.
func TestDrawingTheMenu(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.lacks("a closed menu", "menu-panel")

	c.do(c.g.toggleMenu)
	c.wants("the open menu",
		"menu-panel", "Undo", "Redo", "Manual mode", "New game", "Keyboard shortcuts")
}
