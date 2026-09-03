package web

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
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
		crashed, err := runSafely(fn)
		g.dispatch(func(ctx app.Context) {
			g.busy = false
			if crashed {
				// A corrupt engine state can panic mid-action (e.g. an
				// out-of-range card id). Roll back to the snapshot beginAction
				// recorded so the board stays consistent, then surface the
				// failure instead of freezing the UI on "resolving…".
				if n := len(g.undo); n > 0 {
					last := g.undo[n-1]
					g.undo = g.undo[:n-1]
					g.restore(last)
				}
				g.setStatus(err.Error())
				g.save(ctx)
				return
			}
			if err != nil {
				g.setStatus(err.Error())
			}
			g.afterAction()
			g.save(ctx)
		})
	})
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
			map[string]any{"transform": from, "zIndex": "40"},
			map[string]any{"transform": "none", "zIndex": "40"},
		},
		map[string]any{"duration": 260, "easing": "cubic-bezier(0.2, 0.8, 0.3, 1)"})
}

const maxUndo = 100

// snapshot captures the state, log, and log marks for undo/redo.
func (g *game) snapshot() undoEntry {
	return undoEntry{
		state:  g.g.State.FastCopy(),
		log:    append([]engine.Record(nil), g.g.Log...),
		groups: append([]logMark(nil), g.logGroups...),
	}
}

// beginAction records an undo point and starts a new log group for the root
// action about to run. Every root action (an engine mutation via runAction or a
// manual edit) calls it, so undo steps and log bubbles align with player intent.
func (g *game) beginAction() {
	g.confirmEndTurn = false
	g.btnCursor, g.hasBtnCursor = 0, false
	g.handSlot = g.selHandSlot()
	g.clearFlashes()
	g.undo = append(g.undo, g.snapshot())
	if len(g.undo) > maxUndo {
		g.undo = g.undo[len(g.undo)-maxUndo:]
	}
	g.redo = nil
	g.logGroups = append(g.logGroups, logMark{Start: len(g.g.Log), Player: g.g.State.ActivePlayer})
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

// restore installs a snapshot and resets transient UI.
func (g *game) restore(e undoEntry) {
	g.confirmEndTurn = false
	g.clearFlashes()
	g.g.State = e.state
	g.g.Restore(e.log)
	g.logGroups = e.groups
	g.inPlayPrev = g.inPlaySet()
	g.clearSelection()
	g.forgingKey = -1
	g.settlePhase()
}

func (g *game) canUndo() bool {
	return !g.busy && !g.choosing && !g.choosingOption && len(g.undo) > 0
}

func (g *game) canRedo() bool {
	return !g.busy && !g.choosing && !g.choosingOption && len(g.redo) > 0
}

// undoAction steps back to the state before the last root action.
func (g *game) undoAction(ctx app.Context, _ app.Event) {
	if !g.canUndo() {
		return
	}
	g.redo = append(g.redo, g.snapshot())
	e := g.undo[len(g.undo)-1]
	g.undo = g.undo[:len(g.undo)-1]
	g.restore(e)
	g.save(ctx)
}

// redoAction re-applies the last undone action.
func (g *game) redoAction(ctx app.Context, _ app.Event) {
	if !g.canRedo() {
		return
	}
	g.undo = append(g.undo, g.snapshot())
	e := g.redo[len(g.redo)-1]
	g.redo = g.redo[:len(g.redo)-1]
	g.restore(e)
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
func (g *game) computeFlashes() {
	if len(g.undo) == 0 {
		return
	}
	prev := &g.undo[len(g.undo)-1].state
	if g.flashParity == nil {
		g.flashParity = map[engine.LocalID]bool{}
	}
	flashes := map[engine.LocalID]cardFlash{}
	inPlayNow := map[engine.LocalID]bool{}
	for p := 0; p < 2; p++ {
		for _, id := range g.g.Battleline(p) {
			inPlayNow[id] = true
			g.cardFlags(id, prev, flashes)
		}
		for _, id := range g.g.Artifacts(p) {
			inPlayNow[id] = true
			g.cardFlags(id, prev, flashes)
		}
		if g.g.State.Aember[p] > prev.Aember[p] {
			g.poolParity[p] = !g.poolParity[p]
			g.poolFlash[p] = true
		}
		if g.g.State.Keys[p] > prev.Keys[p] {
			g.keyParity[p] = !g.keyParity[p]
			g.keyFlash[p] = true
		}
		// Cards leaving play cannot pulse — they are gone from the board — so the
		// destination pulses instead, which is also the feedback for a discard.
		if g.g.State.Discard[p].Count > prev.Discard[p].Count {
			g.discardParity[p] = !g.discardParity[p]
			g.discardFlash[p] = true
		}
	}
	// The two combatants clash, whether or not either took damage.
	if g.fighting {
		for _, id := range g.fighters {
			if inPlayNow[id] {
				f := flashes[id]
				f.fight = true
				flashes[id] = f
			}
		}
		g.fighting = false
	}
	for id := range inPlayNow {
		// A card played from hand flies in from its hand slot instead (flyIntoPlay),
		// so it does not also pulse in place.
		if !g.inPlayPrev[id] && (!g.takingOff || id != g.takeoffID) {
			f := flashes[id]
			f.enter = true
			flashes[id] = f
		}
	}
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
			g.flights = append(g.flights, flight{id: id, player: player, zone: zone})
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
	for p := 0; p < 2; p++ {
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
func (g *game) settlePhase() {
	switch {
	case g.g.Winner() >= 0:
		g.phase = phaseOver
	case g.g.State.ActiveHouse == engine.HouseNone:
		g.phase = phaseHouse
	default:
		g.phase = phaseMain
	}
}
