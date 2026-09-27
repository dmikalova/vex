package web

import (
	"errors"
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
	"github.com/dmikalova/vex/internal/session"
)

// This file is the action plumbing: every player input reaches the engine as a
// Command applied to the session (ADR 0040), so each one is recorded in the
// command log undo rewinds over, and each is followed by the settling step that
// animates what changed and presents whatever the engine asks next.
//
// THE ROOT BOUNDARY IS "Pending() is a RequestAction". Between two RequestActions
// lies exactly one root action and every prompt it raised — which is the single
// signal beginAction and afterAction bracket the client's per-action bookkeeping
// with: the undo mark, the flash baseline, and the log group.

// applyRoot takes one root action. It opens the action, hands the Command to the
// session, and settles whatever comes back — a prompt the action raised, or the
// next RequestAction if it resolved outright. A command the session refuses never
// began an action, so its bookkeeping is dropped again and the engine's own reason
// for refusing is shown.
func (g *game) applyRoot(ctx app.Context, cmd engine.Command) {
	if g.s == nil || g.atPrompt() {
		return
	}
	g.status = ""
	g.redoLog = nil
	g.beginAction()
	err := g.s.Apply(cmd)
	switch {
	case errors.Is(err, session.ErrIllegal), errors.Is(err, session.ErrFinished):
		g.abandonAction()
		g.setStatus(g.refusalOf(cmd).Error())
	case err != nil:
		g.brokeAction(ctx, err)
	default:
		g.settleAfterApply(ctx)
	}
}

// answer applies one answer to the pending prompt and settles. Answering is not
// the start of an action — the action that raised the prompt is still the one
// running — so it opens no log group and marks no new undo boundary.
func (g *game) answer(ctx app.Context, cmd engine.Command) {
	if g.s == nil {
		return
	}
	g.redoLog = nil
	if err := g.s.Apply(cmd); err != nil {
		g.brokeAction(ctx, err)
		return
	}
	g.settleAfterApply(ctx)
}

// applyManual records and performs one manual force-edit. A force-edit is not an
// answer to the pending request — it deliberately bypasses the rules and the turn
// loop — so it goes to the session's own manual path and leaves whatever the
// engine is waiting on untouched. It is still a root action for the client's
// purposes: it gets its own log group and its own undo boundary.
func (g *game) applyManual(ctx app.Context, cmd engine.Command) {
	if g.s == nil {
		return
	}
	g.redoLog = nil
	g.beginAction()
	if err := g.s.ApplyManual(cmd); err != nil {
		g.abandonAction()
		g.setStatus(err.Error())
		return
	}
	g.save(ctx)
}

// brokeAction handles an action the engine could not finish: the command is in the
// log but the state it left is half-resolved, so the whole root action is rewound
// and the failure surfaced rather than the client sitting on a broken board.
func (g *game) brokeAction(ctx app.Context, err error) {
	g.rewindLastRoot()
	g.setStatus(fmt.Sprintf("the game hit an unexpected error and rolled back: %v", err))
	g.save(ctx)
}

// settleAfterApply brings the client up to date with whatever the session yielded.
// It first answers anything the player has already decided (autoAnswer), then
// mirrors the badge preview, presents a prompt it has not presented before, and —
// at a root boundary, where the pending request is a RequestAction or the match is
// over — settles the action that has just finished.
func (g *game) settleAfterApply(ctx app.Context) {
	req, live := g.settlePending()
	if !live || req.Kind == engine.RequestAction {
		g.afterAction()
	}
	g.save(ctx)
}

// settlePending brings the client's presentation up to date with whatever the
// session is now waiting on, and reports it: it answers anything the player has
// already decided, mirrors the badge preview, and presents a request it has not
// presented before. It is the half of settling that a rebuild shares with an
// apply, which is why it is not folded into settleAfterApply.
func (g *game) settlePending() (engine.Request, bool) {
	for g.autoAnswer() { //nolint:revive // the loop body is the call's side effect
	}
	req, live := g.pending()
	g.syncBadge(req, live)
	g.presentIfNew(req, live)
	return req, live
}

// autoAnswer answers a request the player has already decided, so the prompt is
// never shown: the play-as-which question a Play creature / Play upgrade button
// pre-answered, and a placement with nowhere else to put the creature. It reports
// whether it answered, so the caller loops until the session asks something real.
func (g *game) autoAnswer() bool {
	req, ok := g.prompt()
	if !ok {
		return false
	}
	switch req.Kind {
	case engine.RequestOption:
		if i, armed := g.armedUpgradeChoice(req.Options); armed {
			return g.s.Apply(engine.Command{
				Kind:  engine.CommandOption,
				Index: i,
			}) == nil
		}
	case engine.RequestPosition:
		// With no other friendly creatures in play there is only one placement, so
		// place without asking. The engine already skips the prompt for an empty line;
		// this guards the web side against ever showing a choice with no alternatives.
		if len(req.Cards) == 0 {
			return g.s.Apply(engine.Command{Kind: engine.CommandPosition}) == nil
		}
	}
	return false
}

// presentIfNew sets the client up for a request it has not presented yet: it tears
// the previous prompt's presentation down and decides how the new one is offered.
// promptAt identifies a request by the number of commands applied before it — each
// request consumes exactly one — so a second settle over the same prompt leaves a
// zone viewer the player closed closed.
func (g *game) presentIfNew(req engine.Request, live bool) {
	if g.s.Len() == g.promptAt {
		return
	}
	g.promptAt = g.s.Len()
	g.clearPrompts()
	g.hasAbilityPick = false
	if !live {
		return
	}
	switch req.Kind {
	case engine.RequestPickCard, engine.RequestPickCardOrDecline:
		// A prompt taking over the board owns the highlight: drop any card the player
		// had selected before it opened so a stale selection ring does not linger on a
		// non-candidate card, reading as still-active while the board dims around the
		// candidates.
		g.clearSelection()
		g.presentPrompt(req.Cards, req.Kind == engine.RequestPickCardOrDecline)
	case engine.RequestReaction:
		g.clearSelection()
		if cards := reactionCards(req.Reactions); len(cards) >= 2 {
			g.presentPrompt(cards, false)
		}
	case engine.RequestOption:
		// A prompt whose options are the whole card database (Etan's Jar) is answered
		// through the card-name typeahead, not a list of a thousand buttons.
		if g.cardNameOptions() {
			g.pickerOpen, g.pickerNaming = true, true
			g.pickerQuery, g.pickerFocused, g.pickerCursor = "", false, 0
		}
	}
}

// rejection explains, in the engine's own words, why a root action the legal set
// did not offer was refused. The session rejects such a command without performing
// it, so there is no engine error to surface; this asks the same Can* readers the
// action bar offers the verb through, which is the client asking the engine rather
// than deciding a rule of its own.
func (g *game) rejection(cmd engine.Command) error {
	eg, p := g.eng(), g.active()
	switch cmd.Kind {
	case engine.CommandPlayCreature, engine.CommandPlayArtifact,
		engine.CommandPlayTactic, engine.CommandPlayUpgrade:
		id, ok := handCardAt(eg, p, cmd.Hand)
		if !ok {
			return errNotNow
		}
		return playTypeError(eg.CanPlay(p, id), eg.Def(id).Type)
	case engine.CommandDiscardFromHand:
		id, ok := handCardAt(eg, p, cmd.Hand)
		if !ok {
			return errNotNow
		}
		return eg.CanDiscard(p, id)
	case engine.CommandReap:
		return eg.CanUseTo(p, cmd.Card, engine.ReapUse)
	case engine.CommandFight:
		return eg.CanUseTo(p, cmd.Card, engine.FightUse)
	case engine.CommandUseAction:
		if eg.TypeOf(cmd.Card) == engine.Artifact {
			return eg.CanUseArtifact(p, cmd.Card)
		}
		return eg.CanUseTo(p, cmd.Card, engine.ActionUse)
	case engine.CommandUnstun:
		return eg.CanUse(p, cmd.Card)
	}
	return errNotNow
}

// refusalOf is rejection with a fallback: a Can* reader can answer nil for a
// command the legal set still did not offer (a play the phase does not allow at
// all), and a refusal always has to say something.
func (g *game) refusalOf(cmd engine.Command) error {
	if err := g.rejection(cmd); err != nil {
		return err
	}
	return errNotNow
}

// errNotNow is the fallback reason for a refused action with no Can* reader of its
// own to ask — a house choice or an end of turn taken in a phase that does not
// offer it, which the handlers already guard against.
var errNotNow = errors.New("that cannot be done right now")

// handCardAt is the card a hand index names, reporting false when the index is
// out of range — a command built against a hand that has since changed.
func handCardAt(eg *engine.Game, player, i int) (engine.LocalID, bool) {
	hand := eg.Hand(player)
	if i < 0 || i >= len(hand) {
		return 0, false
	}
	return hand[i], true
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

// beginAction opens a root action, at the root boundary the action starts from:
// it marks where the action's command segment begins in the session log (the undo
// cursor) and where its log lines begin (the bubble), and snapshots the state
// computeFlashes will diff the result against. rootMarks and logGroups are
// appended together and only here, so the i-th entry of each describes the same
// action and one undo peels one entry off each.
func (g *game) beginAction() {
	g.confirmEndTurn = false
	g.btnCursor, g.hasBtnCursor = 0, false
	g.hasUseTarget = false
	g.handSlot = g.selHandSlot()
	g.clearFlashes()
	g.prevState = g.eng().State.FastCopy()
	g.prevValid = true
	g.rootMarks = append(g.rootMarks, g.s.Len())
	g.logGroups = append(g.logGroups, logMark{
		Start:  len(g.eng().Log),
		Player: g.eng().State.ActivePlayer,
	})
}

// abandonAction takes back what beginAction opened, for a command the session
// refused: nothing was applied, so the action never began.
func (g *game) abandonAction() {
	g.rootMarks = g.rootMarks[:len(g.rootMarks)-1]
	g.logGroups = g.logGroups[:len(g.logGroups)-1]
	g.prevValid = false
	// The verb animations a handler armed before applying (which card reaped, which
	// two fought, which used its action) describe something that did not happen, so
	// they are disarmed rather than left to fire on whatever action comes next.
	g.fighting, g.reaping, g.acting = false, false, false
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

// canUndo reports whether there is a step back to take. During a prompt there is:
// the prompt is part of the action that raised it, so undo rewinds to before that
// action rather than to a half-resolved state the player never saw.
func (g *game) canUndo() bool { return len(g.rootMarks) > 0 }

// canRedo reports whether an undone action can be put back. Not mid-prompt: the
// prompt belongs to a NEW action, which has already dropped the redo history.
func (g *game) canRedo() bool { return !g.atPrompt() && len(g.redoLog) > 0 }

// undoAction steps back to the state before the last root action, by replaying the
// match from a fresh deal up to that action's start (session.Undo). CANCELLING A
// PROMPT IS THE SAME OPERATION and the same call: a prompt belongs to the action
// that raised it, so backing out of one is rewinding past that action — and unlike
// draining the prompt, rewinding never resolves the rest of the effect first, so
// nothing it would have done flickers on the way out.
//
// Undo never reaches command 0. That command is the first-player roll, and undoing
// it would re-roll and re-deal — which is a new match, not a step back — so the
// floor is the first root action's own mark, which is always at least 1.
func (g *game) undoAction(ctx app.Context, _ app.Event) {
	if !g.canUndo() {
		return
	}
	mark := g.rootMarks[len(g.rootMarks)-1]
	seg := g.commandsFrom(mark)
	if !g.rewindTo(mark) {
		return
	}
	g.redoLog = append(g.redoLog, seg)
	g.save(ctx)
}

// rewindLastRoot backs the last root action out with no redo entry, for an action
// that broke rather than one the player took back.
func (g *game) rewindLastRoot() {
	if !g.canUndo() {
		return
	}
	g.rewindTo(g.rootMarks[len(g.rootMarks)-1])
}

// rewindTo rewinds the session to command n and drops the client's own bookkeeping
// for the root action that began there. It reports whether the rewind took; a
// session that cannot replay its own prefix leaves the board as it was.
func (g *game) rewindTo(n int) bool {
	if err := g.s.Undo(n); err != nil {
		return false
	}
	g.rootMarks = g.rootMarks[:len(g.rootMarks)-1]
	g.logGroups = g.logGroups[:len(g.logGroups)-1]
	g.promptAt = g.s.Len()
	g.clearPrompts()
	g.hasAbilityPick = false
	g.resetBadgePreview()
	g.afterRebuild()
	return true
}

// commandsFrom copies the tail of the session's command log from index n, so undo
// can hand it to redo before session.Undo truncates it away.
func (g *game) commandsFrom(n int) []engine.Command {
	cmds := g.s.Record().Commands
	if n < 0 || n > len(cmds) {
		return nil
	}
	return cmds[n:]
}

// redoAction puts the last undone action back by applying its commands again. The
// session has no Redo of its own — Undo truncates the log — so what was truncated
// is the client's to remember and feed back through the ordinary apply path.
func (g *game) redoAction(ctx app.Context, _ app.Event) {
	if !g.canRedo() {
		return
	}
	seg := g.redoLog[len(g.redoLog)-1]
	g.redoLog = g.redoLog[:len(g.redoLog)-1]
	g.beginAction()
	for _, cmd := range seg {
		if err := g.applyRecorded(cmd); err != nil {
			break
		}
	}
	g.promptAt = g.s.Len()
	g.afterRebuild()
	g.save(ctx)
}

// applyRecorded feeds one recorded command back in by KIND — a manual force-edit
// straight to the live game, everything else as an answer the session advances on
// — which is the same routing a session replay uses, so a rebuild here and a
// rebuild there cannot disagree about what a recorded command means.
func (g *game) applyRecorded(cmd engine.Command) error {
	if cmd.Kind.IsManual() {
		return g.s.ApplyManual(cmd)
	}
	return g.s.Apply(cmd)
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
		for _, id := range g.eng().Battleline(p) {
			inPlayNow[id] = true
			g.cardFlags(id, prev, flashes)
		}
		for _, id := range g.eng().Artifacts(p) {
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
	if g.eng().State.Aember[p] > prev.Aember[p] {
		g.poolParity[p] = !g.poolParity[p]
		g.poolFlash[p] = true
	}
	if g.eng().State.KeyCount(p) > prev.KeyCount(p) {
		g.keyParity[p] = !g.keyParity[p]
		g.keyFlash[p] = true
	}
	if g.eng().State.Discard[p].Count > prev.Discard[p].Count {
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
			{"zone-discard", g.eng().Discard(p)},
			{"zone-purge", g.eng().Purge(p)},
			{"zone-archives", g.eng().Archives(p)},
			{"zone-hand", g.eng().Hand(p)},
			{"zone-deck", g.eng().Deck(p)},
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
	now, was := g.eng().State.Cards[id], prev.Cards[id]
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
		for _, id := range g.eng().Battleline(p) {
			set[id] = true
		}
		for _, id := range g.eng().Artifacts(p) {
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
	case g.eng().Winner() >= 0:
		g.phase = phaseOver
	case g.eng().State.ActiveHouse == engine.HouseNone &&
		g.eng().Phase() == engine.PhaseChooseHouse:
		g.phase = phaseHouse
	default:
		g.phase = phaseMain
	}
}
