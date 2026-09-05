package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file is the new-game set selection shown in the action bar over the
// current board: a same-sets shortcut when a previous game exists, then a
// player-by-player picker. beginSetup opens it, pickSet and continueSameSets
// answer it, cancelSetup backs out, and finishSetup deals the chosen match.

// beginSetup opens the set picker for player 0. It snapshots the running match's
// sets first, so the same-sets shortcut names them even as the player's choices
// overwrite setNames.
func (g *game) beginSetup() {
	g.prevSetNames = g.setNames
	g.awaitingSetup = true
	g.setPick = 0
}

// cancelSetup backs out of the picker, leaving the current game running.
func (g *game) cancelSetup(_ app.Context, _ app.Event) {
	g.awaitingSetup = false
}

// hasPrevSets reports whether a previous match's sets are known, i.e. the picker
// was opened over a running game rather than before the first deal.
func (g *game) hasPrevSets() bool {
	return g.prevSetNames[0] != "" && g.prevSetNames[1] != ""
}

// prevSetLabel names the previous match's sets for the same-sets shortcut, as one
// name when both players shared it or "A vs B" when they differed.
func (g *game) prevSetLabel() string {
	if g.prevSetNames[0] == g.prevSetNames[1] {
		return g.prevSetNames[0]
	}
	return g.prevSetNames[0] + " vs " + g.prevSetNames[1]
}

// pickSet records the current player's chosen set and advances: player 0's choice
// moves the picker to player 1, and player 1's choice deals the match.
func (g *game) pickSet(ctx app.Context, name string) {
	g.setNames[g.setPick] = name
	if g.setPick == 0 {
		g.setPick = 1
		return
	}
	g.finishSetup(ctx)
}

// continueSameSets deals a new match from the previous one's sets, skipping the
// per-player picking.
func (g *game) continueSameSets(ctx app.Context) {
	g.setNames = g.prevSetNames
	g.finishSetup(ctx)
}

// finishSetup leaves the picker and deals the match the chosen sets describe,
// seeding the animation baseline and saving it the way OnMount does.
func (g *game) finishSetup(ctx app.Context) {
	g.awaitingSetup = false
	g.newMatch()
	g.inPlayPrev = g.inPlaySet()
	g.save(ctx)
	g.scrollLogToBottom()
}
