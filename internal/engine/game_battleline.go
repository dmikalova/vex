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
		// A power bonus that only reaches a flank may have moved off one of them.
		g.settleDestroyed(player)
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
		// A power bonus that only reaches a flank may have moved onto or off it.
		g.settleDestroyed(player)
		return
	}
}
