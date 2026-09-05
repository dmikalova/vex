package engine

// This file holds read accessors over the flat GameState: a card's derived stats
// (power, armor, assault, hazardous, keywords — each folding in upgrades and
// constant abilities) and the raw reads of pools, keys, and zone contents. These
// are the reads callers, effects (through the Resolver), and tests share.

// Def returns the read-only definition for an id.
func (g *Game) Def(id LocalID) *CardDefinition { return g.cat.def(id) }

// owner returns the owning player index for an id.
func (g *Game) owner(id LocalID) int { return g.cat.owner(id) }

// controller returns the player currently controlling id. By KeyForge rule,
// ownership is immutable and decides where a card goes out of play; control is
// temporary and is represented by which battleline/artifact row the card occupies.
// ControlPlus uses 0 for "owner controls" and stores controller+1 otherwise so
// player 0 can be represented.
func (g *Game) controller(id LocalID) int {
	if c := g.State.Cards[id].ControlPlus; c != 0 {
		return int(c - 1)
	}
	return g.owner(id)
}

// Name returns a card's printed name.
func (g *Game) Name(id LocalID) string { return g.cat.def(id).Name }

// House returns the house a card currently belongs to. A temporary "belongs to
// house" effect applies only while the card remains in play; everywhere else the
// card keeps its printed house.
func (g *Game) House(id LocalID) House {
	if g.inPlay(id) {
		if h := g.State.Cards[id].TempHouse; h != HouseNone {
			return h
		}
		if h := g.State.Cards[id].LastingHouse; h != HouseNone {
			return h
		}
	}
	return g.cat.def(id).House
}

// ActiveHouse returns the house chosen for the current turn.
func (g *Game) ActiveHouse() House { return g.State.ActiveHouse }

// ActivePlayer returns the player whose turn it is.
func (g *Game) ActivePlayer() int { return g.State.ActivePlayer }

// Power returns a creature's current power including attached upgrades.
func (g *Game) Power(id LocalID) int {
	core := &g.State.Cards[id]
	p := g.cat.def(id).Power
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		p += g.staticOn(id, up).PowerBonus
	}
	p += int(core.PowerCounters)
	p += g.constantBonus(id, func(c ConstantAbility) int { return c.PowerBonus })
	return p
}

// Armor absorbs damage. A creature with armor prevents that much of the damage it
// would be dealt: each point of armor stops 1 damage, and armor spent this way does
// not come back until the creature's controller readies at the end of their turn.
// Armor never reduces a creature's power, and healing does not restore spent armor.
// armor returns a creature's armor value including attached upgrades.
func (g *Game) armor(id LocalID) int {
	a := g.cat.def(id).Armor
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		a += g.staticOn(id, up).ArmorBonus
	}
	a += g.constantBonus(id, func(c ConstantAbility) int { return c.ArmorBonus })
	return a
}

// Armor returns a creature's current armor value, including attached upgrades and
// any constant abilities reaching it.
func (g *Game) Armor(id LocalID) int { return g.armor(id) }

// ArmorStripped returns how much armor an effect has taken off a creature this
// turn. It is not the armor the creature spent absorbing damage: only a strip
// counts, so "for each point of armor it lost this way" measures just this way.
func (g *Game) ArmorStripped(id LocalID) int { return int(g.State.Cards[id].ArmorStripped) }

// constantBonus sums the constant-ability contributions to creature id from every
// card in play, using pick to read the relevant bonus (power or armor) from each
// source's constant ability.
func (g *Game) constantBonus(id LocalID, pick func(ConstantAbility) int) int {
	sum := 0
	for p := 0; p < 2; p++ {
		for _, src := range g.allInPlay(p) {
			for _, c := range g.cat.def(src).ConstantAbilities {
				b := pick(c)
				if b == 0 || !g.constantAffects(src, c, id) {
					continue
				}
				if c.Per != nil {
					b *= c.Per.Value(g.constantContext(src))
				}
				sum += b
			}
		}
	}
	return sum
}

// constantContext is the resolution context a constant ability reads from: its
// own source card, seen by that card's controller.
func (g *Game) constantContext(src LocalID) *EffectContext {
	return &EffectContext{Resolver: g, Source: src, Controller: g.controller(src)}
}

// constantAffects reports whether the constant ability c on source src reaches
// creature id, resolving c's Target from src's point of view.
func (g *Game) constantAffects(src LocalID, c ConstantAbility, id LocalID) bool {
	ctx := g.constantContext(src)
	for _, t := range c.target().Select(ctx) {
		if t == id {
			return true
		}
	}
	return false
}

// assault returns a creature's Assault value including attached upgrades.
func (g *Game) assault(id LocalID) int {
	a := g.cat.def(id).Assault
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		a += g.cat.def(up).Static.AssaultBonus
	}
	return a
}

// hazardous returns a creature's Hazardous value including attached upgrades.
func (g *Game) hazardous(id LocalID) int {
	h := g.cat.def(id).Hazardous
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		h += g.cat.def(up).Static.HazardousBonus
	}
	return h
}

// Hazardous returns a creature's current Hazardous value, including attached
// upgrades.
func (g *Game) Hazardous(id LocalID) int { return g.hazardous(id) }

// ElusiveSpent reports whether a creature has already used its Elusive this turn,
// so it will take combat damage for the rest of the turn. The keybar drops the
// Elusive stripe while this holds.
func (g *Game) ElusiveSpent(id LocalID) bool {
	return g.State.Cards[id].ElusiveUsedThisTurn
}

// hasKeyword reports whether a creature has a keyword, either printed on it,
// granted by an attached upgrade, granted by a card's constant ability, or gained
// for the remainder of the turn (Scout).
func (g *Game) hasKeyword(id LocalID, k Keyword) bool {
	if g.State.KeywordsLost&k.bit() != 0 {
		return false
	}
	if g.cat.def(id).hasKeyword(k) {
		return true
	}
	if g.State.Cards[id].GrantedKeywords&k.bit() != 0 {
		return true
	}
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		for _, kw := range g.staticOn(id, up).Keywords {
			if kw == k {
				return true
			}
		}
	}
	for p := 0; p < 2; p++ {
		for _, src := range g.allInPlay(p) {
			for _, c := range g.cat.def(src).ConstantAbilities {
				for _, kw := range c.Keywords {
					if kw == k && g.constantAffects(src, c, id) {
						return true
					}
				}
			}
		}
	}
	return false
}

// staticOn returns the continuous modifier an attached upgrade currently applies
// to its host, which is the zero modifier while the upgrade's condition is unmet —
// Shoulder Armor gives nothing to a creature that has left the flank.
func (g *Game) staticOn(host, upgrade LocalID) StaticModifier {
	m := g.cat.def(upgrade).Static
	if m.WhileOnFlank && !g.onFlankOf(host) {
		return StaticModifier{}
	}
	return m
}

// Damage returns the damage currently on a creature.
func (g *Game) Damage(id LocalID) int { return int(g.State.Cards[id].Damage) }

// AmberOn returns the Æmber sitting on a card (placed by exalt, capture, etc.).
func (g *Game) AmberOn(id LocalID) int { return int(g.State.Cards[id].Amber) }

// Exhausted reports whether a card is exhausted.
func (g *Game) Exhausted(id LocalID) bool { return g.State.Cards[id].Exhausted }

// InPlay reports whether a card is still on the board.
func (g *Game) InPlay(id LocalID) bool { return g.inPlay(id) }

// Stunned reports whether a creature is stunned.
func (g *Game) Stunned(id LocalID) bool { return g.State.Cards[id].Stunned }

// TimesUsedThisTurn reports how many times this creature has been USED this turn:
// reaped, fought, or had its Action: ability used.
func (g *Game) TimesUsedThisTurn(id LocalID) int {
	return int(g.State.Cards[id].TimesUsedThisTurn)
}

// Aember returns a player's Æmber pool.
func (g *Game) Aember(player int) int { return g.State.Aember[player] }

// AemberProtected is the Resolver entry point for aemberProtected.
func (g *Game) AemberProtected(player int) bool { return g.aemberProtected(player) }

// Keys returns a player's forged key count.
func (g *Game) Keys(player int) int { return g.State.Keys[player] }

// TurnHistory reads one of the tallies the engine keeps about what a player did
// during a turn — keys forged, creatures played, enemies killed fighting. See
// TurnStat for what each one means and when it rolls over.
func (g *Game) TurnHistory(player int, of TurnStat) int {
	return int(g.State.TurnHistory[player][of])
}

// PlayedThisTurn returns the cards a player has played this turn, in play order.
func (g *Game) PlayedThisTurn(player int) []LocalID {
	return cloneIDs(g.State.PlayedThisTurn[player].slice())
}

// DiscardedThisTurn returns the cards a player has discarded from hand this turn,
// in discard order.
func (g *Game) DiscardedThisTurn(player int) []LocalID {
	return cloneIDs(g.State.DiscardedThisTurn[player].slice())
}

// KeyColors returns the colours of the keys a player has forged, in forge order.
func (g *Game) KeyColors(player int) []KeyColor {
	n := g.State.Keys[player]
	out := make([]KeyColor, n)
	copy(out, g.State.KeyColors[player][:n])
	return out
}

// Winner returns the winning player index, or -1 if the game is ongoing.
func (g *Game) Winner() int { return g.State.Winner }

// Hand returns a copy of the ids in a player's hand.
func (g *Game) Hand(player int) []LocalID { return cloneIDs(g.State.Hand[player].slice()) }

// Deck returns a copy of the ids in a player's deck, from top to bottom.
func (g *Game) Deck(player int) []LocalID { return cloneIDs(g.State.Deck[player].slice()) }

// Battleline returns a copy of the ids on a player's battleline.
func (g *Game) Battleline(player int) []LocalID {
	return cloneIDs(g.State.Battleline[player].slice())
}

// Discard returns a copy of the ids in a player's discard pile.
func (g *Game) Discard(player int) []LocalID { return cloneIDs(g.State.Discard[player].slice()) }

// Archives returns a copy of the ids in a player's archives.
func (g *Game) Archives(player int) []LocalID { return cloneIDs(g.State.Archives[player].slice()) }

// Purge returns a copy of the ids a player has purged (set aside out of the game).
func (g *Game) Purge(player int) []LocalID { return cloneIDs(g.State.Purge[player].slice()) }

// Artifacts returns a copy of the ids in a player's artifact row.
func (g *Game) Artifacts(
	player int,
) []LocalID {
	return cloneIDs(g.State.Artifacts[player].slice())
}

// Upgrades returns the ids of upgrades attached to a creature, in attach order.
func (g *Game) Upgrades(id LocalID) []LocalID {
	return g.upgradesOf(id)
}

// Under returns the ids of the cards placed under a host — Masterplan, Jargogle,
// Graft — face up or face down, in the order they were placed. These cards are
// out of play, so they never appear in Battleline, Artifacts, or any other zone
// reader; a host's own Under chain is the only way to reach them.
func (g *Game) Under(host LocalID) []LocalID {
	return g.underOf(host)
}

// UnderFaceDown reports whether a card currently placed under a host is
// facedown, as opposed to faceup (Graft always places its card faceup).
func (g *Game) UnderFaceDown(id LocalID) bool {
	return g.State.Cards[id].UnderFaceDown
}

// Peekable reports whether viewer may look at the front of a facedown card
// placed under host — only the host's controller may (master rulebook,
// FACEDOWN CARDS: a facedown card may only be viewed by the controller of the
// card it is placed under).
func (g *Game) Peekable(viewer int, host LocalID) bool {
	return g.Controller(host) == viewer
}

// inPlay reports whether an id is in either player's battleline or artifact row.
// A controlled creature physically sits in its controller's battleline while its
// owner remains unchanged, so this must not assume owner == controller.
func (g *Game) inPlay(id LocalID) bool {
	for p := 0; p < 2; p++ {
		if g.State.Battleline[p].contains(id) || g.State.Artifacts[p].contains(id) {
			return true
		}
	}
	return false
}

// cannotFight reports whether a player is barred from using creatures to fight,
// by a timed bar (Fogbank) or a constant Restrictions.Fighting rule on a card
// they control in play.
func (g *Game) cannotFight(player int) bool {
	if g.State.CannotFight[player].Value {
		return true
	}
	for _, id := range g.allInPlay(player) {
		if g.cat.def(id).Restricts.Fighting {
			return true
		}
	}
	return false
}

// creaturesPlayedThisTurn counts how many of the cards a player played this turn
// were creatures — the tally the ready phase freezes so the next player can ask
// how many creatures their opponent played on their previous turn (Lifeweb).
func (g *Game) creaturesPlayedThisTurn(player int) int {
	n := 0
	for _, id := range g.PlayedThisTurn(player) {
		if g.cat.def(id).Type == Creature {
			n++
		}
	}
	return n
}

// barredFromPlaying reports whether the timed play bar in force on a player covers
// the card type t, either by naming it or by being the AnyType blanket bar.
func (g *Game) barredFromPlaying(player int, t CardType) bool {
	barred := g.State.CannotPlayTypeThis[player].Value
	return barred == t || barred == AnyType
}

// cannotPlayCreatures reports whether a player is barred from playing creatures by
// a constant Restrictions.CannotPlay rule on a card they control in play.
// cannotPlayCreatures reports whether player is barred from playing creatures by a
// constant "cannot play" rule on a card in play.
func (g *Game) cannotPlayCreatures(player int) bool {
	for _, id := range g.allInPlay(player) {
		if g.cat.def(id).Restricts.CannotPlay == Creature {
			return true
		}
	}
	return false
}

// skipsForge reports whether a player is barred from forging a key by a constant
// Restrictions.SkipForge rule on a card they control in play (The Sting).
func (g *Game) skipsForge(player int) bool {
	for _, id := range g.allInPlay(player) {
		if g.cat.def(id).Restricts.SkipForge {
			return true
		}
	}
	return false
}

// forgeAemberGainer returns the opponent's in-play card that gains payer's forge
// spending (The Sting), and whether one is in play.
func (g *Game) forgeAemberGainer(payer int) (LocalID, bool) {
	for _, id := range g.allInPlay(1 - payer) {
		if g.cat.def(id).GainsForgeAember {
			return id, true
		}
	}
	return 0, false
}

// cannotPlayCard reports whether a player cannot play another card this turn
// because they have reached a card-play limit an in-play card imposes (Ember Imp).
func (g *Game) cannotPlayCard(player int) bool {
	for controller := 0; controller < 2; controller++ {
		for _, id := range g.allInPlay(controller) {
			limit := g.cat.def(id).Restricts.PlayCardLimit
			if limit.Amount > 0 && limit.affects(controller, player) &&
				int(g.State.PlayedThisTurn[player].Count) >= limit.Amount {
				return true
			}
		}
	}
	return false
}

// aemberProtected reports whether a card player controls makes their Æmber immune
// to being stolen (The Vaultkeeper).
func (g *Game) aemberProtected(player int) bool {
	for _, id := range g.allInPlay(player) {
		if g.cat.def(id).PreventSteal {
			return true
		}
	}
	return false
}

// houseLockAllows reports whether player may choose house as their active house,
// given every house lock a card in play holds over them (Pitlord requires Dis of
// its controller, Restringuntus bars its opponent from the house it named). A
// "must" lock only binds when the player actually has that house, matching how a
// forced house yields when it cannot be obeyed; a "cannot" lock always binds.
func (g *Game) houseLockAllows(player int, house House) bool {
	for controller := 0; controller < 2; controller++ {
		for _, id := range g.allInPlay(controller) {
			lock := g.cat.def(id).HouseLock
			if !lock.set() {
				continue
			}
			constrained := controller
			if lock.Player == Opponent {
				constrained = 1 - controller
			}
			if constrained != player {
				continue
			}
			locked := lock.locked(g.State.Cards[id].NamedHouse)
			if locked == HouseNone {
				continue
			}
			if lock.Bars {
				if house == locked {
					return false
				}
			} else if house != locked && g.playerHasHouse(player, locked) {
				return false
			}
		}
	}
	return true
}

// keyCostChangeFor returns how much a single in-play card (controlled by
// controller) changes target's key cost — its own change plus any granted by
// attached upgrades.
func (g *Game) keyCostChangeFor(id LocalID, controller, target int) int {
	total := 0
	if kc := g.cat.def(id).KeyCostChange; kc.affects(controller, target) {
		total += g.keyCostAmount(id, kc)
	}
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		if kc := g.cat.def(up).Static.KeyCostChange; kc.affects(controller, target) {
			total += g.keyCostAmount(id, kc)
		}
	}
	return total
}

// keyCostAmount resolves what a key-cost change on src is currently worth: zero
// while its flank condition is unmet, and scaled by its count when it has one.
func (g *Game) keyCostAmount(src LocalID, kc KeyCostChange) int {
	if kc.whileOnFlank && !onFlank(g.constantContext(src), src) {
		return 0
	}
	if kc.per == nil {
		return kc.amount
	}
	return kc.amount * kc.per.Value(g.constantContext(src))
}

// keyCost returns what a player currently pays to forge one key: the base KeyCost
// plus every key-cost change on a card in play that affects that player.
func (g *Game) keyCost(target int) int {
	cost := KeyCost + g.State.KeyCostBump[target].Value
	for controller := 0; controller < 2; controller++ {
		for _, id := range g.allInPlay(controller) {
			cost += g.keyCostChangeFor(id, controller, target)
		}
	}
	return cost
}

// CurrentKeyCost is the exported view of keyCost: the Æmber a player must spend
// to forge one key right now.
func (g *Game) CurrentKeyCost(player int) int { return g.keyCost(player) }

// battlelineCopy returns a fresh slice of a player's battleline ids, safe to hold
// across state mutations (e.g. while dealing damage to each creature).
func (g *Game) battlelineCopy(player int) []LocalID {
	return cloneIDs(g.State.Battleline[player].slice())
}

// allInPlay returns a fresh slice of a player's creatures and artifacts.
func (g *Game) allInPlay(player int) []LocalID {
	b := g.State.Battleline[player].slice()
	a := g.State.Artifacts[player].slice()
	out := make([]LocalID, 0, len(b)+len(a))
	out = append(out, b...)
	out = append(out, a...)
	return out
}

// cloneIDs copies a slice of ids so callers cannot alias the state arrays.
func cloneIDs(src []LocalID) []LocalID {
	out := make([]LocalID, len(src))
	copy(out, src)
	return out
}

// withoutID returns src with one id dropped, for whittling a candidate list down
// as a repeated choice consumes it.
func withoutID(src []LocalID, drop LocalID) []LocalID {
	out := make([]LocalID, 0, len(src))
	for _, id := range src {
		if id != drop {
			out = append(out, id)
		}
	}
	return out
}
