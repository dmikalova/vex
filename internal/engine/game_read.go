package engine

import "slices"

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
// player 0 can be represented. An attached upgrade has no control of its own — it
// acts through its host — so a controller read on one resolves through the host,
// and generic machinery that asks an upgrade for its controller gets the right
// answer without special-casing upgrades.
func (g *Game) controller(id LocalID) int {
	if host, ok := g.hostOf(id); ok {
		return g.controller(host)
	}
	if c := g.State.Cards[id].ControlPlus; c != 0 {
		return int(c - 1)
	}
	return g.owner(id)
}

// Name returns a card's printed name.
func (g *Game) Name(id LocalID) string { return g.cat.def(id).Name }

// AemberBonus returns the number of Æmber pips printed on a card.
func (g *Game) AemberBonus(id LocalID) int {
	return countBonus(g.cat.def(id).Bonuses, BonusAember)
}

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
		for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
			if h := g.cat.def(up).Static.HouseOverride; h != HouseNone {
				return h
			}
		}
	}
	return g.cat.def(id).House
}

// ActiveHouse returns the house chosen for the current turn.
func (g *Game) ActiveHouse() House { return g.State.ActiveHouse }

// AllowedHouses returns the houses the player may legally choose as their active
// house right now; an empty result means they have no active house this turn and
// must choose "No House". It is the read the client's house picker shares with
// ChooseHouse.
func (g *Game) AllowedHouses(player int) []House { return g.allowedHouses(player) }

// ActivePlayer returns the player whose turn it is.
func (g *Game) ActivePlayer() int { return g.State.ActivePlayer }

// Power returns a creature's current power including attached upgrades.
func (g *Game) Power(id LocalID) int {
	core := &g.State.Cards[id]
	p := g.cat.def(id).Power
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		m := g.staticOn(id, up)
		p += g.scaleStatic(m, m.PowerBonus, id)
	}
	p += int(core.PowerCounters)
	p += int(core.TempPowerBonus)
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
		m := g.staticOn(id, up)
		a += g.scaleStatic(m, m.ArmorBonus, id)
	}
	a += int(g.State.Cards[id].TempArmorBonus)
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
				if b == 0 || !g.constantActive(src, c) || !g.constantAffects(src, c, id) {
					continue
				}
				if c.Per != nil {
					b *= c.Per.Value(g.constantContext(src))
				}
				if c.PerTarget != nil {
					b *= c.PerTarget.perTargetValue(g.constantContext(src), id)
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

// constantActive reports whether constant ability c's positional condition is met
// for its source — a WhileOffFlank ability is suspended while its source holds a
// flank.
func (g *Game) constantActive(src LocalID, c ConstantAbility) bool {
	if g.textBlanked(src) {
		return false
	}
	if c.WhileOffFlank && g.onFlankOf(src) {
		return false
	}
	if c.WhileInCenter && !g.InCenterOfBattleline(src) {
		return false
	}
	if c.WhileCondition != nil && !c.WhileCondition.Met(g.constantContext(src)) {
		return false
	}
	return true
}

// textBlanked reports whether a creature's text box is currently blanked (Shadow
// of Dis), so its printed keywords, abilities, and constant grants are ignored —
// its traits and stats are untouched. Only creatures are blanked.
func (g *Game) textBlanked(id LocalID) bool {
	return g.State.TextBlank[g.controller(id)].Value && g.TypeOf(id) == Creature
}

// assault returns a creature's Assault value including attached upgrades.
func (g *Game) assault(id LocalID) int {
	a := 0
	if !g.textBlanked(id) {
		a = g.cat.def(id).Assault
	}
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		a += g.cat.def(up).Static.AssaultBonus
	}
	a += int(g.State.Cards[id].TempAssaultBonus)
	a += int(g.State.Cards[id].AssaultUntilNextTurn)
	return a
}

// hazardous returns a creature's Hazardous value including attached upgrades.
func (g *Game) hazardous(id LocalID) int {
	h := 0
	if !g.textBlanked(id) {
		h = g.cat.def(id).Hazardous
	}
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		h += g.cat.def(up).Static.HazardousBonus
	}
	h += g.constantBonus(id, func(c ConstantAbility) int { return c.HazardousBonus })
	return h
}

// Hazardous returns a creature's current Hazardous value, including attached
// upgrades.
func (g *Game) Hazardous(id LocalID) int { return g.hazardous(id) }

// splashAttack returns a creature's Splash-attack value including attached
// upgrades.
func (g *Game) splashAttack(id LocalID) int {
	s := 0
	if !g.textBlanked(id) {
		s = g.cat.def(id).SplashAttack
	}
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		s += g.cat.def(up).Static.SplashAttackBonus
	}
	return s
}

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
	if g.State.Cards[id].LostKeywords&k.bit() != 0 {
		return false
	}
	if g.State.Cards[id].LostKeywordsUntilNextTurn&k.bit() != 0 {
		return false
	}
	if !g.textBlanked(id) && g.cat.def(id).hasKeyword(k) {
		return true
	}
	// A creature that has gained another card's text box also has that card's
	// printed keywords (Mimic Gel, Creed of Nurture); a blanked text box ignores
	// them just like its own.
	if !g.textBlanked(id) {
		for _, textSource := range g.grantedTextBoxSources(id) {
			if g.cat.def(textSource).hasKeyword(k) {
				return true
			}
		}
	}
	if g.State.Cards[id].GrantedKeywords&k.bit() != 0 {
		return true
	}
	if g.State.Cards[id].KeywordsUntilNextTurn&k.bit() != 0 {
		return true
	}
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		m := g.staticOn(id, up)
		for _, kw := range m.Keywords {
			if kw == k {
				return true
			}
		}
		for _, grant := range m.KeywordGrants {
			if grant.Host && slices.Contains(grant.Keywords, k) {
				return true
			}
		}
	}
	// A neighbor's upgrade may grant its keywords to this creature — Cloaking
	// Dongle gives Elusive to its host and both of the host's neighbors.
	for _, nb := range neighbors(&EffectContext{Resolver: g}, id) {
		for up, ok := g.firstUpgrade(nb); ok; up, ok = g.nextUpgrade(up) {
			for _, grant := range g.staticOn(nb, up).KeywordGrants {
				if grant.Neighbors && slices.Contains(grant.Keywords, k) {
					return true
				}
			}
		}
	}
	for p := 0; p < 2; p++ {
		for _, src := range g.allInPlay(p) {
			for _, c := range g.cat.def(src).ConstantAbilities {
				for _, kw := range c.Keywords {
					if kw == k && g.constantActive(src, c) && g.constantAffects(src, c, id) {
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

// scaleStatic multiplies a static bonus by the modifier's Per count, read off the
// host creature — Light of the Archons scales +1 power/armor by the number of
// upgrades on the host. A modifier with no Per leaves the bonus as is.
func (g *Game) scaleStatic(m StaticModifier, bonus int, host LocalID) int {
	if bonus == 0 || m.Per == nil {
		return bonus
	}
	return bonus * m.Per.perTargetValue(g.constantContext(host), host)
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

// Enraged reports whether a creature is enraged.
func (g *Game) Enraged(id LocalID) bool { return g.State.Cards[id].Enraged }

// Warded reports whether a creature has a ward.
func (g *Game) Warded(id LocalID) bool { return g.State.Cards[id].Warded }

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

// InBattleline reports whether a creature currently sits on either player's
// battleline, excluding artifacts and cards that have left play.
func (g *Game) InBattleline(id LocalID) bool {
	return g.State.Battleline[0].contains(id) || g.State.Battleline[1].contains(id)
}

// InCenterOfBattleline reports whether a creature sits in the exact center of its
// controller's battleline: the single middle creature of an odd-sized line, with
// equal creatures to its left and right. An even-sized line has no center, so a
// creature there is never centered; a lone creature is its own center.
func (g *Game) InCenterOfBattleline(id LocalID) bool {
	bl := g.State.Battleline[g.controller(id)].slice()
	n := len(bl)
	if n%2 == 0 {
		return false
	}
	return bl[n/2] == id
}

// CurrentlyFighting reports whether a creature is one of the two combatants in the
// fight resolving right now — the flag a "while fighting" self-grant reads (Nizak,
// The Forgotten gains invulnerable). It is false whenever no fight is in progress.
func (g *Game) CurrentlyFighting(id LocalID) bool {
	return g.State.FightersPlus[0] == id+1 || g.State.FightersPlus[1] == id+1
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

// mustFightIfAble reports whether any card in play imposes the global "creatures
// must fight when used, if able" rule (Little Rapscal). It affects both players'
// creatures, so it scans every in-play card.
func (g *Game) mustFightIfAble() bool {
	for owner := 0; owner < 2; owner++ {
		for _, id := range g.allInPlay(owner) {
			if g.cat.def(id).Restricts.MustFightIfAble {
				return true
			}
		}
	}
	return false
}

// cannotReapHouse reports whether a creature is barred from reaping because a
// turn-scoped bar (Seismo-entangler) stops its controller reaping with creatures
// of that creature's house this turn.
func (g *Game) cannotReapHouse(player int, id LocalID) bool {
	bar := g.State.CannotReapHouse[player]
	return bar.Value != HouseNone && bar.Value == g.House(id)
}

// creaturesGloballyBarred reports whether a board-wide "creatures cannot
// fight/reap" bar (Into the Night, Sow Salt) currently stops creature id being
// used the given way. The bar reaches every creature whose controller holds it,
// save for creatures of the one house it spares.
func (g *Game) creaturesGloballyBarred(id LocalID, kind UseKind) bool {
	bar := g.State.CreaturesCannot[g.controller(id)].Value
	if bar.Action != kind {
		return false
	}
	return bar.Houses.matches(&EffectContext{Resolver: g}, id)
}

// cannotReap reports whether a player is barred from reaping — either by the
// timed player-wide bar armed for this turn (Inky Gloom) or by a constant
// Restrictions.Reaping rule on a card in play, their own (Reaping Controller) or
// their opponent's (Barrister Joya's Reaping Opponent, "Enemy creatures cannot
// reap.").
func (g *Game) cannotReap(player int) bool {
	if g.State.CannotReap[player].Value {
		return true
	}
	for owner := 0; owner < 2; owner++ {
		for _, id := range g.allInPlay(owner) {
			r := g.cat.def(id).Restricts.Reaping
			switch r {
			case Controller:
				if player == owner {
					return true
				}
			case Opponent:
				if player != owner {
					return true
				}
			case EachPlayer:
				return true
			}
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

// cannotPlayCreatures reports whether player is barred from playing creatures by a
// constant "cannot play" rule on a card in play — either a Restrictions.CannotPlay
// rule they control or a symmetric CannotPlayWhile bar whose condition holds.
func (g *Game) cannotPlayCreatures(player int) bool {
	for _, id := range g.allInPlay(player) {
		if g.cat.def(id).Restricts.CannotPlay == Creature {
			return true
		}
	}
	return g.barredByConditionalPlayBar(player, Creature)
}

// barredByConditionalPlayBar reports whether any card in play — either player's —
// bars player from playing cards of type t through a CannotPlayWhile rule whose
// condition holds for player (Quixxle Stone bars whoever controls more creatures).
func (g *Game) barredByConditionalPlayBar(player int, t CardType) bool {
	for p := 0; p < 2; p++ {
		for _, id := range g.allInPlay(p) {
			bar := g.cat.def(id).CannotPlayWhile
			if bar.When == nil || bar.Type != t {
				continue
			}
			ctx := &EffectContext{Resolver: g, Source: id, Controller: player}
			if bar.When.Met(ctx) {
				return true
			}
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

// forgeKeyNumberBarred reports whether the next key player would forge is barred
// by a constant Restrictions.NoForgeKeyNumber rule on any card in play — the Key
// Imps bar a key ordinal for both players, whoever controls the Imp.
func (g *Game) forgeKeyNumberBarred(player int) bool {
	next := g.State.Keys[player] + 1
	for p := 0; p < 2; p++ {
		for _, id := range g.allInPlay(p) {
			if g.cat.def(id).Restricts.NoForgeKeyNumber == next {
				return true
			}
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

// aemberProtected reports whether a card player controls makes their Æmber unable
// to be stolen (The Vaultkeeper).
func (g *Game) aemberProtected(player int) bool {
	for _, id := range g.allInPlay(player) {
		if g.cat.def(id).AemberCannotBeStolen {
			return true
		}
		if g.cat.def(id).AemberCannotBeStolenWhileItHasAember && g.AmberOn(id) > 0 {
			return true
		}
		if n := g.cat.def(id).AemberCannotBeStolenWhilePoolAtLeast; n > 0 &&
			g.Aember(player) >= int(n) {
			return true
		}
		for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
			if g.cat.def(up).Static.AemberCannotBeStolen {
				return true
			}
		}
	}
	return false
}

// forgeBarredWhileAhead reports whether player is barred from forging because a
// card in play bars every player who leads on forged keys from forging (Heart of
// the Forest).
func (g *Game) forgeBarredWhileAhead(player int) bool {
	if g.State.Keys[player] <= g.State.Keys[1-player] {
		return false
	}
	for controller := 0; controller < 2; controller++ {
		for _, id := range g.allInPlay(controller) {
			if g.cat.def(id).Restricts.NoForgeWhileAheadOnKeys {
				return true
			}
		}
	}
	return false
}

// choosableHouses returns the houses player p may choose as their active house:
// every house on their Archon identity card, plus the house of every card they
// control in play in its own right — a battleline creature or an artifact — read
// live (ADR 0035, rulebook line 454). A controlled off-house card contributes its
// house; an upgrade or a card under a host does not. When p's identity houses are
// unset (an engine test that declares no deck), every real house is choosable.
func (g *Game) choosableHouses(player int) []House {
	var out []House
	add := func(h House) {
		if h != HouseNone && !slices.Contains(out, h) {
			out = append(out, h)
		}
	}
	if len(g.houses[player]) == 0 {
		for h := Brobnar; int(h) < NumHouses; h++ {
			out = append(out, h)
		}
	} else {
		for _, h := range g.houses[player] {
			add(h)
		}
	}
	for _, id := range g.State.Battleline[player].slice() {
		add(g.House(id))
	}
	for _, id := range g.State.Artifacts[player].slice() {
		add(g.House(id))
	}
	return out
}

// allowedHouses returns the houses player p may legally choose right now, resolving
// their whole constraint table together with every continuous house lock a card in
// play holds over them (Pitlord requires Dis of its controller; Restringuntus bars
// its opponent from a named house). It is the one computation ChooseHouse and the
// client's house picker share (ADR 0035):
//   - Start from the choosable houses and remove every cannot — cannot overrides
//     must, so a house that is both required and barred is barred.
//   - A surviving must is one still choosable and not forbidden. If any survive, the
//     allowed set is exactly those (must A, must B leaves {A, B}; adding cannot A
//     leaves {B}). A must for a house the player cannot choose is void.
//
// An empty result means the player has no active house this turn — a valid outcome,
// not an error (they choose No House).
func (g *Game) allowedHouses(player int) []House {
	choosable := g.choosableHouses(player)
	var cannots, musts []House
	addTo := func(dst *[]House, h House) {
		if h != HouseNone && !slices.Contains(*dst, h) {
			*dst = append(*dst, h)
		}
	}
	for i := 0; i < int(g.State.HouseConstraintCount[player]); i++ {
		c := g.State.HouseConstraints[player][i]
		switch c.Kind {
		case constraintCannotHouse:
			addTo(&cannots, c.House)
		case constraintMustHouse:
			addTo(&musts, c.House)
		case constraintMustCreature:
			addTo(&musts, g.House(c.Creature))
		}
	}
	// Continuous house locks fold into the same must/cannot computation.
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
				addTo(&cannots, locked)
			} else {
				addTo(&musts, locked)
			}
		}
	}
	var allowed []House
	for _, h := range choosable {
		if !slices.Contains(cannots, h) {
			allowed = append(allowed, h)
		}
	}
	var surviving []House
	for _, h := range musts {
		if slices.Contains(allowed, h) {
			surviving = append(surviving, h)
		}
	}
	if len(surviving) != 0 {
		return surviving
	}
	return allowed
}

// keyCostChangeFor returns how much a single in-play card (controlled by
// controller) changes target's key cost — its own change plus any granted by
// attached upgrades.
func (g *Game) keyCostChangeFor(id LocalID, controller, target int) int {
	total := 0
	for _, kc := range g.cat.def(id).KeyCostChanges {
		if kc.affects(controller, target) {
			total += g.keyCostAmount(id, kc)
		}
	}
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		if kc := g.cat.def(up).Static.KeyCostChange; kc.affects(controller, target) {
			total += g.keyCostAmount(id, kc)
		}
		for _, kc := range g.cat.def(up).KeyCostChanges {
			if kc.affects(controller, target) {
				total += g.keyCostAmount(up, kc)
			}
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
	if kc.whileOffFlank && onFlank(g.constantContext(src), src) {
		return 0
	}
	if kc.whileCondition != nil && !kc.whileCondition.Met(g.constantContext(src)) {
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
	if sure := g.State.KeyCostPerHouse[target].Value; sure.House != HouseNone {
		cost += sure.Per * g.creaturesOfHouseInPlay(sure.House)
	}
	for controller := 0; controller < 2; controller++ {
		for _, id := range g.allInPlay(controller) {
			cost += g.keyCostChangeFor(id, controller, target)
		}
	}
	// A key cost can never fall below 0, whatever reductions stack (We Can ALL Win).
	return max(cost, 0)
}

// creaturesOfHouseInPlay counts every creature of the house in play, on either
// battleline — the tally a counted key surcharge (Waking Nightmare) reads live at
// each forge.
func (g *Game) creaturesOfHouseInPlay(house House) int {
	n := 0
	for player := 0; player < 2; player++ {
		for _, id := range g.State.Battleline[player].slice() {
			if g.House(id) == house {
				n++
			}
		}
	}
	return n
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

// entersPlayReady reports whether a card player controls grants a card of type t
// and the given house entry into play ready rather than exhausted (Duskwitch for
// creatures, The Curator for artifacts). The grant is friendly only, so callers
// pass the entering card's controller: a card that somehow enters under the
// opponent's control is not readied by your granter and stays exhausted. A grant
// may be gated on the controller's Æmber pool (Fandangle needs 4) and may
// withhold itself from one house (Fandangle readies only your non-Untamed
// creatures).
func (g *Game) entersPlayReady(player int, t CardType, house House) bool {
	for _, id := range g.allInPlay(player) {
		grant := g.cat.def(id).EntersReadyGrant
		if grant.Type != t {
			continue
		}
		if grant.ExceptHouse != HouseNone && grant.ExceptHouse == house {
			continue
		}
		if grant.MinAember > 0 && g.Aember(player) < int(grant.MinAember) {
			continue
		}
		return true
	}
	return false
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
