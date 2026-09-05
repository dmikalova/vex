package engine

import "math"

// maxCardAember is the most Æmber one card can hold — the range of CardCore.Amber,
// which is narrower than a pool's int because it is paid for 128 times over in the
// flat state. A doubling chain (Binate Rupture) can grow a pool past it, so a
// capture of the whole pool is saturated at this ceiling rather than wrapped.
const maxCardAember = math.MaxInt16

// addAmberOn changes the Æmber sitting on a card, saturating at maxCardAember. It
// clamps only the top: wrapping there would turn a huge pile into negative Æmber
// that later leaks back into a pool when the card leaves play. Going below zero is
// a real bug, so it is left to InvariantError to catch rather than hidden here.
func (g *Game) addAmberOn(id LocalID, delta int) {
	total := int(g.State.Cards[id].Amber) + delta
	if total > maxCardAember {
		g.record(AemberLostToCeiling{Card: id, Amount: total - maxCardAember})
		total = maxCardAember
	}
	g.State.Cards[id].Amber = int16(total)
	// A creature can draw its power from the Æmber sitting on it (Yxili Marauder).
	g.settleDestroyed(g.controller(id))
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
	g.State.Aember[player] += amount
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
			if def.Type != Creature ||
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
