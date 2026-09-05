package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file is keyboard navigation of the board: the rows the arrow keys walk,
// the nth card a number key picks, and the cycle Tab steps through — over the
// cards on the board, or over a prompt's candidates while one is up.

// navRows lists the card rows in the order they are drawn, top to bottom: the
// opponent's artifacts and battleline, then the active player's battleline,
// artifacts, and hand. The keyboard walks the board the way it looks.
func (g *game) navRows() [][]engine.LocalID {
	p, opp := g.active(), 1-g.active()
	return [][]engine.LocalID{
		g.sortedArtifacts(opp),
		g.g.Battleline(opp),
		g.g.Battleline(p),
		g.sortedArtifacts(p),
		g.sortedHand(p),
	}
}

// navPos locates the selection in rows, or (-1, -1) when nothing is selected.
func (g *game) navPos(rows [][]engine.LocalID) (row, col int) {
	if !g.hasSel {
		return -1, -1
	}
	for r, ids := range rows {
		if c := indexOfID(ids, g.sel); c >= 0 {
			return r, c
		}
	}
	return -1, -1
}

// navSelect selects the card at row/col. The hand goes through the hand's own
// selection path, so the action bar sees a card in hand rather than one in play.
func (g *game) navSelect(ctx app.Context, rows [][]engine.LocalID, row, col int) {
	if row < 0 || row >= len(rows) {
		return
	}
	ids := rows[row]
	if col < 0 || col >= len(ids) {
		return
	}
	if row == len(rows)-1 {
		g.selectHandID(ctx, ids[col])
		return
	}
	g.selectBoardID(ctx, ids[col])
}

// selectNth selects the nth card (1-based) of the row the selection is already
// in, defaulting to the hand when nothing is selected — the row a player reaches
// into most.
func (g *game) selectNth(ctx app.Context, n int) {
	rows := g.navRows()
	row, _ := g.navPos(rows)
	if row < 0 {
		row = len(rows) - 1
	}
	g.navSelect(ctx, rows, row, n-1)
}

// moveSel walks the selection one card along its row (dCol, wrapping at the
// row's ends) or one row up or down (dRow), skipping empty rows and holding its
// place along the row it steps into. With nothing selected it starts at the
// first card in hand.
func (g *game) moveSel(ctx app.Context, dRow, dCol int) {
	rows := g.navRows()
	row, col := g.navPos(rows)
	if row < 0 {
		g.navSelect(ctx, rows, len(rows)-1, 0)
		return
	}
	if dCol != 0 {
		ids := rows[row]
		g.navSelect(ctx, rows, row, ((col+dCol)%len(ids)+len(ids))%len(ids))
		return
	}
	for r := row + dRow; r >= 0 && r < len(rows); r += dRow {
		if len(rows[r]) == 0 {
			continue
		}
		g.navSelect(ctx, rows, r, min(col, len(rows[r])-1))
		return
	}
}

// tabSel moves on to the next (step 1) or previous (step -1) thing the player
// could act on, wrapping at the ends, so one key walks the turn's decisions
// without aiming. What it walks depends on what is in front of the player: a
// prompt's buttons, a card prompt's candidates, or — in ordinary play — the cards
// that can actually do something followed by End turn, since a card with no move
// is not a stop on the way to the next decision and ending the turn is the last
// one there is.
func (g *game) tabSel(ctx app.Context, step int) {
	if n := g.promptButtons(); n > 0 {
		g.btnCursor, g.hasBtnCursor = cycleIdx(n, g.btnCursor, g.hasBtnCursor, step), true
		return
	}
	if _, ok := g.tabCandidates(); ok {
		g.tabCandidate(step)
		return
	}
	rows := g.navRows()
	var all []engine.LocalID
	for _, ids := range rows {
		for _, id := range ids {
			if g.tabbable(id) {
				all = append(all, id)
			}
		}
	}
	i := cycleIdx(len(all)+1, g.tabPos(all), g.hasBtnCursor || g.hasSel, step)
	if i == len(all) {
		g.clearSelection()
		g.btnCursor, g.hasBtnCursor = 0, true
		return
	}
	g.btnCursor, g.hasBtnCursor = 0, false
	next := all[i]
	for r, ids := range rows {
		if c := indexOfID(ids, next); c >= 0 {
			g.navSelect(ctx, rows, r, c)
			return
		}
	}
}

// tabCandidates returns the candidate list a card prompt currently has up, and
// whether one is up at all — a chooser's candidates, or (sharing the same Tab
// cursor) the enemy creatures a declared fight may still hit. Unifying the two
// here means a fix to how the cursor lands or draws applies to both at once.
func (g *game) tabCandidates() ([]engine.LocalID, bool) {
	switch {
	case g.choosing:
		return g.chooserCandidates, true
	case g.phase == phaseFightTarget:
		return g.g.FightTargets(g.active(), g.attacker), true
	}
	return nil, false
}

// tabCandidate moves the card prompt's cursor one step along its candidates. An
// optional prompt ("exhaust up to 3 creatures") is a cycle of its candidates plus
// one tail stop on Done, because stopping early is one of its answers and a
// player answering from the keyboard has to be able to reach it. A fight target
// has no such tail stop — once a fight is declared some enemy must be hit — so
// it just cycles.
func (g *game) tabCandidate(step int) {
	cands, _ := g.tabCandidates()
	if !g.choosing || !g.chooserDeclinable {
		g.promptCursor, g.hasCursor = cycleID(cands, g.promptCursor, step), true
		return
	}
	i := cycleIdx(len(cands)+1, g.promptPos(cands), g.hasBtnCursor || g.hasCursor, step)
	if i == len(cands) {
		g.hasCursor = false
		g.btnCursor, g.hasBtnCursor = 0, true
		return
	}
	g.promptCursor, g.hasCursor = cands[i], true
	g.btnCursor, g.hasBtnCursor = 0, false
}

// promptPos locates the Tab cursor in an optional card prompt's cycle of
// candidates plus the Done button that tails it, or -1 when it is on neither.
func (g *game) promptPos(cands []engine.LocalID) int {
	if g.hasBtnCursor {
		return len(cands)
	}
	if !g.hasCursor {
		return -1
	}
	return indexOfID(cands, g.promptCursor)
}

// tabPos locates the Tab cursor in the ordinary-play cycle of cards plus the End
// turn button that tails it, or -1 when it is on neither.
func (g *game) tabPos(cards []engine.LocalID) int {
	if g.hasBtnCursor {
		return len(cards)
	}
	if !g.hasSel {
		return -1
	}
	return indexOfID(cards, g.sel)
}

// tabbable reports whether Tab stops on a card: one the active player can act
// with right now — a hand card they can play or discard, or a creature or
// artifact they can use. The opponent's cards are read-only, so Tab passes them.
func (g *game) tabbable(id engine.LocalID) bool {
	if containsID(g.g.Hand(g.active()), id) {
		return g.usableFromHand(id)
	}
	kind := g.boardKindOf(id)
	if kind == selOther {
		return false
	}
	return g.actionable(id, kind)
}

// promptButtons is how many buttons the prompt in front of the player offers,
// in the order they are drawn. Zero means the prompt is not a set of buttons (or
// there is none) and Tab goes back to walking cards.
func (g *game) promptButtons() int {
	switch {
	case g.choosingOption:
		return len(g.optionLabels)
	case g.forgingKey >= 0:
		return len(g.remainingKeyColors(g.forgingKey))
	case g.phase == phaseHouse:
		return len(g.pickableHouses())
	case g.phase == phaseFlank:
		return 2 // left flank, right flank
	}
	return 0
}

// pressButton answers the current prompt with the button Tab stopped on, so a
// yes/no, a house, a key colour, a flank, or ending the turn can be chosen
// without the mouse.
func (g *game) pressButton(ctx app.Context) bool {
	i, n := g.btnCursor, g.promptButtons()
	if !g.hasBtnCursor {
		return false
	}
	if n == 0 {
		// No button prompt is up, so the only buttons Tab parks on are an optional
		// card prompt's Done and, in ordinary play, End turn.
		switch {
		case g.isDoneCursor():
			g.declineChooser(ctx, app.Event{})
		case g.isEndTurnCursor():
			g.endTurn(ctx, app.Event{})
		default:
			return false
		}
		g.btnCursor, g.hasBtnCursor = 0, false
		return true
	}
	if i < 0 || i >= n {
		return false
	}
	switch {
	case g.choosingOption:
		g.chooseOptionIdx(i)(ctx, app.Event{})
	case g.forgingKey >= 0:
		g.pickForgeColor(g.remainingKeyColors(g.forgingKey)[i])(ctx, app.Event{})
	case g.phase == phaseHouse:
		g.pickHouse(g.pickableHouses()[i])(ctx, app.Event{})
	case g.phase == phaseFlank:
		g.playFlank(i == 0)(ctx, app.Event{})
	}
	g.btnCursor, g.hasBtnCursor = 0, false
	return true
}

// isButtonCursor reports whether the nth button of the current prompt draws as
// the one Tab has stopped on.
func (g *game) isButtonCursor(i int) bool {
	return g.hasBtnCursor && g.btnCursor == i
}

// isEndTurnCursor reports whether Tab has walked past the last actionable card
// and parked on End turn. A card prompt is up in front of it, so it does not
// count while one is being answered.
func (g *game) isEndTurnCursor() bool {
	return g.hasBtnCursor && !g.choosing && g.promptButtons() == 0
}

// isDoneCursor reports whether Tab has walked past the last candidate of an
// optional card prompt and parked on its Done button.
func (g *game) isDoneCursor() bool {
	return g.hasBtnCursor && g.choosing && g.chooserDeclinable
}

// confirmPrompt answers whatever prompt is up with the candidate Tab stopped on
// — a button, or a card — so a prompt can be finished without the mouse. It
// reports whether it did anything, so a key bound to both "confirm the Tab
// cursor" and "affirm the default answer" (Space) knows to fall back.
func (g *game) confirmPrompt(ctx app.Context) bool {
	if g.pressButton(ctx) {
		return true
	}
	cands, ok := g.tabCandidates()
	if !ok || !g.hasCursor || !containsID(cands, g.promptCursor) {
		return false
	}
	if g.phase == phaseFightTarget {
		g.fightTargetID(ctx, g.promptCursor)
		return true
	}
	g.chooseCandidate(ctx, g.promptCursor)
	return true
}

// cycleIdx returns the button index one step on from cur in a row of n, wrapping
// around. With no cursor yet it starts at whichever end the step comes from.
func cycleIdx(n, cur int, has bool, step int) int {
	if !has {
		if step > 0 {
			return 0
		}
		return n - 1
	}
	return ((cur+step)%n + n) % n
}

// cycleID returns the id one step after cur in ids, wrapping around. An id that
// is not in the list starts the cycle at whichever end the step comes from.
func cycleID(ids []engine.LocalID, cur engine.LocalID, step int) engine.LocalID {
	if len(ids) == 0 {
		return 0
	}
	i := indexOfID(ids, cur)
	if i < 0 {
		if step > 0 {
			return ids[0]
		}
		return ids[len(ids)-1]
	}
	return ids[((i+step)%len(ids)+len(ids))%len(ids)]
}

// isSelected reports whether a card draws as the selected one: the current
// selection, or the candidate a prompt's Tab cursor has landed on.
func (g *game) isSelected(id engine.LocalID) bool {
	if _, ok := g.tabCandidates(); ok {
		return g.hasCursor && g.promptCursor == id
	}
	return g.hasSel && g.sel == id
}
