package engine

import "math"

// maxAember is the most Æmber one holder — a card or a pool — can carry. It is
// the range of the int16 the flat state stores Æmber in, which is paid for 128
// times over across the battleline. A doubling chain (Binate Rupture) is the only
// thing that can approach it, so a total past it is saturated rather than wrapped.
const maxAember = math.MaxInt16

// clampAember saturates an Æmber total at maxAember and records the Æmber the
// maximum turned away. It clamps only the top: wrapping there would turn a huge
// pile into negative Æmber that later leaks back out. Going below zero is a real
// bug, so it is left to InvariantError to catch rather than hidden here.
func (g *Game) clampAember(total int) int16 {
	if total > maxAember {
		g.record(AemberLostToMaximum{Amount: total - maxAember})
		total = maxAember
	}
	return int16(total)
}

// addAmberOn changes the Æmber sitting on a card, saturating at maxAember.
// A card that has left play takes no write — an ability that places Æmber on its
// own source after that source was destroyed (Strange Gizmo forging mid-window)
// lands on nothing rather than banking Æmber on a card in a discard pile.
func (g *Game) addAmberOn(id LocalID, delta int) {
	c := g.stateOf(id)
	if c == nil {
		return
	}
	c.Amber = g.clampAember(int(c.Amber) + delta)
}

// gainAember adds Æmber from the common supply to a player's pool. It is the
// single seam for pool gains: before the Æmber reaches the pool, a continuous
// replacement such as Ether Spider may capture the incoming Æmber instead, so the
// pool's owner keeps what they already have. Movements that already come from
// another pool or card (steal, capture, returning captured Æmber) are not gains and
// intentionally bypass this helper.
func (g *Game) gainAember(player, amount int) (LocalID, bool) {
	if capturer, ok := g.aemberCaptorFor(player); ok {
		g.addAmberOn(capturer, amount)
		return capturer, true
	}
	g.State.Aember[player] = g.clampAember(int(g.State.Aember[player]) + amount)
	return 0, false
}

// GainAember is the Resolver entry point for gainAember.
func (g *Game) GainAember(
	player, amount int,
) (LocalID, bool) {
	return g.gainAember(player, amount)
}

// aemberCaptorFor returns the in-play creature whose continuous replacement
// captures Æmber that would be added to player's pool, or ok=false when none does.
// A card's replacement names which pool it watches relative to the card's controller
// (Ether Spider watches its Opponent's). When several creatures could each capture
// the Æmber — two Ether Spiders in play — their controller chooses which one does.
func (g *Game) aemberCaptorFor(player int) (LocalID, bool) {
	var captors []LocalID
	for p := 0; p < 2; p++ {
		for _, id := range g.allInPlay(p) {
			def := g.cat.def(id)
			r := def.Replaces
			if g.TypeOf(id) != Creature ||
				r.Of != EventAemberAddedToPool ||
				r.With != Capture {
				continue
			}
			pool := p // the card's own pool (Controller)
			if r.Player == Opponent {
				pool = 1 - p
			}
			if pool == player {
				captors = append(captors, id)
			}
		}
	}
	switch len(captors) {
	case 0:
		return 0, false
	case 1:
		return captors[0], true
	}
	chosen, ok := g.ChooseCreature(
		g.controller(captors[0]),
		0,
		"Choose which creature captures the Æmber",
		captors,
	)
	if !ok {
		return captors[0], true
	}
	return chosen, true
}

// stolenRedirectSource returns the in-play card of either controller that carries
// the continuous replacement redirecting stolen Æmber into a capture
// (Gargantodon), so the log can name what caused the redirect. The redirect is
// global — it applies to every steal regardless of who controls the card — so no
// pool scoping is consulted.
func (g *Game) stolenRedirectSource() (LocalID, bool) {
	for p := 0; p < 2; p++ {
		for _, id := range g.allInPlay(p) {
			r := g.cat.def(id).Replaces
			if r.Of == EventAemberStolen && r.With == Capture {
				return id, true
			}
		}
	}
	return 0, false
}

// StolenAemberCaptor returns a creature player controls that captures Æmber a
// steal would otherwise add to player's pool, together with the card that
// redirected it, or ok=false when no in-play card redirects stolen Æmber or
// player controls no creature to hold it. When player controls several creatures,
// player chooses which one captures.
func (g *Game) StolenAemberCaptor(player int) (captor, cause LocalID, ok bool) {
	cause, ok = g.stolenRedirectSource()
	if !ok {
		return 0, 0, false
	}
	candidates := g.Battleline(player)
	switch len(candidates) {
	case 0:
		return 0, 0, false
	case 1:
		return candidates[0], cause, true
	}
	chosen, picked := g.ChooseCreature(
		player,
		0,
		"Choose which creature captures the stolen Æmber",
		candidates,
	)
	if !picked {
		return candidates[0], cause, true
	}
	return chosen, cause, true
}

// AemberTakenFromSupply reports whether Æmber a steal or capture takes from
// player's pool is drawn from the common supply instead, leaving the pool
// untouched (Po's Pixies), and returns the card that redirects it so the log can
// name the cause. It is the source half of the Æmber-flow replacement spine, the
// mirror of aemberCaptorFor on the destination half: it reads the continuous
// replacement each in-play card carries (Replaces), scoped to the pool it
// watches, rather than a bespoke flag.
func (g *Game) AemberTakenFromSupply(player int) (LocalID, bool) {
	for p := 0; p < 2; p++ {
		for _, id := range g.allInPlay(p) {
			r := g.cat.def(id).Replaces
			if r.Of != EventAemberTakenFromPool || r.With != FromCommonSupply {
				continue
			}
			pool := p
			if r.Player == Opponent {
				pool = 1 - p
			}
			if pool == player {
				return id, true
			}
		}
	}
	return 0, false
}
