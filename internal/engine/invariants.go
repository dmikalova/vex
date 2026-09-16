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
	// some zone list, attached as an upgrade to an in-play creature, or placed
	// under an in-play card. A card that vanishes or duplicates is a leak in a
	// move-between-zones path.
	var count [maxCards]int
	var attached [maxCards]bool
	var underAttached [maxCards]bool
	var giganticAttached [maxCards]bool
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
				// A chain holds upgrades and — since ADR 0026 — creatures played as
				// upgrades. Any other printed type threaded into a chain is corruption.
				if t := g.cat.def(up).Type; t != Upgrade && t != Creature {
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
			for u, ok := g.firstUnder(id); ok; u, ok = g.nextUnder(u) {
				count[u]++
				underAttached[u] = true
				if g.State.Cards[u].UnderHostPlus != underPlus(id) {
					return fmt.Errorf(
						"card %d (%s) is placed under %s but its host back-link disagrees",
						u,
						g.cat.def(u).Name,
						g.Name(id),
					)
				}
			}
			// The slot-less art half of a gigantic sits in no zone; it is accounted
			// through its base half, which holds the battleline slot (ADR 0042).
			if art, ok := g.giganticPartner(id); ok &&
				g.cat.def(id).GiganticRole == GiganticBase {
				count[art]++
				giganticAttached[art] = true
				if g.State.Cards[art].GiganticPartnerPlus != giganticPlus(id) {
					return fmt.Errorf(
						"gigantic art half %d (%s) is linked to %s but its partner back-link disagrees",
						art,
						g.cat.def(art).Name,
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
		// A card placed under a host must be reachable from that host's chain, the
		// same dangling-attachment shape as an upgrade above.
		if g.State.Cards[id].UnderHostPlus != 0 && !underAttached[id] {
			return fmt.Errorf(
				"card %d (%s) has an under-host back-link but no host holds it (dangling under-card)",
				id,
				g.cat.def(LocalID(id)).Name,
			)
		}
		// A linked art half with no base holding it is a dangling gigantic — the
		// leave-play funnel must clear the link on both halves together.
		if g.State.Cards[id].GiganticPartnerPlus != 0 &&
			g.cat.def(LocalID(id)).GiganticRole == GiganticArt && !giganticAttached[id] {
			return fmt.Errorf(
				"card %d (%s) is a linked gigantic art half but no base holds it (dangling gigantic)",
				id,
				g.cat.def(LocalID(id)).Name,
			)
		}
	}

	// A creature whose damage has caught up with its power is destroyed, so one can
	// never be sitting in play in that state — this catches a power change (a buff
	// leaving) that no state-based sweep noticed.
	for p := 0; p < 2; p++ {
		for _, id := range g.State.Battleline[p].slice() {
			if power := g.Power(id); power <= int(g.State.Cards[id].Damage) {
				return fmt.Errorf(
					"creature %d (%s) is in play with %d power and %d damage, want power above damage",
					id,
					g.cat.def(id).Name,
					power,
					g.State.Cards[id].Damage,
				)
			}
		}
	}

	// Per-match state belongs to cards in play. Leaving play zeroes the whole
	// CardCore, so any card outside play (and not attached as an upgrade, and not
	// placed under a host — which is deliberately out of play yet still carries
	// its host and facedown links) carrying damage, Æmber, a stun, counters, or a
	// host link is a leave-play path that forgot to reset it. A gigantic's slot-less
	// art half is in play through its base, so it is skipped here too.
	for id := 0; id < len(g.cat.defs); id++ {
		lid := LocalID(id)
		if g.inPlay(lid) || attached[id] || underAttached[id] || giganticAttached[id] {
			continue
		}
		if core := g.State.Cards[id]; core != (CardCore{}) {
			return fmt.Errorf(
				"card %d (%s) is not in play but still carries in-play state (%+v)",
				id, g.cat.def(lid).Name, core,
			)
		}
	}
	return nil
}
