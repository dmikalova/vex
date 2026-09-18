package engine

// A creature's destroyed state is checked continuously, not only when it is
// dealt damage. Power is dynamic — upgrades, power counters, and constant
// abilities all feed Power — so a creature can become destroyable with no effect
// touching it: the card granting it +power leaves play, and it drops to 0 power
// or down to the damage already marked on it. shouldDestroy names the state;
// settleDestroyed is what notices it has become true.

// simultaneously runs a batch of board changes as one moment and then settles the
// board once, instead of letting each change settle behind it. Holding the
// settling flag is what makes the batch simultaneous: without it, the first
// card's departure can destroy a card still waiting its turn in the batch, so the
// outcome would depend on the order the batch happens to visit its members.
//
// Every effect that moves several cards "at the same time" — Epic Quest archiving
// a creature and the neighbour it was buffing, Timequake shuffling a creature and
// its upgrades away, Transporter Platform returning a creature with its upgrades
// — routes through here rather than looping a single-card move.
// Pinned by TestSimultaneouslySettlesOnce.
func (g *Game) simultaneously(controller int, batch func()) {
	wasSettling, wasDeferring := g.settling, g.deferringLeaves
	g.settling, g.deferringLeaves = true, true
	batch()
	g.settling, g.deferringLeaves = wasSettling, wasDeferring
	if !wasDeferring {
		g.flushDeferredLeaves()
	}
	g.settleDestroyed(controller)
}

// flushDeferredLeaves resolves the "Leaves Play:" windows a batch gathered, once
// every card in the batch has moved. A nested batch does not flush — the
// outermost one owns the moment, so the whole set still resolves together.
//
// The batch's own controller settles the board, but the **active** player orders
// this window: a batch can take cards from both players out at once, and ADR 0013
// gives simultaneous triggers to the active player to order regardless of who
// resolved the effect (TestLeavesPlayWindowIsOrderedByTheActivePlayer).
func (g *Game) flushDeferredLeaves() {
	pending := g.deferredLeaves
	g.deferredLeaves = nil
	if len(pending) > 0 {
		g.resolveWindow(g.orderTriggered(g.State.ActivePlayer, pending))
	}
}

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
	defer func() {
		g.settling = false
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
