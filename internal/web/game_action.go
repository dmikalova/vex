package web

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/match"
)

// This file is the action plumbing: every engine mutation runs through here, so
// each one is snapshotted for undo, resolved off the UI goroutine, and followed
// by the flash bookkeeping that animates what changed.

// runAction resolves an engine mutation on a background goroutine so the UI stays
// responsive while an effect blocks on the chooser. When it finishes it advances
// the phase and clears any transient selection on the UI goroutine.
func (g *game) runAction(ctx app.Context, fn func() error) {
	if g.busy {
		return
	}
	g.beginAction()
	g.busy = true
	g.status = ""
	// Complete via g.dispatch, which is bound in OnMount to the always-mounted
	// root component. Do NOT use ctx.Dispatch here: ctx is tied to the clicked
	// element (often a button this action removes from the DOM), and go-app drops
	// any Dispatch whose source element is no longer mounted — which would leave
	// the UI stuck on "resolving…" after a chooser.
	ctx.Async(func() {
		crashed, err := runSafely(func() error {
			if e := fn(); e != nil {
				return e
			}
			g.handOffEndedTurn()
			return nil
		})
		g.dispatch(func(ctx app.Context) { g.finishAction(ctx, crashed, err) })
	})
}

// finishAction settles a resolved root action on the UI goroutine: it clears the
// busy flag and remakes the chooser cancel channel, then either rolls the action
// back (a crash or a manual-mode Cancel), surfaces an illegal-move rejection, or
// commits it with afterAction. Either way the match is saved.
func (g *game) finishAction(ctx app.Context, crashed bool, err error) {
	g.busy = false
	// A manual-mode Cancel drained this action; roll it back to the snapshot
	// beginAction recorded and remake the cancel channel for the next action.
	cancelled := g.cancelling
	g.cancelling = false
	g.chooser.cancel = make(chan struct{})
	g.chooser.cancelled = false
	if crashed {
		// A corrupt engine state can panic mid-action (e.g. an
		// out-of-range card id). Peel the partial action back off the
		// command log and replay so the board stays consistent, then
		// surface the failure instead of freezing on "resolving…".
		g.rollbackLastRoot()
		g.setStatus(err.Error())
		g.save(ctx)
		return
	}
	if cancelled {
		g.rollbackLastRoot()
		g.save(ctx)
		return
	}
	if err != nil {
		g.setStatus(err.Error())
	}
	g.afterAction()
	g.save(ctx)
}

// runSafely runs a root action, converting a panic from a corrupt engine state
// into an error (crashed == true) instead of letting it kill the WASM goroutine
// and freeze the UI. A returned error with crashed == false is an ordinary
// illegal-move rejection.
func runSafely(fn func() error) (crashed bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("the game hit an unexpected error and rolled back: %v", r)
			crashed = true
		}
	}()
	return false, fn()
}

// setStatus shows a transient message in the controls area and arms a 5s
// auto-clear. statusGen guards the timer so a newer message is not wiped by an
// older one's timer.
func (g *game) setStatus(msg string) {
	g.status = msg
	if msg == "" {
		return
	}
	g.statusGen++
	gen := g.statusGen
	time.AfterFunc(5*time.Second, func() {
		g.dispatch(func(app.Context) {
			if g.statusGen == gen {
				g.status = ""
			}
		})
	})
}

// setNotice raises a standing message that stays up until clearNotice takes it
// down. Use it for a fault the player must act on — the save no longer being
// written, a set the deck generator does not know — where setStatus's 5s fade
// would let the fault pass unread.
func (g *game) setNotice(msg string) { g.notice = msg }

// clearNotice takes the standing notice down, for when the fault behind it is
// resolved.
func (g *game) clearNotice() { g.notice = "" }

// markTakeoff records where a card sits in hand as it is played, so flyIntoPlay
// can start the board card from there. A card with no hand slot on screen (or a
// render with no page behind it, as on the server) simply does not fly.
func (g *game) markTakeoff(id engine.LocalID) {
	g.takingOff = false
	el := app.Window().GetElementByID(handCardID(id))
	if !el.Truthy() {
		return
	}
	r := el.Call("getBoundingClientRect")
	g.takeoff = cardRect{
		x: r.Get("left").Float(),
		y: r.Get("top").Float(),
		w: r.Get("width").Float(),
	}
	g.takeoffID, g.takingOff = id, true
}

// flyIntoPlayDurMS and flyIntoPlayEasing tune the hand-to-board slide of a
// just-played card; flyIntoPlayZIndex lifts it above the settled board while it
// travels. Keep the duration in step with the CSS card transitions.
const (
	flyIntoPlayDurMS  = 260
	flyIntoPlayEasing = "cubic-bezier(0.2, 0.8, 0.3, 1)"
	flyIntoPlayZIndex = "40"
)

// flyIntoPlay slides a just-played card from the hand slot it left to the board
// slot it landed in, so the card the player let go of is the card that arrives
// rather than one fading in beside it. It runs on the render that first shows
// the card in play, when both ends of the move are on screen, and the arming is
// spent either way — a Tactic never lands on the board and so never flies.
func (g *game) flyIntoPlay() {
	if !g.takingOff {
		return
	}
	g.takingOff = false
	el := app.Window().GetElementByID(boardCardID(g.takeoffID))
	if !el.Truthy() {
		return
	}
	to := el.Call("getBoundingClientRect")
	from := fmt.Sprintf("translate(%.1fpx, %.1fpx)",
		g.takeoff.x-to.Get("left").Float(),
		g.takeoff.y-to.Get("top").Float())
	el.Call("animate",
		[]any{
			map[string]any{
				"transform": from,
				"zIndex":    flyIntoPlayZIndex,
			},
			map[string]any{
				"transform": "none",
				"zIndex":    flyIntoPlayZIndex,
			},
		},
		map[string]any{
			"duration": flyIntoPlayDurMS,
			"easing":   flyIntoPlayEasing,
		})
}

// beginAction starts a new log group for the root action about to run and marks
// the animation baseline. Every root action (an engine mutation via runAction or
// a manual edit) calls it, so log bubbles align with player intent and a new
// action drops any pending redo. Undo/redo themselves work off the command log
// (g.inputs), not a separate snapshot stack, so nothing is captured here for them
// beyond prevState, which computeFlashes diffs the resolved state against.
func (g *game) beginAction() {
	g.confirmEndTurn = false
	g.btnCursor, g.hasBtnCursor = 0, false
	g.hasUseTarget = false
	g.handSlot = g.selHandSlot()
	g.clearFlashes()
	g.prevState = g.g.State.FastCopy()
	g.prevValid = true
	g.redoLog = nil
	g.logGroups = append(g.logGroups, logMark{
		Start:  len(g.g.Log),
		Player: g.g.State.ActivePlayer,
	})
}

// clearFlashes drops every queued one-shot animation, so a state change the
// player did not act into (an undo, a new deal) does not replay the last one.
func (g *game) clearFlashes() {
	g.flashes = nil
	g.poolFlash = [2]bool{}
	g.keyFlash = [2]bool{}
	g.discardFlash = [2]bool{}
	g.flights = nil
}

// rebuildFromLog re-deals the match from the seed and replays the current command
// log to regenerate the exact state, typed log, and log groups (ADR 0039). It
// reports whether the replay reproduced the log cleanly; a divergence leaves the
// old game untouched. Manual mode is not part of the log, so it is carried across.
func (g *game) rebuildFromLog() bool {
	if g.defByName == nil {
		g.defByName = cardsByName()
	}
	wasManual := g.g != nil && g.g.Manual()
	if g.chooser != nil {
		g.chooser.drain()
	}
	eg, houses, mavericks, legacies, rosters, err := match.NewWithSets(
		"Player 1",
		"Player 2",
		g.seed,
		g.setNames,
	)
	if err != nil {
		return false
	}
	rc := &replayChooser{inputs: g.inputs}
	eg.SetChooser(0, rc)
	eg.SetChooser(1, rc)
	groups, err := driveReplay(eg, rc, g.defByName)
	if err != nil {
		return false
	}
	g.install(eg, houses, mavericks, legacies, rosters)
	g.logGroups = groups
	if wasManual {
		g.g.SetManual(true)
	}
	return true
}

// afterRebuild resets the transient UI a rebuild throws away and sets the
// animation baseline to the rebuilt board, so the next action diffs against it
// rather than the pre-rebuild game.
func (g *game) afterRebuild() {
	g.confirmEndTurn = false
	g.clearFlashes()
	g.inPlayPrev = g.inPlaySet()
	g.clearSelection()
	g.forgingKey = -1
	g.prevValid = false
	g.settlePhase()
}

// recomputeRootMarks rebuilds rootMarks from inputs, so a truncation or splice of
// the command log leaves the root-boundary index consistent with it.
func (g *game) recomputeRootMarks() {
	g.rootMarks = g.rootMarks[:0]
	for i, in := range g.inputs {
		if in.Kind.isRoot() {
			g.rootMarks = append(g.rootMarks, i)
		}
	}
}

// rollbackLastRoot drops the last recorded root action and everything after it
// from the command log, then rebuilds. A crashed or cancelled action has already
// recorded its root (and maybe some choices) before failing, so undoing it is
// peeling that partial segment back off the log.
func (g *game) rollbackLastRoot() {
	if len(g.rootMarks) == 0 {
		return
	}
	g.inputs = g.inputs[:g.rootMarks[len(g.rootMarks)-1]]
	g.recomputeRootMarks()
	g.rebuildFromLog()
	g.afterRebuild()
}

// atPrompt reports whether a prompt is waiting for an answer, with the action
// that raised it still in flight.
func (g *game) atPrompt() bool {
	return g.choosing || g.choosingOption || g.choosingPosition
}

// canUndo reports whether there is a step back to take. During a prompt there
// is: the prompt is part of the action that raised it, so undo rewinds to before
// that action rather than to a half-resolved state the player never saw. Outside
// a prompt, an action still resolving has no settled state to rewind to.
func (g *game) canUndo() bool {
	if len(g.rootMarks) == 0 {
		return false
	}
	if g.atPrompt() {
		return !g.cancelling
	}
	return !g.busy
}

func (g *game) canRedo() bool {
	return !g.busy && !g.choosing && !g.choosingOption && len(g.redoLog) > 0
}

// undoAction steps back to the state before the last root action by peeling its
// input segment off the command log and replaying what remains. During a prompt
// it instead drains the prompt so the in-flight action backs itself out, which
// rolls the log back to the same place — the prompt's own answers were never
// recorded, so there is nothing to redo forward into.
func (g *game) undoAction(ctx app.Context, _ app.Event) {
	if !g.canUndo() {
		return
	}
	if g.atPrompt() {
		g.cancelChooser(ctx, app.Event{})
		return
	}
	start := g.rootMarks[len(g.rootMarks)-1]
	seg := append([]input(nil), g.inputs[start:]...)
	g.redoLog = append(g.redoLog, seg)
	g.inputs = g.inputs[:start]
	g.recomputeRootMarks()
	g.rebuildFromLog()
	g.afterRebuild()
	g.save(ctx)
}

// redoAction re-applies the last undone action by splicing its segment back onto
// the command log and replaying.
func (g *game) redoAction(ctx app.Context, _ app.Event) {
	if !g.canRedo() {
		return
	}
	seg := g.redoLog[len(g.redoLog)-1]
	g.redoLog = g.redoLog[:len(g.redoLog)-1]
	g.inputs = append(g.inputs, seg...)
	g.recomputeRootMarks()
	g.rebuildFromLog()
	g.afterRebuild()
	g.save(ctx)
}

// afterAction settles the phase after an engine mutation: the game may be won,
// the active house may have been cleared (a new turn began), or play continues.
// The selection is handed on last, once the resting phase is known.
func (g *game) afterAction() {
	g.computeFlashes()
	g.settlePhase()
	g.advanceSelection()
}

// computeFlashes diffs the pre-action snapshot against the resolved state to queue
// one-shot animations: a card that took damage, gained on-card Æmber, changed
// power counters, was stunned or exhausted, or entered play pulses; a player who
// gained pool Æmber or forged a key pulses their score. The parity maps flip on
// each flash so the animation replays on repeats (see the flashes field).
// boardFlashes collects the per-card and per-player pulses for the current board,
// returning the set of cards in play now.
func (g *game) boardFlashes(
	prev *engine.GameState,
	flashes map[engine.LocalID]cardFlash,
) map[engine.LocalID]bool {
	inPlayNow := map[engine.LocalID]bool{}
	for p := range 2 {
		for _, id := range g.g.Battleline(p) {
			inPlayNow[id] = true
			g.cardFlags(id, prev, flashes)
		}
		for _, id := range g.g.Artifacts(p) {
			inPlayNow[id] = true
			g.cardFlags(id, prev, flashes)
		}
		g.playerFlashes(p, prev)
	}
	return inPlayNow
}

// playerFlashes pulses a player's score and discard pill. Cards leaving play cannot
// pulse — they are gone from the board — so the destination pulses instead, which
// is also the feedback for a discard.
func (g *game) playerFlashes(p int, prev *engine.GameState) {
	if g.g.State.Aember[p] > prev.Aember[p] {
		g.poolParity[p] = !g.poolParity[p]
		g.poolFlash[p] = true
	}
	if g.g.State.KeyCount(p) > prev.KeyCount(p) {
		g.keyParity[p] = !g.keyParity[p]
		g.keyFlash[p] = true
	}
	if g.g.State.Discard[p].Count > prev.Discard[p].Count {
		g.discardParity[p] = !g.discardParity[p]
		g.discardFlash[p] = true
	}
}

// useFlashes pulses the cards the running action used and disarms each one-shot.
// A card that left play as a side effect of its own ability cannot pulse in place,
// so — like the discard, pool, and key pulses — it is skipped rather than flashed
// somewhere it no longer is.
func (g *game) useFlashes(
	inPlayNow map[engine.LocalID]bool,
	flashes map[engine.LocalID]cardFlash,
) {
	mark := func(id engine.LocalID, pick func(*cardFlash) *bool) {
		if !inPlayNow[id] {
			return
		}
		f := flashes[id]
		*pick(&f) = true
		flashes[id] = f
	}
	// The two combatants clash, whether or not either took damage.
	if g.fighting {
		for _, id := range g.fighters {
			mark(id, func(f *cardFlash) *bool { return &f.fight })
		}
		g.fighting = false
	}
	if g.reaping {
		mark(g.reapID, func(f *cardFlash) *bool { return &f.reap })
		g.reaping = false
	}
	if g.acting {
		mark(g.actID, func(f *cardFlash) *bool { return &f.act })
		g.acting = false
	}
}

// entryFlashes pulses each card that is newly in play. A card played from hand
// flies in from its hand slot instead (flyIntoPlay), so it does not also pulse in
// place.
func (g *game) entryFlashes(
	inPlayNow map[engine.LocalID]bool,
	flashes map[engine.LocalID]cardFlash,
) {
	for id := range inPlayNow {
		if !g.inPlayPrev[id] && (!g.takingOff || id != g.takeoffID) {
			f := flashes[id]
			f.enter = true
			flashes[id] = f
		}
	}
}

func (g *game) computeFlashes() {
	if !g.prevValid {
		return
	}
	prev := &g.prevState
	if g.flashParity == nil {
		g.flashParity = map[engine.LocalID]bool{}
	}
	flashes := map[engine.LocalID]cardFlash{}
	inPlayNow := g.boardFlashes(prev, flashes)
	g.useFlashes(inPlayNow, flashes)
	g.entryFlashes(inPlayNow, flashes)
	// One parity flip per flashing card drives all its pulses at once.
	for id, f := range flashes {
		g.flashParity[id] = !g.flashParity[id]
		f.odd = g.flashParity[id]
		flashes[id] = f
	}
	g.flashes = flashes
	g.computeFlights(inPlayNow)
	g.inPlayPrev = inPlayNow
}

// computeFlights queues a flying card for each one that was on the board before
// this action and is not now, aimed at the zone pill it landed in. A card that
// has left play cannot animate where it was, so it animates on its way out.
func (g *game) computeFlights(inPlayNow map[engine.LocalID]bool) {
	g.flights = nil
	for id := range g.inPlayPrev {
		if inPlayNow[id] {
			continue
		}
		if player, zone, ok := g.landing(id); ok {
			g.flights = append(g.flights, flight{
				id:     id,
				player: player,
				zone:   zone,
			})
		}
	}
	if len(g.flights) > 0 {
		g.flightParity = !g.flightParity
	}
}

// landing finds the out-of-play zone a card is in now, as the player whose pill
// owns it and that pill's zone icon name. Both players are searched: a card
// leaves play into its owner's zone, which need not be its controller's.
func (g *game) landing(id engine.LocalID) (int, string, bool) {
	for p := range 2 {
		zones := []struct {
			name string
			ids  []engine.LocalID
		}{
			{"zone-discard", g.g.Discard(p)},
			{"zone-purge", g.g.Purge(p)},
			{"zone-archives", g.g.Archives(p)},
			{"zone-hand", g.g.Hand(p)},
			{"zone-deck", g.g.Deck(p)},
		}
		for _, z := range zones {
			if containsID(z.ids, id) {
				return p, z.name, true
			}
		}
	}
	return 0, "", false
}

// cardFlags records which state of one card changed this action (except entering
// play, handled by the caller). It leaves odd unset; computeFlashes finalizes it.
func (g *game) cardFlags(
	id engine.LocalID,
	prev *engine.GameState,
	out map[engine.LocalID]cardFlash,
) {
	now, was := g.g.State.Cards[id], prev.Cards[id]
	f := out[id]
	f.damage = f.damage || now.Damage > was.Damage
	f.amber = f.amber || now.Amber > was.Amber
	f.power = f.power || now.PowerCounters != was.PowerCounters
	f.exhaust = f.exhaust || (now.Exhausted && !was.Exhausted)
	f.stun = f.stun || (now.Stunned && !was.Stunned)
	if f.damage || f.amber || f.power || f.exhaust || f.stun {
		out[id] = f
	}
}

// inPlaySet returns the ids currently in play, so restore/newMatch can seed
// inPlayPrev and avoid flagging the restored board as freshly entered.
func (g *game) inPlaySet() map[engine.LocalID]bool {
	set := map[engine.LocalID]bool{}
	for p := range 2 {
		for _, id := range g.g.Battleline(p) {
			set[id] = true
		}
		for _, id := range g.g.Artifacts(p) {
			set[id] = true
		}
	}
	return set
}

// settlePhase picks the resting interaction phase from the engine state: the game
// is over, a new turn needs a house, or play continues. Transient phases (picking
// a flank or a fight target) depend on a live selection, so they are never a
// resting phase and a resumed match always lands on one of these.
//
// The turn only needs a house while the engine still waits at PhaseChooseHouse; a
// player locked out of every house chooses No House, which leaves ActiveHouse at
// HouseNone but advances the engine past the choice — so key off the phase, not
// ActiveHouse, or a No-House turn would loop back to the picker with no way to end.
func (g *game) settlePhase() {
	switch {
	case g.g.Winner() >= 0:
		g.phase = phaseOver
	case g.g.State.ActiveHouse == engine.HouseNone &&
		g.g.Phase() == engine.PhaseChooseHouse:
		g.phase = phaseHouse
	default:
		g.phase = phaseMain
	}
}
