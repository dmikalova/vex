package engine

// beforeForgePrevented fires the before-forge trigger window for the non-forging
// player's in-play cards — their "when your opponent would forge a key" abilities
// (Keyforgery) — then reports whether one of them cancelled forger's imminent key
// forge. The window runs before any Æmber leaves the pool, so a cancelled forge
// costs the forger nothing. The veto is a transient flag CancelForge sets and this
// reads and clears, mirroring FightCancelled.
func (g *Game) beforeForgePrevented(forger int) bool {
	guard := 1 - forger
	w := g.window()
	for _, id := range g.allInPlay(guard) {
		w.add(id, TriggerBeforeOpponentForgesKey, 0, false)
	}
	if len(w.pending) == 0 {
		return false
	}
	g.State.ForgePrevented = false
	g.resolveWindow(g.orderTriggered(forger, w.pending))
	prevented := g.State.ForgePrevented
	g.State.ForgePrevented = false
	return prevented
}
