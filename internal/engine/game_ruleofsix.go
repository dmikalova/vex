package engine

// The Rule of Six caps how many times the active player may use a single card
// *name* in one turn at six. Every usage of a name shares one pool: playing a
// copy, discarding a copy, using a copy (reap, fight, or Action:), each Destroyed:
// resolution, and each resolution of a self-repeating or Replicator-style chain
// past its free first. A use buys the first resolution of its own ability for
// free; each loop or chained trigger after that charges one usage — a chained
// Replicator trigger against the card that started the chain, so two Replicators
// reaching for each other spend one pool between them. The pool is tracked per
// card in GameState.UsagesThisTurn and summed here across every copy of the name
// regardless of owner, so eight Replicators or three copies of an Automaton all
// draw from the same six — and controlling an opponent's copy of a name you also
// own shares your six rather than opening a second pool. The ledger is safe to key
// by name alone because it resets every turn (StartTurn) and only the active
// player uses cards within a turn, so a usage recorded this turn is always theirs.
// See ADR 0043 for the whole design — counting points, the free first resolution,
// cascade-root attribution, and the two enforcement boundaries.

// nameUsagesThisTurn is how many times this turn any card sharing id's name has
// been used — the running total the Rule of Six caps. Copies are summed by name
// regardless of owner, since only the active player uses cards this turn. An id
// outside the catalog (an effect resolved with no bound source) counts as unused.
func (g *Game) nameUsagesThisTurn(id LocalID) int {
	if int(id) >= len(g.cat.defs) {
		return 0
	}
	name := g.cat.def(id).Name
	total := 0
	for i := range g.cat.defs {
		if g.cat.def(LocalID(i)).Name == name {
			total += int(g.State.UsagesThisTurn[i])
		}
	}
	return total
}

// atRuleOfSix reports whether id's card name has reached its six usages this turn,
// so no further usage of that name may resolve.
func (g *Game) atRuleOfSix(id LocalID) bool {
	return g.nameUsagesThisTurn(id) >= RuleOfSix
}

// recordUsage counts one usage of id toward its card name's Rule-of-Six pool.
func (g *Game) recordUsage(id LocalID) {
	g.State.UsagesThisTurn[id]++
}
