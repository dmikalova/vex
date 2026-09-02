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
		g.logf("%s swaps positions with %s", g.Name(a), g.Name(b))
		// A power bonus that only reaches a flank may have moved off one of them.
		g.settleDestroyed(player)
		return
	}
}
