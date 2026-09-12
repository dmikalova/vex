package engine

// maxOffHousePermits bounds how many this-turn off-house grants one player can
// hold at once — generous for real games, where one or two are ever active.
const maxOffHousePermits = 4

// permitUnlimited is the Remaining value for a grant with no card-count bound:
// United Action lets its controller play any number of cards this turn.
const permitUnlimited uint8 = 255

// OffHousePermit is a this-turn grant, made by an effect, that lets a player play
// or use a bounded number of cards outside their active house — the Star Alliance
// "play a non-Star Alliance card" cycle (Com. Officer Kirby, CXO Taber, United
// Action). It is flat, comparable state: the houses and types it frees are named
// by predicate fields, not a closure, so the permit lives in the snapshotable
// GameState (ADR 0005).
type OffHousePermit struct {
	// Except frees every house but this one — a card's own house, so "non-Star
	// Alliance" is Except: StarAlliance. HouseNone excludes no house.
	Except House
	// Controlled frees only houses the player has a card in play for (United Action).
	Controlled bool
	// Types narrows the card types the permit frees — Com. Officer Kirby frees only a
	// non-creature, so Types is artifact, upgrade, and Tactic. The zero value frees
	// every type.
	Types CardTypes
	// Grant is what the permit frees: playing the card from hand (GrantPlay), using
	// the creature in play (GrantUse), or both (CXO Taber plays or uses).
	Grant HouseGrant
	// Remaining is how many plays or uses are left; permitUnlimited is no bound.
	Remaining uint8
}

// frees reports whether the permit still frees a card of the given house and type
// for the player: it has a use left, the house is not excluded, and the type is
// admitted.
func (p OffHousePermit) frees(g *Game, player int, house House, typ CardType) bool {
	if p.Remaining == 0 {
		return false
	}
	if p.Except != HouseNone && house == p.Except {
		return false
	}
	if p.Controlled && !g.controlsHouseInPlay(player, house) {
		return false
	}
	if !p.Types.has(typ) {
		return false
	}
	return true
}

// controlsHouseInPlay reports whether the player has any card of house in play.
func (g *Game) controlsHouseInPlay(player int, house House) bool {
	for _, id := range g.allInPlay(player) {
		if g.House(id) == house {
			return true
		}
	}
	return false
}

// addOffHousePermit records a this-turn off-house grant for the player, dropping
// it if they are already at the cap.
func (g *Game) addOffHousePermit(player int, p OffHousePermit) {
	if int(g.State.OffHousePermitCount[player]) >= maxOffHousePermits {
		return
	}
	g.State.OffHousePermits[player][g.State.OffHousePermitCount[player]] = p
	g.State.OffHousePermitCount[player]++
}

// clearOffHousePermits drops every this-turn off-house grant a player holds,
// called when their turn ends.
func (g *Game) clearOffHousePermits(player int) {
	g.State.OffHousePermits[player] = [maxOffHousePermits]OffHousePermit{}
	g.State.OffHousePermitCount[player] = 0
}

// offHousePlayPermit returns the index of a stored permit that frees playing def
// from hand this turn, or -1 if none does.
func (g *Game) offHousePlayPermit(player int, def *CardDefinition) int {
	for i := 0; i < int(g.State.OffHousePermitCount[player]); i++ {
		p := g.State.OffHousePermits[player][i]
		if p.Grant&GrantPlay != 0 && p.frees(g, player, def.House, def.Type) {
			return i
		}
	}
	return -1
}

// offHouseUsePermit returns the index of a stored permit that frees using the
// creature id this turn, or -1 if none does.
func (g *Game) offHouseUsePermit(player int, id LocalID) int {
	def := g.cat.def(id)
	for i := 0; i < int(g.State.OffHousePermitCount[player]); i++ {
		p := g.State.OffHousePermits[player][i]
		if p.Grant&GrantUse != 0 && p.frees(g, player, g.House(id), def.Type) {
			return i
		}
	}
	return -1
}

// consumeOffHousePermit spends one play or use of the stored permit at index i.
// An unlimited permit is left untouched.
func (g *Game) consumeOffHousePermit(player, i int) {
	if r := g.State.OffHousePermits[player][i].Remaining; r != permitUnlimited {
		g.State.OffHousePermits[player][i].Remaining = r - 1
	}
}

// nonActivePlayLimit is how many off-house plays the player's continuous "any
// non-active house" permissions grant this turn — Captain Val Jericho grants one,
// but only while its centered condition holds.
func (g *Game) nonActivePlayLimit(player int) int {
	limit := 0
	for _, id := range g.allInPlay(player) {
		p := g.cat.def(id).PlayPermission
		if !p.NonActive {
			continue
		}
		if c := p.Condition; c != nil &&
			!c.Met(&EffectContext{Resolver: g, Source: id, Controller: player}) {
			continue
		}
		limit += p.count()
	}
	return limit
}

// nonActivePlayRemaining is how many any-non-active-house plays the player has
// left this turn.
func (g *Game) nonActivePlayRemaining(player int) int {
	used := int(g.State.NonActivePlaysUsedThisTurn[player])
	limit := g.nonActivePlayLimit(player)
	if used >= limit {
		return 0
	}
	return limit - used
}
