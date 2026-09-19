package engine

import (
	"fmt"
	"iter"
	"slices"
)

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

// DeckList returns the printed names of every card player owns, deduplicated with
// an "xN" count and sorted. It reads the whole catalog rather than the live zones,
// so it lists a player's entire deck — the cards still in the draw pile included —
// which is what a debug replay needs to spot the card behind an invariant even
// when that card never reached the game log.
func (g *Game) DeckList(player int) []string {
	counts := map[string]int{}
	for id := 0; id < g.cat.count(); id++ {
		if g.cat.owner(LocalID(id)) == player {
			counts[g.cat.def(LocalID(id)).Name]++
		}
	}
	names := make([]string, 0, len(counts))
	for name := range counts {
		names = append(names, name)
	}
	slices.Sort(names)
	for i, name := range names {
		if counts[name] > 1 {
			names[i] = fmt.Sprintf("%s x%d", name, counts[name])
		}
	}
	return names
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
	// A stat override masks the real total: The Pale Star makes each creature have 1
	// power, read live and ignoring counters and bonuses, revealed again when it
	// lifts (the stored counters are untouched).
	if v, ok := g.continuousPowerOverride(id); ok {
		return v
	}
	// A creature copying another card's printed stats (Cyber-Clone) has that card's
	// printed power outright, ignoring counters and bonuses just like a live override.
	if src, ok := g.copiedStatsSource(id); ok {
		return g.cat.def(src).Power
	}
	// A variable "X" power reads its neighbors' power, so two such creatures that
	// reference each other would recurse forever. A creature already mid-computation
	// contributes 0 — its power is undeterminable, which KeyForge treats as 0.
	if slices.Contains(g.powerComputing, id) {
		return 0
	}
	g.powerComputing = append(g.powerComputing, id)
	defer func() { g.powerComputing = g.powerComputing[:len(g.powerComputing)-1] }()
	core := &g.State.Cards[id]
	p := g.cat.def(id).Power
	p += g.upgradeStatBonus(id, func(m StaticModifier) int { return m.PowerBonus })
	p += int(core.PowerCounters)
	p += int(core.TempPowerBonus)
	p += g.constantBonus(id, func(c ConstantAbility) int { return c.PowerBonus })
	// A variable "X" power (Picaroon's combined-neighbor power) is a blank-able part
	// of the card's text, so it contributes only while the text is not blanked.
	if px := g.cat.def(id).PowerX; px != nil && !g.textBlanked(id) {
		p += px.Value(&EffectContext{
			Resolver:   g,
			Source:     id,
			Controller: g.controller(id),
		})
	}
	return p
}

// PowerCountersOn returns the net +1/-1 power counters placed on a creature — the
// running tally Chonkers doubles. It is the raw counter total, not the creature's
// power.
func (g *Game) PowerCountersOn(id LocalID) int {
	return int(g.State.Cards[id].PowerCounters)
}

// Armor absorbs damage. A creature with armor prevents that much of the damage it
// would be dealt: each point of armor stops 1 damage, and armor spent this way does
// not come back until the creature's controller readies at the end of their turn.
// Armor never reduces a creature's power, and healing does not restore spent armor.
// armor returns a creature's armor value including attached upgrades.
func (g *Game) armor(id LocalID) int {
	a := g.cat.def(id).Armor
	a += g.upgradeStatBonus(id, func(m StaticModifier) int { return m.ArmorBonus })
	a += int(g.State.Cards[id].TempArmorBonus)
	a += g.constantBonus(id, func(c ConstantAbility) int { return c.ArmorBonus })
	// A creature copying another card's printed stats (Cyber-Clone) gains that card's
	// printed armor on top of its own.
	if src, ok := g.copiedStatsSource(id); ok {
		a += g.cat.def(src).Armor
	}
	return a
}

// Armor returns a creature's current armor value, including attached upgrades and
// any constant abilities reaching it. A stat override masks it (The Pale Star sets
// each creature to 0 armor); the real armor pool is untouched and revealed again
// when the mask lifts.
func (g *Game) Armor(id LocalID) int {
	if v, ok := g.continuousArmorOverride(id); ok {
		return v
	}
	return g.armor(id)
}

// ArmorStripped returns how much armor an effect has taken off a creature this
// turn. It is not the armor the creature spent absorbing damage: only a strip
// counts, so "for each point of armor it lost this way" measures just this way.
func (g *Game) ArmorStripped(id LocalID) int { return int(g.State.Cards[id].ArmorStripped) }

// constantAbilitiesInPlay iterates every constant ability printed on a card either
// player has in play, yielding the card that prints it and the ability itself.
// ConstantAbility carries a dozen independent capabilities and each reader asks
// about one of them, so this walk is shared by a dozen callers across the engine.
//
// It deliberately applies no gates. The gates are what the callers disagree on:
// most want constantActive and constantAffects, a board-wide rule such as
// triggerDisabled skips constantAffects because it changes the rules for everyone,
// and constantBlanksText must not call constantActive at all (that would recurse
// through textBlanked — see its doc comment). Folding the gates in here would make
// those differences invisible, so they stay written out at each call site.
func (g *Game) constantAbilitiesInPlay() iter.Seq2[LocalID, ConstantAbility] {
	return func(yield func(LocalID, ConstantAbility) bool) {
		for p := range 2 {
			for src, c := range g.constantAbilitiesOf(p) {
				if !yield(src, c) {
					return
				}
			}
		}
	}
}

// constantAbilitiesOf is constantAbilitiesInPlay narrowed to the cards one player
// has in play, for the reads whose rule is per-player rather than board-wide.
func (g *Game) constantAbilitiesOf(player int) iter.Seq2[LocalID, ConstantAbility] {
	return func(yield func(LocalID, ConstantAbility) bool) {
		for _, src := range g.allInPlay(player) {
			for _, c := range g.cat.def(src).ConstantAbilities {
				if !yield(src, c) {
					return
				}
			}
		}
	}
}

// anyActiveConstant reports whether a constant ability in play both satisfies want
// and reaches id. This is the shape of every yes/no constant-ability read — the
// caller supplies only the capability it cares about, and the two standard gates
// are applied here.
func (g *Game) anyActiveConstant(id LocalID, want func(ConstantAbility) bool) bool {
	for src, c := range g.constantAbilitiesInPlay() {
		if want(c) && g.constantActive(src, c) && g.constantAffects(src, c, id) {
			return true
		}
	}
	return false
}

// constantBonus sums what the constant abilities in play add to one numeric stat
// of id, pick reading that stat off each ability.
func (g *Game) constantBonus(id LocalID, pick func(ConstantAbility) int) int {
	sum := 0
	for src, c := range g.constantAbilitiesInPlay() {
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
	return sum
}

// constantContext is the resolution context a constant ability reads from: its
// own source card, seen by that card's controller.
func (g *Game) constantContext(src LocalID) *EffectContext {
	return &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: g.controller(src),
	}
}

// constantAffects reports whether the constant ability c on source src reaches
// creature id, resolving c's Target from src's point of view.
func (g *Game) constantAffects(src LocalID, c ConstantAbility, id LocalID) bool {
	ctx := g.constantContext(src)
	return slices.Contains(c.target().Select(ctx), id)
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
// of Dis): its printed keywords, abilities, and constant grants are ignored while
// its traits and stats remain. A duration-scoped blank in the continuous registry
// reaches creatures; a while-in-play ConstantAbility.BlankText reaches artifacts
// (Blossom Drake) — this predicate consults both.
func (g *Game) textBlanked(id LocalID) bool {
	return g.continuousActive(ContinuousTextBlank, id) || g.constantBlanksText(id)
}

// constantBlanksText reports whether a card in play with a BlankText constant
// ability reaches card id. A source blanked itself (by the registry) grants
// nothing, so its blank is skipped — a constant BlankText reaches artifacts, never
// a creature source, so the registry check alone breaks any recursion.
func (g *Game) constantBlanksText(id LocalID) bool {
	for src, c := range g.constantAbilitiesInPlay() {
		if !c.BlankText || g.continuousActive(ContinuousTextBlank, src) {
			continue
		}
		if g.constantAffects(src, c, id) {
			return true
		}
	}
	return false
}

// constantRemovesTraits reports whether a card in play with a RemovesTraits
// constant ability reaches card id, stripping its traits while the source stays
// in play (Grey Aberrant). It mirrors constantBlanksText but honours the source's
// own active gates via constantActive (blank, off-flank, WhileCondition), since a
// trait-removing source is itself a creature that can be blanked or conditioned.
func (g *Game) constantRemovesTraits(id LocalID) bool {
	return g.anyActiveConstant(id, func(c ConstantAbility) bool { return c.RemovesTraits })
}

// constantSelectiveArchivePickup reports whether player controls an active card in
// play whose constant ability grants the selective archive pickup (The Archivist),
// returning that source card. Like triggerDisabled it reads player-wide, ignoring
// Target: the source's mere presence changes the pickup rule for its controller.
func (g *Game) constantSelectiveArchivePickup(player int) (LocalID, bool) {
	for src, c := range g.constantAbilitiesOf(player) {
		if c.SelectiveArchivePickup && g.constantActive(src, c) {
			return src, true
		}
	}
	return 0, false
}

// assault returns a creature's Assault value including attached upgrades.
func (g *Game) assault(id LocalID) int {
	a := 0
	if !g.textBlanked(id) {
		a = g.cat.def(id).Assault
	}
	a += g.upgradeStatBonus(id, func(m StaticModifier) int { return m.AssaultBonus })
	a += int(g.State.Cards[id].TempAssaultBonus)
	a += int(g.State.Cards[id].AssaultUntilNextTurn)
	a += g.constantBonus(id, func(c ConstantAbility) int { return c.AssaultBonus })
	return a
}

// hazardous returns a creature's Hazardous value including attached upgrades.
func (g *Game) hazardous(id LocalID) int {
	h := 0
	if !g.textBlanked(id) {
		h = g.cat.def(id).Hazardous
	}
	h += g.upgradeStatBonus(id, func(m StaticModifier) int { return m.HazardousBonus })
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
	s += g.upgradeStatBonus(id, func(m StaticModifier) int { return m.SplashAttackBonus })
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
	if g.keywordLost(id, k) {
		return false
	}
	return g.keywordFromText(id, k) ||
		g.keywordGranted(id, k) ||
		g.keywordFromUpgrades(id, k) ||
		g.keywordFromNeighborUpgrades(id, k) ||
		g.keywordFromConstantAbilities(id, k)
}

// keywordLost reports whether a keyword is currently taken away from a creature —
// by a continuous "loses its keywords" effect, or a lasting one lasting this turn
// or until the creature's next turn.
func (g *Game) keywordLost(id LocalID, k Keyword) bool {
	return g.continuousLostKeywords(id)&k.bit() != 0 ||
		g.State.Cards[id].LostKeywords&k.bit() != 0 ||
		g.State.Cards[id].LostKeywordsUntilNextTurn&k.bit() != 0
}

// keywordFromText reports whether a keyword comes from a creature's own printed
// text: its own keyword line, a card whose text box it has gained (Mimic Gel, Creed
// of Nurture), or a card whose printed stats it copies (Cyber-Clone). A blanked
// text box ignores all three.
func (g *Game) keywordFromText(id LocalID, k Keyword) bool {
	if g.textBlanked(id) {
		return false
	}
	if g.cat.def(id).hasKeyword(k) {
		return true
	}
	for _, textSource := range g.grantedTextBoxSources(id) {
		if g.cat.def(textSource).hasKeyword(k) {
			return true
		}
	}
	if src, ok := g.copiedStatsSource(id); ok && g.cat.def(src).hasKeyword(k) {
		return true
	}
	return false
}

// keywordGranted reports whether a keyword was gained for the remainder of the turn
// or until the creature's next turn (Scout).
func (g *Game) keywordGranted(id LocalID, k Keyword) bool {
	return g.State.Cards[id].GrantedKeywords&k.bit() != 0 ||
		g.State.Cards[id].KeywordsUntilNextTurn&k.bit() != 0
}

// keywordFromUpgrades reports whether an upgrade attached to the creature grants
// the keyword, either as a plain static keyword or a host-directed keyword grant.
func (g *Game) keywordFromUpgrades(id LocalID, k Keyword) bool {
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		m := g.staticOn(id, up)
		if slices.Contains(m.Keywords, k) {
			return true
		}
		for _, grant := range m.KeywordGrants {
			if grant.Host && slices.Contains(grant.Keywords, k) {
				return true
			}
		}
	}
	return false
}

// keywordFromNeighborUpgrades reports whether a neighbor's upgrade grants the
// keyword to this creature — Cloaking Dongle gives Elusive to its host and both of
// the host's neighbors.
func (g *Game) keywordFromNeighborUpgrades(id LocalID, k Keyword) bool {
	for _, nb := range neighbors(&EffectContext{Resolver: g}, id) {
		for up, ok := g.firstUpgrade(nb); ok; up, ok = g.nextUpgrade(up) {
			for _, grant := range g.staticOn(nb, up).KeywordGrants {
				if grant.Neighbors && slices.Contains(grant.Keywords, k) {
					return true
				}
			}
		}
	}
	return false
}

// keywordFromConstantAbilities reports whether any in-play card's constant ability
// grants the keyword to this creature while active and affecting it.
func (g *Game) keywordFromConstantAbilities(id LocalID, k Keyword) bool {
	return g.anyActiveConstant(id, func(c ConstantAbility) bool {
		return slices.Contains(c.Keywords, k)
	})
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

// upgradeStatBonus sums what the upgrades attached to host contribute to one stat,
// pick reading that stat off each modifier. Every stat goes through here so none
// can skip a modifier's own gates: staticOn suspends a WhileOnFlank upgrade off the
// flank, and scaleStatic applies its Per count. Assault, Hazardous, and Splash used
// to read Static directly and silently ignored both
// (TestUpgradeStatBonusHonoursWhileOnFlankForEveryStat).
func (g *Game) upgradeStatBonus(host LocalID, pick func(StaticModifier) int) int {
	sum := 0
	for up, ok := g.firstUpgrade(host); ok; up, ok = g.nextUpgrade(up) {
		m := g.staticOn(host, up)
		sum += g.scaleStatic(m, pick(m), host)
	}
	return sum
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
func (g *Game) Aember(player int) int { return int(g.State.Aember[player]) }

// AemberProtected is the Resolver entry point for aemberProtected.
func (g *Game) AemberProtected(player int) bool { return g.aemberProtected(player) }

// Keys returns a player's forged key count.
func (g *Game) Keys(player int) int { return g.State.KeyCount(player) }

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
	n := g.Keys(player)
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

// ZoneOf reports which out-of-play pile a card sits in and whose it is, so a
// caller need not scan every pile itself to place a card a prompt reaches into.
// It checks each player's discard, hand, archives, deck, and purge piles in turn;
// a card in play or not found returns ok false.
func (g *Game) ZoneOf(id LocalID) (player int, zone Zone, ok bool) {
	for p := range 2 {
		for _, z := range []struct {
			zone Zone
			ids  []LocalID
		}{
			{Discard, g.State.Discard[p].slice()},
			{Hand, g.State.Hand[p].slice()},
			{Archives, g.State.Archives[p].slice()},
			{Deck, g.State.Deck[p].slice()},
			{Purged, g.State.Purge[p].slice()},
		} {
			if slices.Contains(z.ids, id) {
				return p, z.zone, true
			}
		}
	}
	return 0, zoneUnset, false
}

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

// inPlay reports whether an id is in play: in either player's battleline or
// artifact row, or attached as an upgrade to a card that is. An attached upgrade
// sits on its host's chain rather than in a row, but it is a card in play in its
// own right, so every reachability re-check must see it (TestInPlayCountsUpgrades).
// A controlled creature physically sits in its controller's battleline while its
// owner remains unchanged, so this must not assume owner == controller.
func (g *Game) inPlay(id LocalID) bool {
	if host, ok := g.hostOf(id); ok {
		id = host
	}
	for p := range 2 {
		if g.State.Battleline[p].contains(id) || g.State.Artifacts[p].contains(id) {
			return true
		}
	}
	return false
}

// activeInDiscard reports whether a card resolves an ability from a discard pile:
// it carries TriggersFromDiscard and currently sits in its owner's discard pile
// (Relentless Creeper returns itself to hand after its controller chooses Dis).
func (g *Game) activeInDiscard(id LocalID) bool {
	return g.cat.def(id).TriggersFromDiscard && g.State.Discard[g.owner(id)].contains(id)
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
	for owner := range 2 {
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
	for owner := range 2 {
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

// barredByNamedCard reports whether a card named by a permanent in play (Etan's
// Jar) bars def from being played. The bar is symmetric — it applies to whichever
// player attempts the play — and matches by name, so it covers every copy of the
// named card in either deck. It lifts when the naming permanent leaves play,
// because resetCore clears the stored name.
func (g *Game) barredByNamedCard(def *CardDefinition) bool {
	for p := range 2 {
		for _, id := range g.allInPlay(p) {
			named := g.State.Cards[id].NamedCardPlus
			if named != 0 && g.cat.def(LocalID(named-1)).Name == def.Name {
				return true
			}
		}
	}
	return false
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
	for p := range 2 {
		for _, id := range g.allInPlay(p) {
			bar := g.cat.def(id).CannotPlayWhile
			if bar.When == nil || bar.Type != t {
				continue
			}
			ctx := &EffectContext{
				Resolver:   g,
				Source:     id,
				Controller: player,
			}
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

// keyForgeCapReached reports whether player already holds the maximum number of
// keys (MaxKeys), so a "forge a key" trigger that fires after the fourth forge
// cannot push the count past the KeyColors slots.
func (g *Game) keyForgeCapReached(player int) bool {
	return g.Keys(player) >= MaxKeys
}

// forgeKeyNumberBarred reports whether the next key player would forge is barred
// by a constant Restrictions.NoForgeKeyNumber rule on any card in play — the Key
// Imps bar a key ordinal for both players, whoever controls the Imp.
func (g *Game) forgeKeyNumberBarred(player int) bool {
	next := g.Keys(player) + 1
	for p := range 2 {
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
	for controller := range 2 {
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
		ctx := &EffectContext{
			Resolver:   g,
			Source:     id,
			Controller: player,
		}
		if c := g.cat.def(id).AemberCannotBeStolen; c != nil && c.Met(ctx) {
			return true
		}
		for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
			if c := g.cat.def(up).Static.AemberCannotBeStolen; c != nil && c.Met(ctx) {
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
	if g.Keys(player) <= g.Keys(1-player) {
		return false
	}
	for controller := range 2 {
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
	for _, id := range g.allInPlay(player) {
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
	cannots, musts := g.houseConstraintLists(player)
	var allowed []House
	for _, h := range g.choosableHouses(player) {
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

// houseConstraintLists gathers the houses player is barred from (cannots) and
// required to pick from (musts), from both the explicit house constraints and the
// continuous house locks of in-play cards. Both fold into the same two lists.
func (g *Game) houseConstraintLists(player int) (cannots, musts []House) {
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
	for controller := range 2 {
		for _, id := range g.allInPlay(controller) {
			if h, bars, ok := g.lockedHouse(id, controller, player); ok {
				if bars {
					addTo(&cannots, h)
				} else {
					addTo(&musts, h)
				}
			}
		}
	}
	return cannots, musts
}

// lockedHouse returns the house an in-play card locks for player, and whether that
// lock bars the house (a cannot) or requires it (a must). ok is false when the card
// carries no house lock, its lock does not constrain player, or it locks no house.
func (g *Game) lockedHouse(id LocalID, controller, player int) (house House, bars, ok bool) {
	lock := g.cat.def(id).HouseLock
	if !lock.set() {
		return HouseNone, false, false
	}
	constrained := controller
	if lock.Player == Opponent {
		constrained = 1 - controller
	}
	if constrained != player {
		return HouseNone, false, false
	}
	locked := lock.locked(g.State.Cards[id].NamedHouse)
	if locked == HouseNone {
		return HouseNone, false, false
	}
	return locked, lock.Bars, true
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
	if sure := g.State.KeyCostPerHouse[target].Value; sure.Per != 0 {
		cost += sure.Per * g.creaturesMatchingInPlay(sure.House)
	}
	for controller := range 2 {
		for _, id := range g.allInPlay(controller) {
			cost += g.keyCostChangeFor(id, controller, target)
		}
	}
	// A key cost can never fall below 0, whatever reductions stack (We Can ALL Win).
	return max(cost, 0)
}

// creaturesMatchingInPlay counts every creature the matcher admits in play, on
// either battleline — the tally a counted key surcharge (Waking Nightmare) reads
// live at each forge.
func (g *Game) creaturesMatchingInPlay(m HouseMatcher) int {
	ctx := &EffectContext{Resolver: g}
	n := 0
	for player := range 2 {
		for _, id := range g.State.Battleline[player].slice() {
			if m.matches(ctx, id) {
				n++
			}
		}
	}
	return n
}

// CurrentKeyCost is the exported view of keyCost: the Æmber a player must spend
// to forge one key right now.
func (g *Game) CurrentKeyCost(player int) int { return g.keyCost(player) }

// AtCheck reports whether a pool of Æmber meets a player's current key cost —
// enough to forge a key, KeyForge's "Check!". The amount is passed in so a live
// standing checks the player's pool while the game log checks a recorded amount.
func (g *Game) AtCheck(player, aember int) bool {
	return aember >= g.CurrentKeyCost(player)
}

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
