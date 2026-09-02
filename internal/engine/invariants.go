package engine

import "fmt"

// InvariantError reports the first violation of a flat-state invariant that must
// hold in any legal game between actions, or nil when the state is sound. It reads
// only GameState and the catalog — no mutation, no I/O — so it is cheap enough to
// run after every step of a simulation (see internal/sim) and inside the engine at
// turn boundaries in an -tags assert build (see assertInvariants). It is exported
// so the simulator can reuse the one true definition rather than restating it.
func (g *Game) InvariantError() error {
	for p := 0; p < 2; p++ {
		if a := g.State.Aember[p]; a < 0 {
			return fmt.Errorf("player %d has negative Æmber (%d)", p, a)
		}
		if k := g.State.Keys[p]; k < 0 || k > KeysToWin {
			return fmt.Errorf(
				"player %d has out-of-range key count (%d, want 0..%d)",
				p,
				k,
				KeysToWin,
			)
		}
		if c := g.State.Chains[p]; c < 0 {
			return fmt.Errorf("player %d has negative chains (%d)", p, c)
		}
	}
	if w := g.State.Winner; w < -1 || w > 1 {
		return fmt.Errorf("winner is out of range (%d, want -1, 0 or 1)", w)
	}

	// Card conservation: every registered card must sit in exactly one place —
	// some zone list, or attached as an upgrade to an in-play creature. A card that
	// vanishes or duplicates is a leak in a move-between-zones path.
	var count [maxCards]int
	var attached [maxCards]bool
	tally := func(ids []LocalID) {
		for _, id := range ids {
			count[id]++
		}
	}
	for p := 0; p < 2; p++ {
		tally(g.State.Hand[p].slice())
		tally(g.State.Deck[p].slice())
		tally(g.State.Battleline[p].slice())
		tally(g.State.Discard[p].slice())
		tally(g.State.Artifacts[p].slice())
		tally(g.State.Archives[p].slice())
		tally(g.State.Purge[p].slice())
		for _, id := range g.allInPlay(p) {
			for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
				count[up]++
				attached[up] = true
				if t := g.cat.def(up).Type; t != Upgrade {
					return fmt.Errorf(
						"card %d (%s) is in %s's upgrade chain but is a %s, not an Upgrade",
						up,
						g.cat.def(up).Name,
						g.Name(id),
						t,
					)
				}
				if g.State.Cards[up].HostPlus != upgradePlus(id) {
					return fmt.Errorf(
						"upgrade %d (%s) is in %s's chain but its host back-link disagrees",
						up,
						g.cat.def(up).Name,
						g.Name(id),
					)
				}
			}
		}
	}
	for id := 0; id < len(g.cat.defs); id++ {
		if count[id] != 1 {
			return fmt.Errorf("card %d (%s) is in %d places, want exactly 1",
				id, g.cat.def(LocalID(id)).Name, count[id])
		}
		if a := g.State.Cards[id].Amber; a < 0 {
			return fmt.Errorf("card %d (%s) has negative Æmber on it (%d)",
				id, g.cat.def(LocalID(id)).Name, a)
		}
		// An upgrade in play must be attached: a card that thinks it has a host but
		// no creature holds it in a chain is a dangling attachment.
		if g.State.Cards[id].HostPlus != 0 && !attached[id] {
			return fmt.Errorf(
				"card %d (%s) has a host back-link but no creature holds it (dangling upgrade)",
				id,
				g.cat.def(LocalID(id)).Name,
			)
		}
	}
	return nil
}
