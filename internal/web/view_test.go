package web

import (
	"regexp"
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

// The lower (active-player) bar renders after the hand row, so a lifted card
// growing up from the hand covers the board's empty space rather than the bar
// (item 2 / ADR 0015). The bar is the last score-pill in the markup; the hand
// row is the only row labelled "Hand".
func TestTheLowerBarRendersBelowTheHand(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	h := c.html()
	hand := strings.Index(h, `row-label-zone">Hand`)
	bar := strings.LastIndex(h, "score-pill")
	if hand < 0 {
		t.Fatal("the played turn shows no hand row")
	}
	if bar < hand {
		t.Errorf(
			"the lower player bar renders before the hand row (bar at %d, hand at %d)",
			bar,
			hand,
		)
	}
}

// The action bar carries an inline Undo icon beside End turn, so a misplay is one
// click from being taken back. It rides with End turn: the opening house prompt
// has no End turn and so shows no Undo either.
func TestTheActionBarShowsAnInlineUndo(t *testing.T) {
	c := newClient(t)
	c.lacks("the house prompt", "end-turn-bar", "undo.svg")

	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.wants("a played turn", "end-turn-bar", "undo.svg", "End turn")
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

	// The opponent fields a creature during their turn, so our ready creature has
	// a legal fight target when we take the board again (Fight is offered only when
	// a target exists).
	c.pass()
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.pass()
	c.manualTurn(testHouse)

	c.g.selectBoardID(c.ctx, id)
	c.wants("a ready creature selected", "card-focus", "Reap", "Fight")
	// A lift that can still act is not dimmed.
	if lift := liftMarkup(t, c.html()); strings.Contains(lift, "card--dimmed") {
		t.Error("the lifted usable card was dimmed")
	}
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
	// The flank verbs sit on the lifted card; the dock shows the resting controls
	// greyed out rather than a blank box.
	c.wants("the flank dock", "end-turn-bar", "End turn", "disabled")
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
	// The verb button renders as >Reap<; a card's printed "Reap:" rules text
	// elsewhere on the board must not be mistaken for the button being offered.
	c.lacks("an exhausted creature", ">Reap<")
	// A lift with no verb it can take dims, so "cannot be used" reads at a glance.
	if lift := liftMarkup(t, c.html()); !strings.Contains(lift, "card--dimmed") {
		t.Error("the lifted unusable card was not dimmed")
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
	// The host reserves one tab's width of margin per attached card, so a neighbour
	// in the strip is pushed clear of the tab instead of covering it.
	c.wants("a host that reserves room for its upgrade tab", "--up-tabs:1")
}

// An attached upgrade that is a chooser candidate — one Destroy Them All may
// destroy — gets the targetable ring on its own tab, since the tab is the only
// thing to click for a card that shares its host's board slot.
func TestChoosingAnAttachedUpgrade(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)

	up := c.g.g.Register(
		engine.NewCard("Test Upgrade", testHouse, engine.Upgrade, engine.Common), c.g.active())
	c.g.g.AttachUpgrade(host, up)

	c.g.choosing = true
	c.g.chooserCandidates = []engine.LocalID{up}

	c.wants("an upgrade candidate", "card-tab--target")
}

// A power counter on a creature shows a +1 (or -1) token in its status row, with
// the count when more than one rides the card, so how many tokens sit on it is
// legible for the interactions that care.
func TestDrawingAPowerCounterToken(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)

	c.lacks("a creature with no counters", "power-counter-plus.svg", "power-counter-minus.svg")

	c.g.g.AddPowerCounter(host, 3)
	c.wants("a creature with three +1 counters", "power-counter-plus.svg", ">3<")
	c.lacks("a creature with +1 counters", "power-counter-minus.svg")

	c.g.g.AddPowerCounter(host, -5) // net -2
	c.wants("a creature at net -2", "power-counter-minus.svg", ">2<")
	c.lacks("a creature at net -2", "power-counter-plus.svg")
}

// Ward and Enrage show as condition tokens in a creature's status bar, the same
// way stun does, so a warded or enraged creature reads at a glance.
func TestDrawingWardAndEnrageTokens(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	c.lacks("an unwarded, unenraged creature", "ward.svg", "enrage.svg")

	c.g.g.State.Cards[id].Warded = true
	c.wants("a warded creature", "ward.svg")

	c.g.g.State.Cards[id].Enraged = true
	c.wants("an enraged creature", "enrage.svg")
}

// The armor shown on a creature is what it has left to absorb damage this turn,
// so it falls as hits land rather than staying pinned at the printed maximum.
func TestArmorShowsWhatIsLeftToAbsorb(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	// Whitespace is normalised so the count can be pinned to the shield icon it
	// labels — the markup is pretty-printed across lines, and the power stat can
	// carry the same number. The icon's attributes are matched order-independently
	// because go-app writes an element's attributes in map order, which is
	// randomised per render: what the test binds is the count to a shield icon, not
	// the order class and src happen to land in.
	norm := func(s string) string { return strings.Join(strings.Fields(s), "") }
	shield := regexp.MustCompile(`(\d+)<img[^>]*shield\.svg`)
	armorShown := func(h string) string {
		if m := shield.FindStringSubmatch(h); m != nil {
			return m[1]
		}
		return ""
	}

	c.g.g.State.Cards[id].ArmorRemaining = 5
	if got := armorShown(norm(c.html())); got != "5" {
		t.Errorf("a creature with armor intact showed %q armor remaining, want 5", got)
	}

	c.g.g.State.Cards[id].ArmorRemaining = 2
	if got := armorShown(norm(c.html())); got != "2" {
		t.Errorf("a creature after absorbing showed %q armor remaining, want 2", got)
	}
}

// Attached cards dim with an exhausted host and brighten when it readies, so an
// upgrade reads as spent alongside the creature that has already acted.
func TestAttachedTabsDimWithAnExhaustedHost(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host) // a creature enters play exhausted

	up := c.g.g.Register(
		engine.NewCard("Test Upgrade", testHouse, engine.Upgrade, engine.Common), c.g.active())
	c.g.g.AttachUpgrade(host, up)

	c.wants("an exhausted host's tabs", "card-tabs--dim")

	c.g.selectBoardID(c.ctx, host)
	c.do(c.g.manualReady)
	c.lacks("a ready host's tabs", "card-tabs--dim")
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
	// The dock has no flank question of its own, so it shows the resting controls
	// disabled instead of an empty box.
	c.wants("the flank dock", "end-turn-bar", "disabled")
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

// The end-turn confirm reveals its usable rows once when it arms and rearms the
// reveal after it disarms, so the strips scroll a jiggling card into view the
// moment moves-left is warned rather than every render while it stays armed.
func TestTheEndTurnConfirmScrollsUsableRowsOnce(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.deal(testCreature)

	c.g.scrollUsableRowsIntoView()
	if c.g.confirmScrolled {
		t.Fatal("the usable rows should not be scrolled before the confirm arms")
	}

	c.g.confirmEndTurn = true
	c.g.scrollUsableRowsIntoView()
	if !c.g.confirmScrolled {
		t.Fatal("arming the confirm should reveal the usable rows once")
	}

	c.g.confirmEndTurn = false
	c.g.scrollUsableRowsIntoView()
	if c.g.confirmScrolled {
		t.Fatal("disarming the confirm should rearm the reveal for next time")
	}
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
	// The player may review their own deck (sorted so its order stays hidden).
	own := c.g.zoneNames(
		c.g.active(),
		"Deck",
		[]engine.LocalID{c.g.g.AddToDeck(*def, c.g.active())},
	)
	if len(own) != 1 || own[0] != testCreature {
		t.Errorf("own deck roster listed %v, want [%s]", own, testCreature)
	}
	if got := c.g.zoneNames(1-c.g.active(), "Deck", []engine.LocalID{1}); got != nil {
		t.Errorf("an opponent's deck leaked its names: %v", got)
	}
	if got := c.g.zoneNames(1-c.g.active(), "Hand", []engine.LocalID{1}); got != nil {
		t.Errorf("an opponent's hand leaked its names: %v", got)
	}
}

// TestArchivesPillHasARosterPopover checks the Archives pill gets the same
// persistent roster popover the readable piles get — so a touchscreen can raise
// its label the way it raises Discard/Purge — while still leaking no card names,
// and that a truly hidden zone (an opponent's deck) gets no popover at all.
func TestArchivesPillHasARosterPopover(t *testing.T) {
	c := newClient(t)
	def, ok := c.g.defByName[testCreature]
	if !ok {
		t.Fatalf("no card named %q", testCreature)
	}
	p := c.g.active()
	id := c.g.g.AddToArchives(*def, p)
	roster := c.g.zoneRoster(p, "Archives", []engine.LocalID{id})
	if roster == nil {
		t.Fatal("the Archives pill has no roster popover, so it shows nothing on touch")
	}
	html := app.HTMLString(roster)
	if !strings.Contains(html, "zone-roster-label") || !strings.Contains(html, "Archives") {
		t.Errorf("the Archives popover does not name the zone: %s", html)
	}
	if strings.Contains(html, testCreature) {
		t.Errorf("the Archives popover leaked a card name: %s", html)
	}
	// An opponent's deck is hidden and not a labelled zone of its own, so it stays a
	// plain tip with no popover.
	if got := c.g.zoneRoster(1-p, "Deck", []engine.LocalID{id}); got != nil {
		t.Error("a hidden opponent deck got a roster popover")
	}
}

// TestZoneCountsWireTheRosterClamp checks the zone-count pills that open a roster
// popover carry the hover handler that clamps it on screen, so the rightmost pill
// (Purge) does not spill off the viewport edge.
func TestZoneCountsWireTheRosterClamp(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	html := app.HTMLString(app.Div().Body(c.g.zoneCounts(c.g.active())...))
	if !strings.Contains(html, "zone-has-roster") {
		t.Fatal("no zone pill opened a roster to clamp")
	}
}

func TestDrawingTheSetPicker(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.do(c.g.restartMenu)
	c.wants("the set picker", "New game", "choose a set", "Cancel")
}

// New game must work from the end-of-game panel too: over a finished game the
// picker takes the controls, rather than the win result shadowing it.
func TestNewGameFromTheEndOfGamePanel(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.g.State.Winner = 0
	c.g.phase = phaseOver
	c.wants("the end-of-game panel", "wins!", "New game")

	c.do(c.g.openSetup)
	c.wants("the set picker over a finished game", "choose a set")
	c.lacks("the set picker over a finished game", "wins!")
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

// The end-of-game banner names the winner in that player's own colour, so
// "Player 1 wins!" reads in the winner's colour rather than a fixed one. Every
// place a player's name is drawn wears the same per-player colour class.
func TestWinBannerNamesTheWinnerInTheirColour(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.g.State.Winner = 0
	c.g.phase = phaseOver
	c.wants("the win banner in player 0's colour", "player-name--p0", "wins!")

	c.g.g.State.Winner = 1
	c.wants("the win banner in player 1's colour", "player-name--p1", "wins!")
}

// A Deploy creature's placement is asked on the lifted copy of the card being
// placed — like the flank question — extended with the interior deploy options,
// laid out left-to-right as [left flank] [deploy left] [deploy right] [right
// flank] so the ends are the flanks and the middle slots it between creatures.
func TestDeployPromptLiftsPlacementButtons(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	line := c.board()

	// Stage the placement prompt on a lifted hand creature the way ChoosePosition
	// posts it while the play action is parked, so the render is exercised without a
	// background goroutine.
	deployed := c.deal(testCreature)
	c.g.selectHandID(c.ctx, deployed)
	c.g.choosingPosition = true
	c.g.positionLine = line
	c.g.positionRight = false

	c.wants("the deploy lift", "card-focus",
		"Left flank", "Deploy left", "Deploy right", "Right flank")
	// The dock stays out of it — the verbs are on the card, as with the flank prompt.
	c.lacks("the deploy dock", "Deploy left of a creature")

	h := c.html()
	prev := -1
	for _, label := range []string{"Left flank", "Deploy left", "Deploy right", "Right flank"} {
		at := strings.Index(h, label)
		if at < 0 {
			t.Fatalf("the deploy lift is missing %q", label)
		}
		if at <= prev {
			t.Errorf("%q is out of left-to-right order in the deploy lift", label)
		}
		prev = at
	}
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

// The end-of-turn standing draws all three key slots, colouring the ones forged
// and dimming the rest, rather than a plain "N keys" number.
func TestPlayerStandingDrawsThreeKeySlots(t *testing.T) {
	c := newClient(t)
	e := engine.PlayerStanding{
		Player:    0,
		Aember:    4,
		KeyColors: []engine.KeyColor{engine.KeyColorRed},
	}
	h := app.HTMLString(app.Div().Body(c.g.playerStandingSegments(e)...))
	if !strings.Contains(h, "key-red") {
		t.Error("the standing did not colour the forged key")
	}
	if n := strings.Count(h, "key-unforged"); n != 2 {
		t.Errorf("the standing drew %d unforged key slots, want 2", n)
	}
}

// The end-of-turn standing lights its Æmber amount with the check highlight when
// the player holds enough to forge their next key, and leaves it plain otherwise,
// so the log echoes the board's "Check!" glow.
func TestPlayerStandingHighlightsAemberAtCheck(t *testing.T) {
	c := newClient(t)
	cost := c.g.g.CurrentKeyCost(0)

	atCheck := engine.PlayerStanding{Player: 0, Aember: cost}
	h := app.HTMLString(app.Div().Body(c.g.playerStandingSegments(atCheck)...))
	if !strings.Contains(h, "log-aember") {
		t.Error("a standing at check did not highlight the Æmber amount")
	}

	below := engine.PlayerStanding{Player: 0, Aember: cost - 1}
	h = app.HTMLString(app.Div().Body(c.g.playerStandingSegments(below)...))
	if strings.Contains(h, "log-aember") {
		t.Error("a standing below check highlighted the Æmber amount")
	}
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
// carry a zero. The check is scoped to the creature's own stat line: a card's
// ability-glyph strip (ADR 0022) can carry a damage glyph too — a hand card that
// deals damage draws one — so a page-wide search would not tell a stat from a
// glyph.
func TestDrawingDamage(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)
	statHTML := func() string {
		return app.HTMLString(app.Div().Body(c.g.statLine(id)...))
	}
	if strings.Contains(statHTML(), "damage.svg") {
		t.Error("an undamaged creature still shows a damage stat")
	}

	c.g.g.State.Cards[id].Damage = 1
	if !strings.Contains(statHTML(), "damage.svg") {
		t.Error("a damaged creature does not show a damage stat")
	}
}

// An artifact shows no power or armor stat — only creatures carry those values.
// When a card turns into a creature (Auto-Legionary, ADR 0033) its power reads,
// so the stat line follows the card's current type, not its printed one.
func TestArtifactHidesPowerUntilItBecomesACreature(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testArtifact)
	c.playFromHand(id)
	statHTML := func() string {
		return app.HTMLString(app.Div().Body(c.g.statLine(id)...))
	}
	if strings.Contains(statHTML(), "power.svg") {
		t.Error("an artifact should not show a power stat")
	}

	c.g.g.PutIntoBattlelineAsCreature(id, true)
	if !strings.Contains(statHTML(), "power.svg") {
		t.Error("a card turned into a creature should show its power")
	}
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

// A triggered ability's glyph line closes its trigger glyph with a colon, the
// icon counterpart of the printed "Play:" before the effect.
func TestTriggerGlyphLineShowsAColon(t *testing.T) {
	line := iconLine(glyphLine{
		triggers: []string{"glyph-play"},
		glyphs:   []glyph{{asset: "damage"}},
	})
	if h := app.HTMLString(line); !strings.Contains(h, "card-glyph-trigger-colon") {
		t.Errorf("a trigger glyph line does not close with a colon: %s", h)
	}
}

// A triggerless line (a keyword-only strip) draws no colon.
func TestKeywordGlyphLineHasNoColon(t *testing.T) {
	line := iconLine(glyphLine{glyphs: []glyph{{asset: "kw-elusive"}}})
	if h := app.HTMLString(line); strings.Contains(h, "card-glyph-trigger-colon") {
		t.Errorf("a triggerless glyph line should draw no colon: %s", h)
	}
}
