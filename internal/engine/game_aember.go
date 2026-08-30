package engine

// gainAember adds Æmber from the common supply to a player's pool. It is the
// single seam for pool gains: before the Æmber reaches the pool, a continuous
// replacement such as Ether Spider may capture the incoming Æmber instead, so the
// pool's owner keeps what they already have. Movements that already come from
// another pool or card (steal, capture, returning captured Æmber) are not gains and
// intentionally bypass this helper.
func (g *Game) gainAember(player, amount int) (LocalID, bool) {
	if capturer, ok := g.aemberCaptorFor(player); ok {
		g.State.Cards[capturer].Amber += int16(amount)
		return capturer, true
	}
	g.State.Aember[player] += amount
	return 0, false
}

// GainAember is the Resolver entry point for gainAember.
func (g *Game) GainAember(player, amount int) (LocalID, bool) { return g.gainAember(player, amount) }

// aemberCaptorFor returns the first in-play creature whose continuous replacement
// captures Æmber that would be added to player's pool, or ok=false when none does.
// A card's replacement names which pool it watches relative to the card's controller
// (Ether Spider watches its Opponent's).
func (g *Game) aemberCaptorFor(player int) (LocalID, bool) {
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
				return id, true
			}
		}
	}
	return 0, false
}
