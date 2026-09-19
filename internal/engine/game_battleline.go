package engine

// SwapCards exchanges two cards' places. When both sit in the same battleline,
// only the ordered slots move and the creatures keep all their card state. When
// one card is in play and the other rests in a discard pile, the two exchange
// places across zones: the resting card enters play in the in-play card's slot,
// and the in-play card leaves to that discard pile (swapAcrossZones). A pair the
// swap cannot place — two creatures on opposite battlelines, or two resting cards
// — is left where it is.
func (g *Game) SwapCards(a, b LocalID) {
	for player := range g.State.Battleline {
		line := &g.State.Battleline[player]
		ai := line.indexOf(a)
		bi := line.indexOf(b)
		if ai < 0 || bi < 0 {
			continue
		}
		line.IDs[ai], line.IDs[bi] = line.IDs[bi], line.IDs[ai]
		g.record(PositionsSwapped{A: a, B: b})
		return
	}
	g.swapAcrossZones(a, b)
}

// swapAcrossZones exchanges an in-play creature with one resting in a discard
// pile: the resting card enters play in the creature's battleline slot, exhausted
// with full armor and firing its enters-play abilities (not its Play: ability),
// while the creature leaves play to the discard pile the resting card came from.
// The creature sheds all its card state on the way out exactly as any creature
// leaving play does — its Æmber goes to its opponent, its upgrades and counters
// are discarded (leavePlayTeardown). The move is a plain relocation, not a
// destruction of its own: a creature leaving play this way counts as destroyed
// only when an open Destroyed window already enrolled it (Gebuk). It does nothing
// unless exactly one card is in a battleline and the other rests in a discard.
func (g *Game) swapAcrossZones(a, b LocalID) {
	inPlay, resting := a, b
	if !g.inPlay(inPlay) {
		inPlay, resting = b, a
	}
	if !g.inPlay(inPlay) || g.inPlay(resting) {
		return
	}
	restingOwner := g.owner(resting)
	if !g.State.Discard[restingOwner].contains(resting) {
		return
	}
	controller := g.controller(inPlay)
	line := &g.State.Battleline[controller]
	idx := line.indexOf(inPlay)
	if idx < 0 {
		return // the in-play card is an artifact; cross-zone artifact swap is unsupported
	}
	g.record(CardsSwapped{A: inPlay, B: resting, FromPlayer: restingOwner, FromZone: Discard})
	// Filing, not an attempt: the swap already settled, so it owes no ward check.
	g.fileFromPlay(inPlay, func(half LocalID, o int) { g.State.Discard[o].add(half) })
	g.State.Discard[restingOwner].remove(resting)
	core := &g.State.Cards[resting]
	core.Exhausted = true
	core.ArmorRemaining = int16(g.armor(resting))
	// leavePlayTeardown may cascade and shrink the line below the slot idx
	// captured before removal; clamp so the reinsert lands on the flank.
	line.insertAt(min(idx, int(line.Count)), resting)
	g.emitEnters(resting)
}

// MoveToFlank moves a creature to a flank of its own controller's battleline: the
// right flank when right is true, otherwise the left. Only the ordered slot moves;
// the creature keeps all its damage, upgrades, status, and control.
func (g *Game) MoveToFlank(id LocalID, right bool) {
	for player := range g.State.Battleline {
		line := &g.State.Battleline[player]
		if !line.remove(id) {
			continue
		}
		if right {
			line.add(id)
		} else {
			line.insertAt(0, id)
		}
		g.record(MovedToFlank{Creature: id, Right: right})
		return
	}
}

// MoveWithinBattleline repositions a creature anywhere in its own controller's
// battleline: it is removed from the line and reinserted before the creature now
// at index pos (pos == len leaves it on the right flank). chooser is the player
// making the placement decision, which may be the creature's opponent (Malison
// moves an enemy creature). Only the ordered slot moves; the creature keeps all
// its damage, upgrades, status, and control.
func (g *Game) MoveWithinBattleline(chooser int, id LocalID) {
	for player := range g.State.Battleline {
		line := &g.State.Battleline[player]
		idx := line.indexOf(id)
		if idx < 0 {
			continue
		}
		// Choose the slot against the line without the moved creature, but do
		// not remove it yet. Removing it first exposes an intermediate board
		// whose settle boundary (ADR 0029, raised when the placement prompt is
		// presented) can destroy a neighbor the moved creature was holding up by
		// position, shrinking the line and leaving the chosen slot out of range.
		// Snapshot the other creatures, choose, then move atomically.
		full := line.slice()
		others := make([]LocalID, 0, len(full)-1)
		others = append(others, full[:idx]...)
		others = append(others, full[idx+1:]...)
		pos := g.choosePosition(chooser, id, "Choose where to move "+g.Name(id), others)
		line.remove(id)
		line.insertAt(pos, id)
		g.record(MovedWithinBattleline{Creature: id})
		return
	}
}

// placeGainedOnFlank moves a creature that just entered controller's battleline
// through a control gain or a gift — a Treachery creature entering the opponent's
// line, a seized creature reverting to its owner — onto the flank the active
// player chooses. It is the seam these raw control movers route through so none
// silently assumes a flank: the active player always makes the placement call
// (even onto the opponent's line), an artifact has no flank, and a one-creature
// line offers no choice. The creature is already listed under controller (the raw
// mover appended it); this only repositions it to the chosen flank.
func (g *Game) placeGainedOnFlank(id LocalID, controller int) {
	if g.TypeOf(id) != Creature || len(g.State.Battleline[controller].slice()) <= 1 {
		return
	}
	right := g.ChooseOption(g.State.ActivePlayer, id,
		FlankPromptPrefix+g.Name(id), []string{FlankLeftLabel, FlankRightLabel}) == 1
	g.MoveToFlank(id, right)
}
