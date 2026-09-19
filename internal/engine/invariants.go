package engine

import "fmt"

// InvariantError reports the first violation of a flat-state invariant that must
// hold in any legal game between actions, or nil when the state is sound. It reads
// only GameState and the catalog — no mutation, no I/O — so it is cheap enough to
// run after every step of a simulation (see internal/sim) and inside the engine at
// turn boundaries in an -tags assert build (see assertInvariants). It is exported
// so the simulator can reuse the one true definition rather than restating it.
func (g *Game) InvariantError() error {
	if err := g.checkPlayerTotals(); err != nil {
		return err
	}
	var pl cardPlacement
	if err := g.tallyPlacement(&pl); err != nil {
		return err
	}
	if err := g.checkCardPresence(&pl); err != nil {
		return err
	}
	if err := g.checkRestingOwnership(); err != nil {
		return err
	}
	if err := g.checkDamageBelowPower(); err != nil {
		return err
	}
	return g.checkStaleInPlayState(&pl)
}

// checkPlayerTotals validates the per-player scalars — Æmber, forged-key prefix,
// chains — and the winner range.
func (g *Game) checkPlayerTotals() error {
	for p := range 2 {
		if a := g.State.Aember[p]; a < 0 {
			return fmt.Errorf("player %d has negative Æmber (%d)", p, a)
		}
		// The key count is the non-empty prefix of KeyColors, so a gap would hide
		// forged keys from every reader rather than merely miscount them.
		for i := g.State.KeyCount(p) + 1; i < MaxKeys; i++ {
			if c := g.State.KeyColors[p][i]; c != KeyColorNone {
				return fmt.Errorf(
					"player %d has a %s key forged after an unforged slot (index %d)",
					p,
					c,
					i,
				)
			}
		}
		if c := g.State.Chains[p]; c < 0 {
			return fmt.Errorf("player %d has negative chains (%d)", p, c)
		}
	}
	if w := g.State.Winner; w < -1 || w > 1 {
		return fmt.Errorf("winner is out of range (%d, want -1, 0 or 1)", w)
	}
	return nil
}

// cardPlacement records where each registered card was found during the
// conservation tally: how many places hold it, and whether it is attached as an
// upgrade, placed under a host, or the linked art half of a gigantic.
type cardPlacement struct {
	count            [maxCards]int
	attached         [maxCards]bool
	underAttached    [maxCards]bool
	giganticAttached [maxCards]bool
}

// tallyPlacement walks every zone and in-play chain to fill pl, returning the
// first structural corruption (a mistyped or mis-linked attachment) it finds.
// Card conservation: every registered card must sit in exactly one place — some
// zone list, attached as an upgrade to an in-play creature, or placed under an
// in-play card. checkCardPresence then confirms the tally is exactly one.
func (g *Game) tallyPlacement(pl *cardPlacement) error {
	tally := func(ids []LocalID) {
		for _, id := range ids {
			pl.count[id]++
		}
	}
	for p := range 2 {
		tally(g.State.Hand[p].slice())
		tally(g.State.Deck[p].slice())
		tally(g.State.Battleline[p].slice())
		tally(g.State.Discard[p].slice())
		tally(g.State.Artifacts[p].slice())
		tally(g.State.Archives[p].slice())
		tally(g.State.Purge[p].slice())
		for _, id := range g.allInPlay(p) {
			if err := g.tallyUpgradeChain(id, pl); err != nil {
				return err
			}
			if err := g.tallyUnderChain(id, pl); err != nil {
				return err
			}
			if err := g.tallyGigantic(id, pl); err != nil {
				return err
			}
		}
	}
	return nil
}

// tallyUpgradeChain records id's upgrade chain into pl, checking each link is a
// legal type carrying an agreeing host back-link.
func (g *Game) tallyUpgradeChain(id LocalID, pl *cardPlacement) error {
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		pl.count[up]++
		pl.attached[up] = true
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
	return nil
}

// tallyUnderChain records the cards placed under id into pl, checking each carries
// an agreeing under-host back-link.
func (g *Game) tallyUnderChain(id LocalID, pl *cardPlacement) error {
	for u, ok := g.firstUnder(id); ok; u, ok = g.nextUnder(u) {
		pl.count[u]++
		pl.underAttached[u] = true
		if g.State.Cards[u].UnderHostPlus != underPlus(id) {
			return fmt.Errorf(
				"card %d (%s) is placed under %s but its host back-link disagrees",
				u,
				g.cat.def(u).Name,
				g.Name(id),
			)
		}
	}
	return nil
}

// tallyGigantic records id's linked art half into pl when id is the base half.
// The slot-less art half of a gigantic sits in no zone; it is accounted through
// its base half, which holds the battleline slot (ADR 0042).
func (g *Game) tallyGigantic(id LocalID, pl *cardPlacement) error {
	art, ok := g.giganticPartner(id)
	if !ok || g.cat.def(id).GiganticRole != GiganticBase {
		return nil
	}
	pl.count[art]++
	pl.giganticAttached[art] = true
	if g.State.Cards[art].GiganticPartnerPlus != giganticPlus(id) {
		return fmt.Errorf(
			"gigantic art half %d (%s) is linked to %s but its partner back-link disagrees",
			art,
			g.cat.def(art).Name,
			g.Name(id),
		)
	}
	return nil
}

// checkCardPresence confirms the conservation tally: every card sits in exactly
// one place, carries no negative Æmber, and holds no dangling host back-link.
func (g *Game) checkCardPresence(pl *cardPlacement) error {
	for id := 0; id < len(g.cat.defs); id++ {
		if pl.count[id] != 1 {
			return fmt.Errorf("card %d (%s) is in %d places, want exactly 1",
				id, g.cat.def(LocalID(id)).Name, pl.count[id])
		}
		if a := g.State.Cards[id].Amber; a < 0 {
			return fmt.Errorf("card %d (%s) has negative Æmber on it (%d)",
				id, g.cat.def(LocalID(id)).Name, a)
		}
		// An upgrade in play must be attached: a card that thinks it has a host but
		// no creature holds it in a chain is a dangling attachment.
		if g.State.Cards[id].HostPlus != 0 && !pl.attached[id] {
			return fmt.Errorf(
				"card %d (%s) has a host back-link but no creature holds it (dangling upgrade)",
				id,
				g.cat.def(LocalID(id)).Name,
			)
		}
		// A card placed under a host must be reachable from that host's chain, the
		// same dangling-attachment shape as an upgrade above.
		if g.State.Cards[id].UnderHostPlus != 0 && !pl.underAttached[id] {
			return fmt.Errorf(
				"card %d (%s) has an under-host back-link but no host holds it (dangling under-card)",
				id,
				g.cat.def(LocalID(id)).Name,
			)
		}
		// A linked art half with no base holding it is a dangling gigantic — the
		// leave-play funnel must clear the link on both halves together.
		if g.State.Cards[id].GiganticPartnerPlus != 0 &&
			g.cat.def(LocalID(id)).GiganticRole == GiganticArt && !pl.giganticAttached[id] {
			return fmt.Errorf(
				"card %d (%s) is a linked gigantic art half but no base holds it (dangling gigantic)",
				id,
				g.cat.def(LocalID(id)).Name,
			)
		}
	}
	return nil
}

// checkRestingOwnership confirms a card only ever rests in its owner's zones.
// Ownership decides where a card goes out of play, so a card can only ever rest
// in its *owner's* hand, deck, discard pile, or purge pile — control is temporary
// and never follows a card out of play. Archives is the one exception: an effect
// can archive a card into the opponent's archives. In-play rows are excluded
// because that is exactly where control shows, and an under-chain is excluded
// because a card may be placed under an enemy card.
func (g *Game) checkRestingOwnership() error {
	for p := range 2 {
		for _, z := range []struct {
			name string
			ids  []LocalID
		}{
			{"hand", g.State.Hand[p].slice()},
			{"deck", g.State.Deck[p].slice()},
			{"discard pile", g.State.Discard[p].slice()},
			{"purge pile", g.State.Purge[p].slice()},
		} {
			for _, id := range z.ids {
				if o := g.owner(id); o != p {
					return fmt.Errorf(
						"card %d (%s) is owned by player %d but sits in player %d's %s",
						id, g.cat.def(id).Name, o, p, z.name,
					)
				}
			}
		}
	}
	return nil
}

// checkDamageBelowPower confirms no in-play creature sits with damage at or above
// its power. A creature whose damage has caught up with its power is destroyed, so
// one can never be sitting in play in that state — this catches a power change (a
// buff leaving) that no state-based sweep noticed.
func (g *Game) checkDamageBelowPower() error {
	for p := range 2 {
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
	return nil
}

// checkStaleInPlayState confirms no out-of-play card still carries in-play state.
// Per-match state belongs to cards in play. Leaving play zeroes the whole
// CardCore, so any card outside play (and not attached as an upgrade, and not
// placed under a host — which is deliberately out of play yet still carries its
// host and facedown links) carrying damage, Æmber, a stun, counters, or a host
// link is a leave-play path that forgot to reset it. A gigantic's slot-less art
// half is in play through its base, so it is skipped here too.
func (g *Game) checkStaleInPlayState(pl *cardPlacement) error {
	for id := 0; id < len(g.cat.defs); id++ {
		lid := LocalID(id)
		if g.inPlay(lid) || pl.attached[id] || pl.underAttached[id] || pl.giganticAttached[id] {
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
