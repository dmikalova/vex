package engine

import (
	"fmt"
	"slices"
)

// maxControlEntries bounds the global control stack. A fresh take supersedes its
// own source's earlier entry on the same card (see supersedeControl), so the stack
// holds at most one entry per distinct source controlling a card, not one per
// control change over the game. It therefore tracks the board, not the game
// length: a card taken permanently again and again (every permanent take shares
// the card as its source) stays at a single entry. 64 is far above any real board;
// a 65th entry panics rather than silently dropping a control effect.
const maxControlEntries = 64

// ControlEntry is one "take control" effect in the global stack
// (GameState.Controls), ordered by when it was applied. The current controller of
// a card is the Controller of its most recently pushed entry; earlier entries are
// the fallback a later one reverts to when it is removed (LIFO). An entry is
// dropped when its Source leaves play — an UntilCardLeavesPlay control names the
// seized card itself as Source, so it lasts as long as the card stays in play.
// Ownership is immutable and decided separately.
type ControlEntry struct {
	// Card is the seized card.
	Card LocalID
	// Controller is the player index (0 or 1) this effect hands the card to.
	Controller uint8
	// Source is the card whose lasting effect holds the control; when Source leaves
	// play the entry is dropped. An UntilCardLeavesPlay control names Card itself as its Source.
	Source LocalID
}

// takeControl records a new control effect over a card and moves it into the new
// controller's play area — a creature into their battleline, an artifact into
// their artifact row — without changing ownership. Ownership is immutable in
// KeyForge and still decides the out-of-play zone the card returns to. source is
// the card whose leaving play ends the control; an UntilCardLeavesPlay control names the
// seized card itself as source, so it lasts until the card leaves play. Control stacks
// LIFO: a later take takes precedence, and removing it falls back to the entry
// beneath (see releaseControlHeldBy).
func (g *Game) takeControl(id LocalID, controller int, source LocalID) {
	if !g.inPlay(id) {
		return
	}
	g.supersedeControl(id, source)
	g.pushControl(id, controller, source)
	g.placeUnderController(id, controller)
	g.record(ControlTaken{
		Player: controller,
		Card:   id,
	})
}

// supersedeControl drops a source's earlier control entry on a card before that
// source takes it again, so re-taking replaces rather than stacks a duplicate.
// This keeps each (card, source) pair to one entry, bounding the stack to the
// distinct in-play sources controlling a card — the board, not the game length.
// Every permanent take names the card itself as its source, so repeated permanent
// takes share a key and collapse to one entry the same way.
func (g *Game) supersedeControl(id, source LocalID) {
	g.dropControls(func(e ControlEntry) bool { return e.Card == id && e.Source == source })
}

// pushControl appends a control effect to the top of the stack. An entry past
// capacity panics — a caught invariant, never a silent drop.
func (g *Game) pushControl(id LocalID, controller int, source LocalID) {
	if int(g.State.ControlCount) >= maxControlEntries {
		panic(fmt.Sprintf("control table full: cannot take control of card %d", id))
	}
	g.State.Controls[g.State.ControlCount] = ControlEntry{
		Card:       id,
		Controller: uint8(controller),
		Source:     source,
	}
	g.State.ControlCount++
}

// controllerFromStack returns the controller named by the top of a card's control
// stack, or the card's owner when it carries no control entry.
func (g *Game) controllerFromStack(id LocalID) int {
	for i := int(g.State.ControlCount) - 1; i >= 0; i-- {
		if g.State.Controls[i].Card == id {
			return int(g.State.Controls[i].Controller)
		}
	}
	return g.owner(id)
}

// placeUnderController relists a card under controller — a creature in their
// battleline, an artifact in their artifact row — and refreshes the cached
// ControlPlus. It is a raw zone move that does not settle the board or touch the
// control stack; the caller pushes or pops the stack entry that justifies the move
// and settles once the whole batch has moved.
func (g *Game) placeUnderController(id LocalID, controller int) {
	g.removeFromPlayRows(id)
	if controller != g.owner(id) {
		g.State.Cards[id].ControlPlus = uint8(controller + 1)
	} else {
		g.State.Cards[id].ControlPlus = 0
	}
	if g.TypeOf(id) == Artifact {
		g.State.Artifacts[controller].add(id)
	} else {
		g.State.Battleline[controller].add(id)
	}
}

// releaseControlHeldBy drops every control effect the given source holds over
// other cards and re-derives each affected card's controller from the entry
// beneath — the LIFO fallback. It is the leave-play half of a control effect:
// control lasts exactly as long as its source stays in play. An UntilCardLeavesPlay control
// names the seized card itself as its source, so it is shed by clearControls when
// that card leaves play, not here. Every affected card is still in play — a card
// sheds its own control entries when it leaves (clearControls) — so re-deriving
// never touches a card that has left play.
func (g *Game) releaseControlHeldBy(source LocalID) {
	affected := g.dropControls(func(e ControlEntry) bool {
		return e.Source == source && e.Card != source
	})
	for _, id := range affected {
		controller := g.controllerFromStack(id)
		g.placeUnderController(id, controller)
		// Reverting control is a control change, so the active player chooses which
		// flank the creature returns to (it never silently keeps the old slot).
		g.placeGainedOnFlank(id, controller)
		g.record(ControlReturned{
			Card:  id,
			Owner: controller,
		})
	}
}

// clearControls drops every control effect naming a card as its seized subject,
// compacting the stack in place. A card sheds its control entries when it leaves
// play; removeFromPlay calls this for every exit. It does not relist the card —
// the card is on its way out of play.
func (g *Game) clearControls(id LocalID) {
	g.dropControls(func(e ControlEntry) bool { return e.Card == id })
}

// dropControls removes every stack entry the predicate matches, keeping the
// survivors contiguous and order-by-construction, and returns the distinct seized
// cards whose stack changed so callers can re-derive their controllers.
func (g *Game) dropControls(match func(ControlEntry) bool) []LocalID {
	var affected []LocalID
	w := 0
	for i := 0; i < int(g.State.ControlCount); i++ {
		e := g.State.Controls[i]
		if match(e) {
			if !slices.Contains(affected, e.Card) {
				affected = append(affected, e.Card)
			}
			continue
		}
		g.State.Controls[w] = e
		w++
	}
	for i := w; i < int(g.State.ControlCount); i++ {
		g.State.Controls[i] = ControlEntry{}
	}
	g.State.ControlCount = uint8(w)
	return affected
}
