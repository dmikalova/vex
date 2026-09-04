package web

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests cover keyboard navigation: the rows the arrows walk, the nth card
// a number key picks, the cycle Tab steps through, and the layer Escape backs
// out of.

func TestNavRowsAreDrawnOrder(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	mine := c.deal(testCreature)
	c.playFromHand(mine)

	rows := c.g.navRows()
	if len(rows) != 5 {
		t.Fatalf("navRows has %d rows, want 5", len(rows))
	}
	if !containsID(rows[2], mine) {
		t.Error("the active player's battleline is not the third row")
	}
	if !containsID(rows[4], c.hand()[0]) {
		t.Error("the active player's hand is not the last row")
	}
}

func TestNavPosFindsTheSelection(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	rows := c.g.navRows()

	if row, col := c.g.navPos(rows); row != -1 || col != -1 {
		t.Errorf("with nothing selected navPos is (%d, %d), want (-1, -1)", row, col)
	}

	c.g.selectHandID(c.ctx, c.hand()[0])
	row, col := c.g.navPos(c.g.navRows())
	if row != 4 || col < 0 {
		t.Errorf("a card in hand is at (%d, %d), want the last row", row, col)
	}
}

// A selection that is not in any row — a card selected out of the zone viewer —
// reads as nowhere, so the arrows start over rather than walk from a card that
// is not on screen.
func TestNavPosOfACardInNoRow(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.selectZoneCard(c.ctx, c.g.g.Deck(c.g.active())[0])
	if row, _ := c.g.navPos(c.g.navRows()); row != -1 {
		t.Errorf("a card in no row is at row %d, want -1", row)
	}
}

// With nothing selected the arrows reach into the hand, which is the row a
// player reaches into most.
func TestArrowsStartInHand(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.press("ArrowLeft")
	if !c.g.hasSel || c.g.selKind != selHand {
		t.Fatalf("the first arrow left kind=%v has=%v", c.g.selKind, c.g.hasSel)
	}
	if c.g.sel != c.g.sortedHand(c.g.active())[0] {
		t.Error("the first arrow did not land on the first card in hand")
	}
}

// Walking along a row wraps at its ends, so one key reaches every card.
func TestArrowsWrapAlongARow(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	hand := c.g.sortedHand(c.g.active())
	c.g.selectHand(hand[0])

	c.press("ArrowLeft")
	if c.g.sel != hand[len(hand)-1] {
		t.Errorf("stepping left off the start landed on %d, want the last card %d",
			c.g.sel, hand[len(hand)-1])
	}
	c.press("ArrowRight")
	if c.g.sel != hand[0] {
		t.Errorf("stepping right off the end landed on %d, want the first card %d",
			c.g.sel, hand[0])
	}
}

// Stepping between rows skips the empty ones and holds its place along the row
// it steps into.
func TestArrowsSkipEmptyRows(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	mine := c.deal(testCreature)
	c.playFromHand(mine)

	c.g.selectHand(c.g.sortedHand(c.g.active())[0])
	c.press("ArrowUp")
	if c.g.sel != mine {
		t.Errorf("stepping up from hand landed on %d, want the battleline card %d",
			c.g.sel, mine)
	}
	c.press("ArrowDown")
	if c.g.selKind != selHand {
		t.Errorf("stepping back down landed on kind %v, want a card in hand", c.g.selKind)
	}
}

// Stepping off the top or bottom of the board has nowhere to go, so the
// selection stays put.
func TestArrowsStopAtTheEdges(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	first := c.g.sortedHand(c.g.active())[0]
	c.g.selectHand(first)
	c.press("ArrowDown")
	if c.g.sel != first {
		t.Error("stepping down off the hand moved the selection")
	}
}

func TestNumberKeysPickTheNthCard(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	hand := c.g.sortedHand(c.g.active())
	c.press("2")
	if c.g.sel != hand[1] {
		t.Errorf("2 selected %d, want the second card in hand %d", c.g.sel, hand[1])
	}
	// A number past the end of the row has no card to pick, so nothing moves.
	c.press("9")
	if len(hand) < 9 && c.g.sel != hand[1] {
		t.Error("a number past the end of the row moved the selection")
	}
}

// Tab walks the cards the player can actually act with and then parks on End
// turn, which is the last decision a turn has.
func TestTabWalksActionableCardsThenEndTurn(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	var actionable int
	for _, ids := range c.g.navRows() {
		for _, id := range ids {
			if c.g.tabbable(id) {
				actionable++
			}
		}
	}
	if actionable == 0 {
		t.Skip("the deal left the opening turn with nothing to act with")
	}

	for i := 0; i < actionable; i++ {
		c.press("Tab")
		if !c.g.hasSel {
			t.Fatalf("Tab %d landed on no card", i+1)
		}
		if !c.g.tabbable(c.g.sel) {
			t.Errorf("Tab stopped on %d, which has no move", c.g.sel)
		}
	}
	c.press("Tab")
	if !c.g.isEndTurnCursor() {
		t.Error("Tab past the last actionable card did not park on End turn")
	}
	if c.g.hasSel {
		t.Error("parking on End turn left a card selected")
	}
	c.shiftPress("Tab")
	if c.g.isEndTurnCursor() {
		t.Error("Shift+Tab did not step back off End turn")
	}
}

// The opponent's cards are read-only, so Tab passes them.
func TestTabPassesTheOpponentsCards(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.pass()
	c.manualTurn(testHouse)
	theirs := c.deal(testCreature)
	c.playFromHand(theirs)
	c.pass()
	c.manualTurn(testHouse)

	if c.g.tabbable(theirs) {
		t.Error("Tab stops on a card the active player cannot act with")
	}
}

func TestCycleIdx(t *testing.T) {
	tests := []struct {
		name         string
		n, cur, step int
		has          bool
		want         int
	}{
		{"forwards with no cursor starts at the first", 3, 0, 1, false, 0},
		{"backwards with no cursor starts at the last", 3, 0, -1, false, 2},
		{"forwards steps on", 3, 1, 1, true, 2},
		{"forwards wraps", 3, 2, 1, true, 0},
		{"backwards wraps", 3, 0, -1, true, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cycleIdx(tt.n, tt.cur, tt.has, tt.step); got != tt.want {
				t.Errorf("cycleIdx(%d, %d, %v, %d) = %d, want %d",
					tt.n, tt.cur, tt.has, tt.step, got, tt.want)
			}
		})
	}
}

func TestCycleID(t *testing.T) {
	ids := []engine.LocalID{7, 8, 9}
	tests := []struct {
		name string
		ids  []engine.LocalID
		cur  engine.LocalID
		step int
		want engine.LocalID
	}{
		{"an empty list has nothing to step to", nil, 7, 1, 0},
		{"an id outside the list starts at the first", ids, 99, 1, 7},
		{"an id outside the list starts at the last, backwards", ids, 99, -1, 9},
		{"steps on", ids, 8, 1, 9},
		{"wraps forwards", ids, 9, 1, 7},
		{"wraps backwards", ids, 7, -1, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cycleID(tt.ids, tt.cur, tt.step); got != tt.want {
				t.Errorf("cycleID(%v, %d, %d) = %d, want %d",
					tt.ids, tt.cur, tt.step, got, tt.want)
			}
		})
	}
}

// While a fight target is being picked, Tab cycles the enemy creatures the fight
// may still hit, and Enter commits the one it stopped on.
func TestTabAndEnterAnswerAFightTarget(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	attacker := c.deal(testCreature)
	c.playFromHand(attacker)

	c.pass()
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))

	c.pass()
	c.manualTurn(testHouse)
	c.g.selectBoardID(c.ctx, attacker)
	c.do(c.g.startFight)
	if c.g.phase != phaseFightTarget {
		t.Fatalf("phase is %v, want phaseFightTarget", c.g.phase)
	}

	c.press("Tab")
	if !c.g.hasCursor {
		t.Fatal("Tab did not put a cursor on a candidate")
	}
	target := c.g.promptCursor
	if !c.g.isSelected(target) {
		t.Error("the candidate the cursor is on does not draw as selected")
	}
	c.press("Enter")
	if c.g.phase != phaseMain {
		t.Errorf("Enter left the phase at %v, want phaseMain", c.g.phase)
	}
	if containsID(c.g.g.Battleline(1-c.g.active()), target) {
		t.Error("the candidate Enter committed to survived the fight")
	}
}

// Tab over the house prompt walks its buttons, and Enter presses the one it
// stopped on.
func TestTabAndEnterAnswerTheHousePrompt(t *testing.T) {
	c := newClient(t)
	if got := c.g.promptButtons(); got != len(c.g.pickableHouses()) {
		t.Fatalf("the house prompt offers %d buttons, want %d",
			got, len(c.g.pickableHouses()))
	}
	want := c.g.pickableHouses()[0]

	c.press("Tab")
	if !c.g.isButtonCursor(0) {
		t.Fatal("Tab did not stop on the first house button")
	}
	c.press("Enter")
	if c.g.g.State.ActiveHouse != want {
		t.Errorf("Enter chose %v, want %v", c.g.g.State.ActiveHouse, want)
	}
	if c.g.hasBtnCursor {
		t.Error("the answered prompt left its button cursor up")
	}
}

// The flank prompt is two buttons, which Tab walks like any other.
func TestTabAnswersTheFlankPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	second := c.deal(testCreature)
	c.g.selectHandID(c.ctx, second)
	c.do(c.g.play)
	if c.g.promptButtons() != 2 {
		t.Fatalf("the flank prompt offers %d buttons, want 2", c.g.promptButtons())
	}
	c.press("Tab")
	c.press("Enter")
	if c.board()[0] != second {
		t.Errorf("the left flank holds %d, want %d", c.board()[0], second)
	}
}

// pressButton has nothing to press when Tab has not parked on a button.
func TestPressButtonWithNoCursor(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	if c.g.pressButton(c.ctx) {
		t.Error("pressButton acted with no cursor on a button")
	}
}

// Escape backs out one layer per press, innermost first, so it never does more
// than the player expects.
func TestEscapeBacksOutOneLayerAtATime(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.selectHandID(c.ctx, c.hand()[0])
	c.g.confirmEndTurn = true
	c.g.confirmRestart = true
	c.g.zonesPlayer = 0
	c.g.pickerOpen = true
	c.g.menuOpen = true
	c.g.keysOpen = true

	layers := []struct {
		name string
		open func(g *game) bool
	}{
		{"the shortcut sheet", func(g *game) bool { return g.keysOpen }},
		{"the menu", func(g *game) bool { return g.menuOpen }},
		{"the card picker", func(g *game) bool { return g.pickerOpen }},
		{"the zone viewer", func(g *game) bool { return g.zonesPlayer >= 0 }},
		{"the restart confirmation", func(g *game) bool { return g.confirmRestart }},
		{"the end-turn confirmation", func(g *game) bool { return g.confirmEndTurn }},
		{"the selection", func(g *game) bool { return g.hasSel }},
	}
	for _, layer := range layers {
		if !layer.open(c.g) {
			t.Fatalf("%s was already closed before its turn to be dismissed", layer.name)
		}
		c.press("Escape")
		if layer.open(c.g) {
			t.Fatalf("Escape did not dismiss %s", layer.name)
		}
	}
	for _, layer := range layers {
		if layer.open(c.g) {
			t.Errorf("%s came back open", layer.name)
		}
	}
}

// The key-forge picker is a layer of its own, between the confirmations and the
// prompts.
func TestEscapeClosesTheForgePicker(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.do(c.g.manualForgeKey(0))
	if c.g.forgingKey != 0 {
		t.Fatalf("the forge picker is at %d, want player 0", c.g.forgingKey)
	}
	c.press("Escape")
	if c.g.forgingKey != -1 {
		t.Error("Escape did not close the forge picker")
	}
}

// z walks the zone viewer through both players and then closed, so one key
// reaches every pile without a shortcut per player.
func TestZoneViewerCycles(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	me := c.g.active()

	c.press("z")
	if c.g.zonesPlayer != me {
		t.Errorf("the first z opened player %d, want the active player %d",
			c.g.zonesPlayer, me)
	}
	c.press("z")
	if c.g.zonesPlayer != 1-me {
		t.Errorf("the second z opened player %d, want the opponent %d",
			c.g.zonesPlayer, 1-me)
	}
	c.press("z")
	if c.g.zonesPlayer != -1 {
		t.Errorf("the third z left the viewer at %d, want it closed", c.g.zonesPlayer)
	}
}

// The view-only keys work while a prompt is up, because they race with nothing.
func TestViewKeysWorkUnderAPrompt(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.choosing = true

	c.press("?")
	if !c.g.keysOpen {
		t.Error("? did not open the shortcut sheet under a prompt")
	}
	collapsed := c.g.sidebarCollapsed
	c.press("h")
	if c.g.sidebarCollapsed == collapsed {
		t.Error("h did not move the sidebar under a prompt")
	}
}

// A key with no game behind it is a no-op rather than a crash: the listener is
// installed before the first deal finishes.
func TestKeysBeforeTheDeal(t *testing.T) {
	c := newBlankClient(t)
	c.g.onKey(c.ctx, "p", false)
}

// Space confirms the Tab cursor when there is one; with none it takes the
// affirmative move for whatever is selected.
func TestSpaceAffirmsTheSelection(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.g.selectHandID(c.ctx, id)
	c.press(" ")
	if !containsID(c.board(), id) {
		t.Error("Space did not play the selected card in hand")
	}
}

// Playing from the keyboard hands the selection on to the next card there is
// still something to do with, so a run of plays is one key each.
func TestKeyboardPlayAdvancesTheSelection(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	first := c.deal(testCreature)
	c.deal(testCreature)
	c.g.selectHand(first)
	c.press("p")
	if c.g.phase == phaseFlank {
		c.press("r")
	}
	if !c.g.hasSel {
		t.Fatal("the selection was not handed on after a keyboard play")
	}
	if c.g.sel == first {
		t.Error("the selection stayed on the card that was played")
	}
	if !c.g.usableFromHand(c.g.sel) {
		t.Errorf("the selection landed on %d, which has no move", c.g.sel)
	}
}

// A mouse play does not move the selection: a player who just let go of a card
// would find it jumping somewhere else surprising.
func TestMousePlayDoesNotAdvanceTheSelection(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.g.selectHandID(c.ctx, id)
	c.do(c.g.play)
	if c.g.hasSel {
		t.Error("a mouse play handed the selection on")
	}
}
