package web

import (
	"errors"
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file holds the handlers for taking a turn: selecting a card, choosing a
// house, playing (by click or by drag), reaping, using, fighting, and ending the
// turn.

// ---- selection ----
//
// Card click handlers take the card id as a parameter rather than capturing it in
// a per-card closure. go-app compares event handlers by function pointer, so
// every closure built from the same func literal looks identical to it and a
// captured id would never be refreshed when the board re-renders (e.g. when the
// turn flips). cardView instead keeps the id in a field and passes the live value
// to these methods on click.

// selectBoardID selects a card in play, deriving whether it belongs to the active
// player (actionable) or is another card (read-only).
func (g *game) selectBoardID(_ app.Context, id engine.LocalID) {
	if g.busy || g.choosing || g.phase == phaseFightTarget {
		return
	}
	g.sel, g.selKind, g.selHand, g.hasSel = id, g.boardKindOf(id), -1, true
	g.status = ""
}

// selectHandID selects a card in the active player's hand, recovering its hand
// index from the id.
func (g *game) selectHandID(_ app.Context, id engine.LocalID) {
	if g.busy || g.choosing || g.phase == phaseFightTarget {
		return
	}
	g.selectHand(id)
}

// selectHand makes a card in hand the selection, recovering its hand index from
// the id. It is the selection itself, without the guards a click has to pass.
func (g *game) selectHand(id engine.LocalID) {
	idx := indexOfID(g.g.Hand(g.active()), id)
	if idx < 0 {
		return
	}
	g.sel, g.selKind, g.selHand, g.hasSel = id, selHand, idx, true
	g.status = ""
}

// selHandSlot is the place the selected card holds in the hand as drawn, or -1
// when the selection is not a card in hand. beginAction records it so
// advanceSelection can hand the selection on once the card is gone.
func (g *game) selHandSlot() int {
	if !g.hasSel || g.selKind != selHand {
		return -1
	}
	return indexOfID(g.sortedHand(g.active()), g.sel)
}

// advanceSelection is what the selection does once an action has resolved. A
// card played or discarded from hand passes the selection to whatever card now
// holds its place in hand (the last card, when it was the last), so playing a
// run of cards from the keyboard keeps its place instead of starting over from
// nothing. It only lands on a card there is still something to do with, so
// playing the last playable card leaves nothing selected rather than parking on
// a card the turn cannot touch. Anything else — a reap, a fight, a new turn —
// simply clears.
func (g *game) advanceSelection() {
	slot, gone := g.handSlot, !g.hasSel || !containsID(g.g.Hand(g.active()), g.sel)
	g.handSlot = -1
	g.clearSelection()
	if slot < 0 || !gone || g.phase != phaseMain {
		return
	}
	hand := g.sortedHand(g.active())
	if len(hand) == 0 {
		return
	}
	from := min(slot, len(hand)-1)
	for i := range hand {
		if id := hand[(from+i)%len(hand)]; g.usableFromHand(id) {
			g.selectHand(id)
			return
		}
	}
}

// clickAway drops the selection when a click lands on the board's background
// rather than on a card or a player bar. The click bubbles up from whatever it
// hit, so the target is asked what it belongs to.
func (g *game) clickAway(_ app.Context, e app.Event) {
	if !g.hasSel || g.busy || g.choosing ||
		g.phase == phaseFightTarget || g.phase == phaseFlank {
		return
	}
	if t := e.Get("target"); t.Truthy() {
		if t.Call("closest", ".card, .score-pill").Truthy() {
			return
		}
	}
	g.clearSelection()
}

// boardKindOf reports how a card in play should be treated when selected: one of
// the active player's own creatures/artifacts, or another (read-only) card.
func (g *game) boardKindOf(id engine.LocalID) selKind {
	active := g.active()
	switch {
	case containsID(g.g.Battleline(active), id):
		return selYourCreature
	case containsID(g.g.Artifacts(active), id):
		return selYourArtifact
	default:
		return selOther
	}
}

// ---- turn / play handlers ----

func (g *game) pickHouse(h engine.House) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if g.busy || g.choosing {
			return
		}
		p := g.active()
		g.runAction(ctx, func() error { return g.g.ChooseHouse(p, h) })
	}
}

func (g *game) endTurn(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.choosingOption {
		return
	}
	// Only end from the resting main phase; mid-action phases have their own flow.
	if g.phase != phaseMain {
		return
	}
	// If the player could still act, arm a confirm and wait for a second end-turn.
	// The End turn button turns into a red "Confirm end turn", which says so where
	// the player is already looking — no transient message needed.
	if !g.confirmEndTurn && g.hasMoves() {
		g.confirmEndTurn = true
		return
	}
	g.confirmEndTurn = false
	g.status = ""
	p := g.active()
	g.runAction(ctx, func() error {
		opp := 1 - p
		g.g.EndPlayPhase(p) // the end-of-turn phase narrates where both players stand
		g.g.StartTurn(opp)  // forge + start-of-turn triggers for the next player
		return nil
	})
}

// hasMoves reports whether the active player could still act this turn: a playable
// hand card, a usable creature, or a usable artifact. It drives the end-turn
// confirmation — with nothing left to do, ending needs no confirm.
func (g *game) hasMoves() bool {
	p := g.active()
	for _, id := range g.g.Hand(p) {
		if g.g.CanPlay(p, id) == nil {
			return true
		}
	}
	for _, id := range g.g.Battleline(p) {
		if g.actionable(id, selYourCreature) {
			return true
		}
	}
	for _, id := range g.g.Artifacts(p) {
		if g.actionable(id, selYourArtifact) {
			return true
		}
	}
	return false
}

// play resolves the selected hand card. Creatures ask for a flank first (unless
// the battleline is empty); everything else plays immediately.
func (g *game) play(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain || g.selKind != selHand {
		return
	}
	p := g.active()
	idx := g.selHand
	def := g.g.Def(g.sel)
	g.markTakeoff(g.sel)
	switch def.Type {
	case engine.Creature:
		if len(g.g.Battleline(p)) == 0 {
			g.runAction(
				ctx,
				func() error { _, err := g.g.PlayCreature(p, idx, false); return playTypeError(err, def.Type) },
			)
			return
		}
		g.phase = phaseFlank
	case engine.Artifact:
		g.runAction(
			ctx,
			func() error { _, err := g.g.PlayArtifact(p, idx); return playTypeError(err, def.Type) },
		)
	case engine.Tactic:
		g.runAction(ctx, func() error { return playTypeError(g.g.PlayAction(p, idx), def.Type) })
	case engine.Upgrade:
		g.runAction(
			ctx,
			func() error { _, err := g.g.PlayUpgrade(p, idx); return playTypeError(err, def.Type) },
		)
	}
}

// playTypeError makes the generic "cannot play this type" restriction explicit
// about which card type is barred (e.g. "Tactic cards cannot be played").
func playTypeError(err error, t engine.CardType) error {
	if errors.Is(err, engine.ErrCannotPlayType) {
		return fmt.Errorf("%s cards cannot be played", t)
	}
	return err
}

func (g *game) playFlank(left bool) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if g.busy || g.choosing || g.phase != phaseFlank || g.selKind != selHand {
			return
		}
		p, idx := g.active(), g.selHand
		g.markTakeoff(g.sel)
		g.runAction(
			ctx,
			func() error { _, err := g.g.PlayCreature(p, idx, left); return playTypeError(err, engine.Creature) },
		)
	}
}

func (g *game) discard(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain || g.selKind != selHand {
		return
	}
	p, idx := g.active(), g.selHand
	g.runAction(ctx, func() error { return g.g.DiscardFromHand(p, idx) })
}

// ---- drag and drop (hand → board) ----

// startHandDrag begins dragging a playable hand card. It selects the card so the
// drop shares the same target, marks a drag in progress so the board shows as a
// drop zone, and hides the hover preview (mouseleave does not fire during drag).
func (g *game) startHandDrag(ctx app.Context, id engine.LocalID) {
	if g.busy || g.choosing || g.choosingOption || g.phase != phaseMain {
		return
	}
	g.hasHover, g.hoverDef = false, nil
	g.selectHandID(ctx, id)
	g.dragging = true
}

// endHandDrag clears the drag state when the pointer is released, whether or not
// the card landed on the board.
func (g *game) endHandDrag(_ app.Context, _ engine.LocalID) {
	g.dragging = false
	g.hasHover, g.hoverDef = false, nil
}

// dropOnBoard plays the dragged hand card when it is released over the play area,
// following the normal flow — a creature still prompts for its flank. A card that
// cannot be played from hand is left where it is.
func (g *game) dropOnBoard(ctx app.Context, e app.Event) {
	e.PreventDefault()
	if !g.dragging {
		return
	}
	g.dragging = false
	if g.busy || g.choosing || g.choosingOption || g.phase != phaseMain || g.selKind != selHand {
		return
	}
	if !g.playableFromHand(g.sel) {
		return
	}
	g.play(ctx, e)
}

// ---- creature actions ----

func (g *game) reap(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain || g.selKind != selYourCreature {
		return
	}
	p, id := g.active(), g.sel
	g.runAction(ctx, func() error { return g.g.Reap(p, id) })
}

func (g *game) useAction(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain {
		return
	}
	if g.selKind != selYourCreature && g.selKind != selYourArtifact {
		return
	}
	p, id := g.active(), g.sel
	g.runAction(ctx, func() error { return g.g.UseAction(p, id) })
}

// startFight enters fight-target selection for the selected creature, after
// checking it can actually be used (so the player is not left picking a target
// for an exhausted or out-of-house attacker).
func (g *game) startFight(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain || g.selKind != selYourCreature {
		return
	}
	if err := g.g.CanUse(g.active(), g.sel); err != nil {
		g.setStatus(err.Error())
		return
	}
	g.attacker = g.sel
	g.status = ""
	// With a single legal target there is nothing to choose, so the fight resolves
	// straight away instead of asking for the only possible answer.
	if targets := g.g.FightTargets(g.active(), g.attacker); len(targets) == 1 {
		g.fightTargetID(ctx, targets[0])
		return
	}
	g.phase = phaseFightTarget
}

func (g *game) fightTargetID(ctx app.Context, defender engine.LocalID) {
	if g.busy || g.choosing {
		return
	}
	p, att := g.active(), g.attacker
	g.phase = phaseMain
	g.fighters = [2]engine.LocalID{att, defender}
	g.fighting = true
	g.runAction(ctx, func() error { return g.g.Fight(p, att, defender) })
}

func (g *game) cancelTargeting(_ app.Context, _ app.Event) {
	g.phase = phaseMain
	g.attacker = 0
}
