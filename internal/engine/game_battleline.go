package engine

// SwapBattlelinePositions exchanges two creatures' positions in the same
// battleline. Only the ordered battleline slots move; the creatures keep all
// damage, upgrades, status, control, and other card state.
func (g *Game) SwapBattlelinePositions(a, b LocalID) {
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
		if !line.remove(id) {
			continue
		}
		pos := g.choosePosition(chooser, id, "Choose where to move "+g.Name(id), line.slice())
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
