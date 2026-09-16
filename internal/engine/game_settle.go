package engine

// A creature's destroyed state is checked continuously, not only when it is
// dealt damage. Power is dynamic — upgrades, power counters, and constant
// abilities all feed Power — so a creature can become destroyable with no effect
// touching it: the card granting it +power leaves play, and it drops to 0 power
// or down to the damage already marked on it. shouldDestroy names the state;
// settleDestroyed is what notices it has become true.

// settleDestroyed destroys every creature the board has left in a destroyable
// state, repeating until nothing more dies — each destruction can take away the
// buff that was keeping the next creature alive. It is a no-op while a
// destruction batch or an earlier sweep is already running, so a batch keeps its
// simultaneous "Destroyed:" timing instead of being split up.
func (g *Game) settleDestroyed(controller int) {
	if g.settling {
		return
	}
	g.settling = true
	g.replacedThisSweep = nil
	defer func() {
		g.settling = false
		g.replacedThisSweep = nil
	}()
	for {
		var dying []LocalID
		for p := range 2 {
			for _, id := range g.Battleline(p) {
				if g.shouldDestroy(id) {
					dying = append(dying, id)
				}
			}
			for _, id := range g.Artifacts(p) {
				if g.artifactShouldSelfDestroy(id) {
					dying = append(dying, id)
				}
			}
		}
		if len(dying) == 0 {
			return
		}
		g.destroyBatch(controller, dying)
	}
}

// destructionReplacedThisSweep reports whether a creature's destruction was already
// replaced earlier in this sweep. A replacement that heals damage but does not lift
// the reason the creature is destroyable — it sits at 0 power — would otherwise be
// re-detected and re-replaced every pass, hanging the sweep, so its second
// destruction in the same sweep resolves for real instead of being replaced again.
func (g *Game) destructionReplacedThisSweep(id LocalID) bool {
	for _, r := range g.replacedThisSweep {
		if r == id {
			return true
		}
	}
	return false
}
