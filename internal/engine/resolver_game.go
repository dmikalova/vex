package engine

import "sort"

// This file holds the *Game implementation of the Resolver port declared in
// resolver.go: every method a card can reach through EffectContext, grouped by
// the same roles as the interfaces (reads, economy, creature state, combat,
// zones, turn-scoped grants, choices, logging). The interfaces name what an
// effect may do; these bodies carry it out by delegating to the unexported *Game
// internals. See resolver.go for the port and its role interfaces.

// The read accessors Aember, AmberOn, Damage, Power, Name, PlayerName,
// Battleline, and Artifacts are defined in game.go; the remaining Resolver
// methods follow.

// Owner returns the player who owns a card.
func (g *Game) Owner(id LocalID) int { return g.owner(id) }

// Controller returns the player a card in play currently answers to.
func (g *Game) Controller(id LocalID) int { return g.controller(id) }

// HasTrait reports whether a card has a trait, printed, gained through another
// card's text box (Mimic Gel copies a creature, Creed of Nurture lends one), or
// granted until the controller's next turn (the Mutation cycle grants Mutant). A
// RemovesTraits constant in play (Grey Aberrant) strips every trait source, so the
// card then has none.
func (g *Game) HasTrait(id LocalID, trait Trait) bool {
	if g.constantRemovesTraits(id) {
		return false
	}
	if trait != traitUnset && g.State.Cards[id].TraitUntilNextTurn == trait {
		return true
	}
	if g.cat.def(id).hasTrait(trait) {
		return true
	}
	for _, textSource := range g.grantedTextBoxSources(id) {
		if g.cat.def(textSource).hasTrait(trait) {
			return true
		}
	}
	// A creature copying another card's printed stats (Cyber-Clone) gains that card's
	// printed traits.
	if src, ok := g.copiedStatsSource(id); ok && g.cat.def(src).hasTrait(trait) {
		return true
	}
	return false
}

// TraitCount reports how many traits a card has. A RemovesTraits constant in play
// (Grey Aberrant) strips them, so the card then counts zero.
func (g *Game) TraitCount(id LocalID) int {
	if g.constantRemovesTraits(id) {
		return 0
	}
	traits := g.cat.def(id).Traits
	// A creature copying another card's printed stats (Cyber-Clone) gains that card's
	// printed traits, so count each copied trait it does not already print.
	if src, ok := g.copiedStatsSource(id); ok {
		count := len(traits)
		for _, t := range g.cat.def(src).Traits {
			if !g.cat.def(id).hasTrait(t) {
				count++
			}
		}
		return count
	}
	return len(traits)
}

// HasBonusIcons reports whether a card prints at least one bonus icon. A gigantic
// base half prints none itself; while in play it exposes its linked art half's
// icons, so the combined creature counts as having them (ADR 0042).
func (g *Game) HasBonusIcons(id LocalID) bool {
	if len(g.cat.def(id).Bonuses) > 0 {
		return true
	}
	art, ok := g.giganticPartner(id)
	return ok && len(g.cat.def(art).Bonuses) > 0
}

// BonusIconCount reports how many bonus icons a card prints, following a gigantic
// base half to its linked art half (ADR 0042).
func (g *Game) BonusIconCount(id LocalID) int {
	return len(g.bonusIconsOf(id))
}

// SharesTrait reports whether two cards have at least one trait in common. A card
// whose traits are stripped by a RemovesTraits constant in play (Grey Aberrant)
// shares no trait with anything.
func (g *Game) SharesTrait(a, b LocalID) bool {
	if g.constantRemovesTraits(a) || g.constantRemovesTraits(b) {
		return false
	}
	other := g.cat.def(b)
	for _, tr := range g.cat.def(a).Traits {
		if other.hasTrait(tr) {
			return true
		}
	}
	return false
}

// HasKeyword reports whether a creature has a keyword, printed or granted.
func (g *Game) HasKeyword(id LocalID, k Keyword) bool { return g.hasKeyword(id, k) }

// ProtectedByTaunt reports whether target is shielded from attacker by a
// neighboring taunter (the exported CombatResolver port method).
func (g *Game) ProtectedByTaunt(attacker, target LocalID) bool {
	return g.protectedByTaunt(attacker, target)
}

// LoseKeyword takes a keyword away from every creature in play for the remainder
// of the turn (Sniffer).
func (g *Game) LoseKeyword(k Keyword) {
	g.addContinuous(ContinuousEffect{
		Kind:       ContinuousKeywordLost,
		Scope:      ScopeAllCreatures,
		Controller: int8(g.State.ActivePlayer),
		Keywords:   k.bit(),
	}, RemainderOfPlayerTurn)
	g.record(KeywordLostByAll{Keyword: k})
}

// GrantKeyword gives one creature a keyword for the remainder of the turn (Scout).
func (g *Game) GrantKeyword(id LocalID, k Keyword) {
	c := g.stateOf(id)
	if c == nil || c.GrantedKeywords&k.bit() != 0 {
		return
	}
	c.GrantedKeywords |= k.bit()
	g.record(CreatureGainedKeyword{Creature: id, Keyword: k})
}

// LoseKeywordFrom takes a keyword away from one creature for the remainder of the
// turn (Niffle Grounds).
func (g *Game) LoseKeywordFrom(id LocalID, k Keyword) {
	c := g.stateOf(id)
	if c == nil || c.LostKeywords&k.bit() != 0 {
		return
	}
	c.LostKeywords |= k.bit()
	g.record(CreatureLostKeyword{Creature: id, Keyword: k})
}

// LoseKeywordUntilNextTurn takes a keyword away from one creature until the start
// of its controller's next turn, so the loss survives the opponent's turn
// (Reckless Rizzo). A creature destroyed mid-ability has left play, so the loss
// lands on nothing rather than leaving lasting state on a card in the discard pile.
func (g *Game) LoseKeywordUntilNextTurn(id LocalID, k Keyword) {
	c := g.stateOf(id)
	if c == nil || c.LostKeywordsUntilNextTurn&k.bit() != 0 {
		return
	}
	c.LostKeywordsUntilNextTurn |= k.bit()
	g.record(CreatureLostKeyword{Creature: id, Keyword: k})
}

// GrantKeywordUntilNextTurn gives one creature a keyword until the start of its
// controller's next turn (Hideaway Hole).
func (g *Game) GrantKeywordUntilNextTurn(id LocalID, k Keyword) {
	c := g.stateOf(id)
	if c == nil || c.KeywordsUntilNextTurn&k.bit() != 0 {
		return
	}
	c.KeywordsUntilNextTurn |= k.bit()
	g.record(CreatureGainedKeyword{Creature: id, Keyword: k})
}

// ConsideredFlank reports whether a creature counts as a flank creature for the
// turn regardless of its battleline position (Spectral Tunneler).
func (g *Game) ConsideredFlank(id LocalID) bool { return g.State.Cards[id].ConsideredFlank }

// ConsiderFlank makes one creature count as a flank creature for the remainder of
// the turn (Spectral Tunneler).
func (g *Game) ConsiderFlank(id LocalID) {
	if g.State.Cards[id].ConsideredFlank {
		return
	}
	g.State.Cards[id].ConsideredFlank = true
	g.record(CreatureConsideredFlank{Creature: id})
}

// GainStats gives one creature power and/or armor for the remainder of the turn
// (Abond the Armorsmith). The added armor also tops up the armor still available
// to absorb damage this turn, so the extra points can stop damage right away.
func (g *Game) GainStats(id LocalID, power, armor int) {
	c := g.stateOf(id)
	if c == nil {
		return
	}
	c.TempPowerBonus += int16(power)
	c.TempArmorBonus += int16(armor)
	if armor > 0 {
		c.ArmorRemaining += int16(armor)
	}
	g.record(CreatureGainedStats{Creature: id, Power: power, Armor: armor})
}

// GainAssault gives one creature Assault for the remainder of the turn (Creed of
// Nature grants assault equal to a chosen creature's power).
func (g *Game) GainAssault(id LocalID, amount int) {
	c := g.stateOf(id)
	if c == nil {
		return
	}
	c.TempAssaultBonus += int16(amount)
	g.record(CreatureGainedAssault{Creature: id, Amount: amount})
}

// GrantAssaultUntilNextTurn gives one creature Assault until the start of its
// controller's next turn (the Mutation cycle grants assault 3), so it clears with
// the keyword and trait a mutation grants in the same breath.
func (g *Game) GrantAssaultUntilNextTurn(id LocalID, amount int) {
	c := g.stateOf(id)
	if c == nil {
		return
	}
	c.AssaultUntilNextTurn += int16(amount)
	g.record(CreatureGainedAssault{Creature: id, Amount: amount})
}

// GrantTraitUntilNextTurn gives one creature a trait until the start of its
// controller's next turn (the Mutation cycle grants Mutant). It holds a single
// trait, so a second grant to the same creature overwrites the first.
func (g *Game) GrantTraitUntilNextTurn(id LocalID, trait Trait) {
	c := g.stateOf(id)
	if c == nil || c.TraitUntilNextTurn == trait {
		return
	}
	c.TraitUntilNextTurn = trait
	g.record(CreatureGainedTrait{Creature: id, Trait: trait})
}

// ForgeKeyAtExtraCost forges one key at its current cost plus extra and reports
// whether a key was forged, so a forge card purges itself only when it did.
func (g *Game) ForgeKeyAtExtraCost(player, extra int) bool {
	return g.forgeKeyAtExtraCost(player, extra)
}

// ForgeKeyFree forges one key without paying its current cost and reports whether
// a key was forged.
func (g *Game) ForgeKeyFree(player int) bool { return g.forgeKeyFree(player) }

// ForgeKeyFreeForced forges one free key for a player who did not choose to forge,
// with the active player picking its colour, and reports whether a key was forged.
func (g *Game) ForgeKeyFreeForced(player int) bool { return g.forgeKeyFreeForced(player) }

// IsCreature reports whether a card is a creature, by its current type.
func (g *Game) IsCreature(id LocalID) bool { return g.TypeOf(id) == Creature }

// TypeOf returns a card's current type. A card in an upgrade chain reads as an
// Upgrade whatever its printed type — a creature played as an upgrade (ADR 0026)
// is an upgrade while attached. Otherwise an in-play card that converted its type
// (Auto-Legionary turning into a creature) reads its LastingType, and everything
// else reads its printed type.
func (g *Game) TypeOf(id LocalID) CardType {
	if g.State.Cards[id].HostPlus != 0 {
		return Upgrade
	}
	if g.inPlay(id) {
		if t := g.State.Cards[id].LastingType; t != TypeUnset {
			return t
		}
	}
	return g.cat.def(id).Type
}

// GiganticRoleOf returns which half of a gigantic creature a card is, reading its
// printed role.
func (g *Game) GiganticRoleOf(id LocalID) GiganticRole {
	return g.cat.def(id).GiganticRole
}

// SetAember sets a player's Æmber pool, clamped to the range a holder can carry.
// Pool Æmber can feed a
// creature's power (Marmo Swarm gains +1 power per Æmber in its controller's pool),
// so lowering a pool can leave a creature with lethal damage; the resolution
// boundary settles that, not this write (ADR 0029).
func (g *Game) SetAember(player, amount int) {
	if amount < 0 {
		amount = 0
	}
	g.State.Aember[player] = g.clampAember(amount)
}

// NoteAemberStolenFrom adds to the running tally of Æmber stolen from a player
// this turn, clamped to the int8 the history holds, so a later card can ask
// whether their opponent robbed them on their previous turn.
func (g *Game) NoteAemberStolenFrom(player, amount int) {
	if amount <= 0 {
		return
	}
	total := int(g.State.TurnHistory[player][AemberStolenFromThisTurn]) + amount
	if total > 127 {
		total = 127
	}
	g.State.TurnHistory[player][AemberStolenFromThisTurn] = int8(total)
}

// stateOf returns a card's mutable in-play state, or nil once it has left play.
// An ability keeps resolving after its source or target dies — Zyzzix the Many
// adds power counters to itself after its own upgrade's damage destroyed it — and
// KeyForge lets the rest of it resolve, so the parts that need a card on the board
// have to land on nothing. Writing anyway leaves a card in the discard pile
// carrying counters or a stun that nothing will ever clear.
func (g *Game) stateOf(id LocalID) *CardCore {
	if !g.inPlay(id) {
		return nil
	}
	return &g.State.Cards[id]
}

// SetDamage sets the damage on a creature, clamped at zero.
func (g *Game) SetDamage(id LocalID, amount int) {
	if amount < 0 {
		amount = 0
	}
	if c := g.stateOf(id); c != nil {
		c.Damage = int16(amount)
	}
}

// StripArmor empties a creature's remaining armor and tallies what was taken, so
// a following effect can scale with it. It adds to any earlier strip this turn
// rather than replacing it, so two effects that each strip armor both count.
func (g *Game) StripArmor(id LocalID) {
	c := g.stateOf(id)
	if c == nil {
		return
	}
	c.ArmorStripped += c.ArmorRemaining
	c.ArmorRemaining = 0
}

// SetStunned sets a creature's stun status.
func (g *Game) SetStunned(id LocalID, stunned bool) {
	if c := g.stateOf(id); c != nil {
		c.Stunned = stunned
	}
}

// SetEnraged sets a creature's enrage status.
func (g *Game) SetEnraged(id LocalID, enraged bool) {
	if c := g.stateOf(id); c != nil {
		c.Enraged = enraged
	}
}

// SetWarded sets a creature's ward status.
func (g *Game) SetWarded(id LocalID, warded bool) {
	if c := g.stateOf(id); c != nil {
		c.Warded = warded
	}
}

// SetDamageImmune marks a creature unable to be dealt damage for the duration,
// installed as a continuous effect on that single creature.
func (g *Game) SetDamageImmune(id LocalID, d Duration) {
	if g.stateOf(id) == nil {
		return
	}
	g.addContinuous(ContinuousEffect{
		Kind:       ContinuousDamageImmune,
		Scope:      ScopeSubject,
		Subject:    id,
		Controller: int8(g.controller(id)),
	}, d)
}

// SetSideDamageImmune makes every creature player controls unable to be dealt
// damage for the duration (Shield of Justice, Lucky Dice). The mask is read live
// at damage time, so a creature played or taken after this resolves is protected
// too.
func (g *Game) SetSideDamageImmune(player int, d Duration) {
	g.addContinuous(ContinuousEffect{
		Kind:       ContinuousDamageImmune,
		Scope:      ScopeFriendlyCreatures,
		Controller: int8(player),
	}, d)
}

// SetStatOverride masks every creature's power and/or armor to a fixed value for
// the duration (The Pale Star: 1 power, 0 armor). It is read live, so a creature
// that enters later is masked too; the stored counters and armor pool are
// untouched and revealed again when it lifts. An unset StatMask leaves that stat
// alone.
func (g *Game) SetStatOverride(power, armor StatMask, d Duration) {
	g.addContinuous(ContinuousEffect{
		Kind:       ContinuousStatOverride,
		Scope:      ScopeAllCreatures,
		Controller: int8(g.State.ActivePlayer),
		Power:      power.Value,
		Armor:      armor.Value,
		HasPower:   power.Set,
		HasArmor:   armor.Set,
	}, d)
}

// SetExhausted sets a creature's exhausted status.
func (g *Game) SetExhausted(id LocalID, exhausted bool) {
	if c := g.stateOf(id); c != nil {
		c.Exhausted = exhausted
	}
}

// BelongToHouseForRemainderOfTurn makes a card belong to house until its
// controller's turn ends.
func (g *Game) BelongToHouseForRemainderOfTurn(id LocalID, house House) {
	if c := g.stateOf(id); c != nil {
		c.TempHouse = house
	}
}

// SetLastingHouse makes a card belong to house until it leaves play.
func (g *Game) SetLastingHouse(id LocalID, house House) {
	if c := g.stateOf(id); c != nil {
		c.LastingHouse = house
	}
}

// PutIntoBattlelineAsCreature turns an in-play card into a creature and moves it to
// a flank of its controller's battleline. Only artifacts convert this way today
// (Auto-Legionary, Animator), so the card is pulled from the artifact row and
// inserted at the chosen flank; it keeps its exhaustion, Æmber, and power counters,
// and its LastingType makes it read as a creature. Its ArmorRemaining is topped up
// to its full armor so it can absorb hits as a creature this turn. When temporary is
// set the conversion lasts only the current turn (Animator): CreatureUntilTurnEnd
// marks it so the ready phase reverts it to an artifact at end of turn; otherwise it
// stays a creature until it leaves play (Auto-Legionary).
// A repeated use finds the card already a creature in the battleline; removing it
// from there first makes the second use reposition it to the chosen flank rather
// than insert a duplicate.
func (g *Game) PutIntoBattlelineAsCreature(id LocalID, right bool, temporary bool) {
	controller := g.controller(id)
	if !g.State.Artifacts[controller].remove(id) {
		g.State.Battleline[controller].remove(id)
	}
	c := &g.State.Cards[id]
	c.LastingType = Creature
	c.CreatureUntilTurnEnd = temporary
	c.ArmorRemaining = int16(g.armor(id))
	if right {
		g.State.Battleline[controller].add(id)
	} else {
		g.State.Battleline[controller].insertAt(0, id)
	}
	g.record(TurnedIntoCreature{Card: id, Right: right})
}

// revertTemporaryCreatures returns every card that turned into a creature only for
// the current turn (Animator) back to an artifact in its controller's row. It runs
// in the ready phase over both players' cards, since a card can be animated on the
// opponent's turn, so a turn-scoped conversion lifts at end of turn like every other
// RemainderOfPlayerTurn effect.
func (g *Game) revertTemporaryCreatures() {
	var revert []LocalID
	for owner := 0; owner < 2; owner++ {
		for _, id := range g.allInPlay(owner) {
			if g.State.Cards[id].CreatureUntilTurnEnd {
				revert = append(revert, id)
			}
		}
	}
	for _, id := range revert {
		g.revertToArtifact(id)
	}
}

// revertToArtifact moves a temporarily animated card out of its controller's
// battleline and back into their artifact row, dropping the combat state it held
// only as a creature (damage, stun, enrage, ward, armor) but keeping its permanent
// power counters, its Æmber, and its exhaustion.
func (g *Game) revertToArtifact(id LocalID) {
	controller := g.controller(id)
	g.State.Battleline[controller].remove(id)
	c := &g.State.Cards[id]
	c.LastingType = TypeUnset
	c.CreatureUntilTurnEnd = false
	c.Damage = 0
	c.Stunned = false
	c.Enraged = false
	c.Warded = false
	c.ArmorRemaining = 0
	c.ArmorStripped = 0
	g.State.Artifacts[controller].add(id)
	g.record(RevertedToArtifact{Card: id})
}

// SetNamedHouse records the house a card named as it entered play, which its
// HouseLock then constrains for as long as the card stays in play.
func (g *Game) SetNamedHouse(id LocalID, house House) {
	if c := g.stateOf(id); c != nil {
		c.NamedHouse = house
	}
}

// SetNamedCard records the card a permanent named as it entered play, matched by
// name against every card attempted to be played while the permanent stays in play
// (Etan's Jar). It stores the named card's LocalID+1 so the zero value is "named
// nothing"; resetCore clears it when the permanent leaves play.
func (g *Game) SetNamedCard(id, named LocalID) {
	if c := g.stateOf(id); c != nil {
		c.NamedCardPlus = uint8(named) + 1
	}
}

// NameableCards returns one representative card id per distinct card name present
// in the match, ordered by name. It is the match half of a "name a card" choice:
// what the player may name is the whole card database (NameableNames), but only a
// name a card in this match carries can ever bar anything.
func (g *Game) NameableCards() []LocalID {
	seen := make(map[string]bool, len(g.cat.defs))
	out := make([]LocalID, 0, len(g.cat.defs))
	for i := range g.cat.defs {
		id := LocalID(i)
		if name := g.cat.def(id).Name; !seen[name] {
			seen[name] = true
			out = append(out, id)
		}
	}
	sort.Slice(out, func(a, b int) bool {
		return g.cat.def(out[a]).Name < g.cat.def(out[b]).Name
	})
	return out
}

// SetNameableNames injects the names a player may name (Etan's Jar). The match
// passes the whole implemented card database, which the engine cannot read itself
// (ADR 0003). The list is sorted and deduped here so the choice is deterministic
// for replay whatever order it arrives in.
func (g *Game) SetNameableNames(names []string) {
	seen := make(map[string]bool, len(names))
	out := make([]string, 0, len(names))
	for _, n := range names {
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Strings(out)
	g.nameableNames = out
}

// NameableNames returns the names a player may name, falling back to the names
// present in this match when none were injected (a test-built game). Naming from
// the match alone would show a player their opponent's whole deck list, so a real
// match injects the database.
func (g *Game) NameableNames() []string {
	if len(g.nameableNames) > 0 {
		return g.nameableNames
	}
	ids := g.NameableCards()
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = g.cat.def(id).Name
	}
	return out
}

// CardNamed returns a representative card id in this match carrying name, or
// ok=false when no card in the match does.
func (g *Game) CardNamed(name string) (LocalID, bool) {
	for i := range g.cat.defs {
		id := LocalID(i)
		if g.cat.def(id).Name == name {
			return id, true
		}
	}
	return 0, false
}

// SetFightDamageRedirect redirects the attacker's fight damage in the current
// fight to another creature; the combat step reads and clears it.
func (g *Game) SetFightDamageRedirect(id LocalID) { g.State.FightDamageRedirect = id }

// CancelCurrentFight makes the fight in progress not occur; the combat step reads
// and clears it before Assault, Hazardous, and fight damage.
func (g *Game) CancelCurrentFight() { g.State.FightCancelled = true }

// CancelCurrentForge makes the key forge in progress not occur; the before-forge
// window reads and clears it before any Æmber leaves the pool.
func (g *Game) CancelCurrentForge() { g.State.ForgePrevented = true }

// AddAmberOn changes the Æmber sitting on a card.
func (g *Game) AddAmberOn(id LocalID, delta int) { g.addAmberOn(id, delta) }

// DealDamage is the Resolver entry point for the internal dealDamage.
func (g *Game) DealDamage(controller int, targets []DamageTarget) {
	g.dealDamage(controller, targets...)
}

// DestroyEach is the Resolver entry point for destroyEach.
func (g *Game) DestroyEach(controller int, ids []LocalID) { g.destroyEach(controller, ids) }

// DestroyEachFrom credits a source card for the destruction, so the batch it
// targets narrates as one grouped line.
func (g *Game) DestroyEachFrom(controller int, source LocalID, ids []LocalID) {
	prevS, prevH := g.destroyingSource, g.hasDestroyingSource
	g.destroyingSource, g.hasDestroyingSource = source, true
	g.destroyEach(controller, ids)
	g.destroyingSource, g.hasDestroyingSource = prevS, prevH
}

// TakeControl is the Resolver entry point for takeControl.
func (g *Game) TakeControl(id LocalID, controller int, source LocalID) {
	g.takeControl(id, controller, source)
}

// PutIntoPlay is the Resolver entry point for putIntoPlay.
func (g *Game) PutIntoPlay(id LocalID, controller int) {
	g.putIntoPlay(id, controller)
}

// PlayerHasHouse reports whether house is one of the player's identity houses.
func (g *Game) PlayerHasHouse(player int, house House) bool {
	return g.playerHasHouse(player, house)
}

// Draw is the Resolver entry point for the internal draw.
func (g *Game) Draw(controller, count int) {
	if n := g.draw(controller, count); n > 0 {
		g.record(CardsDrawnBy{Player: controller, Count: n})
	}
}

// RefillHand refills a player's hand as if it were the end of their turn,
// honoring their chains and draw modifiers (Punctuated Equilibrium).
func (g *Game) RefillHand(player int) { g.drawStep(player) }

// PutOnTopOfDeck is the Resolver entry point for putOnTopOfDeck.
func (g *Game) PutOnTopOfDeck(id LocalID) { g.putOnTopOfDeck(id) }

// PutIntoHand is the Resolver entry point for putIntoHand.
func (g *Game) PutIntoHand(id LocalID) { g.putIntoHand(id) }

// ReturnUpgradesToHand is the Resolver entry point for returnUpgradesToHand.
func (g *Game) ReturnUpgradesToHand(host LocalID) { g.returnUpgradesToHand(host) }

// Simultaneously is the Resolver entry point for simultaneously.
func (g *Game) Simultaneously(controller int, batch func()) {
	g.simultaneously(controller, batch)
}

// ArchiveUpgrade is the Resolver entry point for archiveUpgrade.
func (g *Game) ArchiveUpgrade(upgrade LocalID) { g.archiveUpgrade(upgrade) }

// PutIntoArchives is the Resolver entry point for putIntoArchives.
func (g *Game) PutIntoArchives(id LocalID) { g.putIntoArchives(id) }

// PutIntoArchivesEach is the Resolver entry point for putIntoArchivesEach.
func (g *Game) PutIntoArchivesEach(controller int, ids []LocalID) {
	g.putIntoArchivesEach(controller, ids)
}

// PutIntoDeckShuffled is the Resolver entry point for putIntoDeckShuffled.
func (g *Game) PutIntoDeckShuffled(id LocalID) { g.putIntoDeckShuffled(id) }

// ShuffleFriendlyCardsInPlayIntoDeck is the Resolver entry point for
// shuffleFriendlyInPlayIntoDeck.
func (g *Game) ShuffleFriendlyCardsInPlayIntoDeck(player int) [2]int {
	return g.shuffleFriendlyInPlayIntoDeck(player)
}

// BeginShuffleBatch opens a shuffle batch: cards shuffled into a deck until
// EndShuffleBatch are collected rather than narrated one by one.
func (g *Game) BeginShuffleBatch() {
	g.shuffleBatch, g.batchingShuffle = nil, true
}

// EndShuffleBatch closes the batch and narrates the collected cards grouped by
// owner as one CardsShuffledIntoDeckBy line each, attributed to the frame's source,
// in the order the owners were first shuffled.
func (g *Game) EndShuffleBatch() {
	batch := g.shuffleBatch
	g.shuffleBatch, g.batchingShuffle = nil, false
	byOwner := map[int][]LocalID{}
	var owners []int
	for _, id := range batch {
		o := g.owner(id)
		if _, seen := byOwner[o]; !seen {
			owners = append(owners, o)
		}
		byOwner[o] = append(byOwner[o], id)
	}
	for _, o := range owners {
		g.record(CardsShuffledIntoDeckBy{Owner: o, Cards: byOwner[o]})
	}
}

// ArchiveFromHand moves a card from its owner's hand to their archives.
func (g *Game) ArchiveFromHand(id LocalID) { g.archiveFromHand(g.owner(id), id) }

// ArchiveEnemyFromHand moves a card from its owner's hand into player's archives —
// an abduction from hand (Hidden Stash).
func (g *Game) ArchiveEnemyFromHand(player int, id LocalID) { g.archiveEnemyFromHand(player, id) }

// ArchiveFromPurge moves a card from a player's purge pile to their archives.
func (g *Game) ArchiveFromPurge(owner int, id LocalID) { g.archiveFromPurge(owner, id) }

// ArchiveFromDiscard moves a card from a player's discard pile to their archives.
func (g *Game) ArchiveFromDiscard(owner int, id LocalID) { g.archiveFromDiscard(owner, id) }

// DiscardArchives moves all of a player's archived cards to their discard pile.
func (g *Game) DiscardArchives(owner int) { g.discardArchives(owner) }

// PurgeFromDiscard moves a card from a player's discard pile to their purge pile.
func (g *Game) PurgeFromDiscard(owner int, id LocalID) { g.purgeFromDiscard(owner, id) }

// PurgeFromHand moves a card from a player's hand to their purge pile.
func (g *Game) PurgeFromHand(owner int, id LocalID) { g.purgeFromHand(owner, id) }

// PurgeFromArchives moves a card from a player's archives to their purge pile.
func (g *Game) PurgeFromArchives(owner int, id LocalID) { g.purgeFromArchives(owner, id) }

// PurgeFromDeck moves a card from a player's deck to their purge pile.
func (g *Game) PurgeFromDeck(owner int, id LocalID) { g.purgeFromDeck(owner, id) }

// PurgeFromPlay is the Resolver entry point for purgeFromPlay.
func (g *Game) PurgeFromPlay(id LocalID) { g.purgeFromPlay(id) }

// RedirectResolvingCard sends a card whose play is still resolving to dest when
// that play completes, instead of to its owner's discard pile (Sucker Punch
// archives itself, Library Access purges itself). It writes through stateOf's
// in-play guard because a resolving card is deliberately in no zone at all.
func (g *Game) RedirectResolvingCard(id LocalID, dest Destination) {
	g.State.Cards[id].ResolvingDest = dest
}

// AddPowerCounter changes the net power counters on a creature. A -1 counter can
// lower power to the damage already marked; the resolution boundary settles that,
// not this write (ADR 0029).
func (g *Game) AddPowerCounter(id LocalID, delta int) {
	if c := g.stateOf(id); c != nil {
		c.PowerCounters += int16(delta)
	}
}

// PutFromDiscardIntoHand moves a card from its owner's discard pile to their hand.
func (g *Game) PutFromDiscardIntoHand(id LocalID) {
	o := g.owner(id)
	g.moveOwnCard(id, o, Discard, Hand, CardPutFromDiscardIntoHand{Player: o, Card: id})
}

// MoveFromDeckToHand moves a card from its owner's deck to their hand.
func (g *Game) MoveFromDeckToHand(id LocalID) {
	o := g.owner(id)
	g.moveOwnCard(id, o, Deck, Hand, CardPutFromDeckIntoHand{Player: o, Card: id})
}

// MoveFromDeckToDiscard moves a card from its owner's deck to their discard pile.
func (g *Game) MoveFromDeckToDiscard(id LocalID) {
	o := g.owner(id)
	g.moveOwnCard(id, o, Deck, Discard, CardMoved{Player: o, Card: id, From: Deck, To: Discard})
}

// moveOwnCard moves a card between two of the same player's resting zones.
func (g *Game) moveOwnCard(id LocalID, player int, from, to Zone, entry LogEntry) {
	g.moveCard(id,
		zoneRef{Player: player, Zone: from},
		zoneRef{Player: player, Zone: to},
		entry)
}

// ArchiveFromDeck moves a card from its owner's deck to their archives.
func (g *Game) ArchiveFromDeck(id LocalID) { g.archiveFromDeck(g.owner(id), id) }

// PutDeckCardOnBottom moves a card from its owner's deck to the bottom of that
// same deck. The move is private information, so it records no log line.
func (g *Game) PutDeckCardOnBottom(id LocalID) {
	d := &g.State.Deck[g.owner(id)]
	d.remove(id)
	d.add(id)
}

// SetDeckTop rewrites the top len(order) cards of player's deck to order, with
// order[0] on top. The reorder is private information, so it records no log line.
func (g *Game) SetDeckTop(player int, order []LocalID) {
	d := &g.State.Deck[player]
	copy(d.IDs[:], order)
}

// ShuffleZonesIntoDeck moves each named zone's cards into a player's deck and
// shuffles once.
func (g *Game) ShuffleZonesIntoDeck(player int, zones []Zone) {
	rec := ShuffledIntoDeck{Player: player}
	for _, z := range zones {
		switch z {
		case Hand:
			rec.HandCount += int(g.State.Hand[player].Count)
		case Archives:
			rec.ArchivesCount += int(g.State.Archives[player].Count)
		default: // Discard
			rec.DiscardCards = append(rec.DiscardCards, g.State.Discard[player].slice()...)
		}
	}
	g.shuffleZonesIntoDeck(player, zones)
	g.record(rec)
}

// MoveFromDiscardToTopOfDeck moves a card from its owner's discard pile to the
// top of their deck.
func (g *Game) MoveFromDiscardToTopOfDeck(id LocalID) {
	o := g.owner(id)
	g.State.Discard[o].remove(id)
	g.State.Deck[o].addFront(id)
	g.record(CardPutFromDiscardOnTopOfDeck{Player: o, Card: id})
}

// MoveFromDeckToTopOfDeck repositions a card already in its owner's deck to the
// top of that deck.
func (g *Game) MoveFromDeckToTopOfDeck(id LocalID) {
	o := g.owner(id)
	g.State.Deck[o].remove(id)
	g.State.Deck[o].addFront(id)
	g.record(CardPutOnTopOfDeck{Card: id, Owner: o})
}

// shuffleIntoDeckFrom moves a card out of one of its owner's zones into their
// deck and shuffles. During a shuffle batch the card is collected for a single
// grouped narration rather than narrated on its own (Not Finished with You), so
// every source zone narrates the same way.
func (g *Game) shuffleIntoDeckFrom(id LocalID, from *deckList) {
	o := g.owner(id)
	from.remove(id)
	g.State.Deck[o].add(id)
	g.Shuffle(o)
	if g.batchingShuffle {
		g.shuffleBatch = append(g.shuffleBatch, id)
		return
	}
	g.record(CardShuffledIntoDeck{Card: id, Owner: o})
}

// ShuffleFromDiscardIntoDeck moves a card from its owner's discard pile into their
// deck and shuffles.
func (g *Game) ShuffleFromDiscardIntoDeck(id LocalID) {
	g.shuffleIntoDeckFrom(id, &g.State.Discard[g.owner(id)])
}

// ShuffleFromHandIntoDeck moves a card from its owner's hand into their deck and
// shuffles.
func (g *Game) ShuffleFromHandIntoDeck(id LocalID) {
	g.shuffleIntoDeckFrom(id, &g.State.Hand[g.owner(id)])
}

// GainChains adds chains to a player, which reduce their draws until shed.
func (g *Game) GainChains(controller, amount int) {
	g.State.Chains[controller] += amount
	g.record(ChainsGained{
		Player: controller,
		Amount: amount,
		Total:  g.State.Chains[controller],
	})
}

// OrderByChoice is the Resolver entry point for orderByChoice.
func (g *Game) OrderByChoice(controller int, prompt string, ids []LocalID) []LocalID {
	return g.orderByChoice(controller, prompt, ids)
}

// ChooseCreature asks a player to choose one creature from candidates, attributing
// the prompt to the source card. A sole candidate is taken automatically.
func (g *Game) ChooseCreature(
	player int,
	source LocalID,
	prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	return g.pickCreature(player, g.sourceName(source), prompt, candidates)
}

// ChooseCard asks a player to choose one card from candidates, attributing the
// prompt to the source card. A sole candidate is taken automatically.
func (g *Game) ChooseCard(
	player int,
	source LocalID,
	prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	return g.pickCard(player, g.sourceName(source), prompt, candidates)
}

// ChooseCardOptional asks a player to choose one card from candidates or to
// decline, attributing the prompt to the source card. A sole candidate is still
// offered rather than forced, because declining is a legal answer.
func (g *Game) ChooseCardOptional(
	player int,
	source LocalID,
	prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	return g.pickOptional(player, g.sourceName(source), prompt, candidates)
}

// ChooseOption asks a player to choose one of several labeled options, attributing
// (does not implement OptionChooser), the first option is taken.
func (g *Game) ChooseOption(player int, source LocalID, prompt string, options []string) int {
	return g.chooseOption(player, g.sourceName(source), prompt, options)
}

// ChooseRandom picks one uniformly random card from candidates using the game's
// RNG, reporting ok=false for an empty slice. It is the shared draw behind a
// Random selection.
func (g *Game) ChooseRandom(candidates []LocalID) (LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	return candidates[g.State.PRNG.Intn(len(candidates))], true
}

// PreviewBadge forwards a selection-badge hint to a player's client if it can
// show one (implements BadgeChooser). It is display-only, so a chooser without
// the capability ignores it.
func (g *Game) PreviewBadge(player int, badge SelectionBadge) {
	if bc, ok := g.chooserFor(player).(BadgeChooser); ok {
		bc.PreviewBadge(badge)
	}
}

// chooseOption is the shared option-choice path: it attributes the prompt to a
// source name (empty for a source-less prompt such as a turn-structure choice)
// and defaults to the first option when the chooser has no preference. A sole
// option is taken automatically without consulting the chooser.
func (g *Game) chooseOption(player int, source, prompt string, options []string) int {
	if len(options) == 1 {
		return 0
	}
	// Boundary: settle before presenting the choice (ADR 0029).
	g.settleDestroyed(player)
	if oc, ok := g.chooserFor(player).(OptionChooser); ok {
		return oc.ChooseOption(source, renderPrompt(source, prompt), options)
	}
	return 0
}

// sourceName returns the name of a source card for prompt attribution, or "" when
// the id is not a registered card (e.g. an unset source in a unit test).
func (g *Game) sourceName(source LocalID) string {
	if int(source) < len(g.cat.defs) {
		return g.cat.defs[source].Name
	}
	return ""
}

// FightWith makes attacker fight defender, ability-driven (ignoring active player
// and house). A creature can only be used while ready and while its name is under
// the Rule of Six, so an exhausted or six-times-used attacker does nothing.
func (g *Game) FightWith(attacker, defender LocalID) {
	if g.usableByAbility(attacker) {
		g.fight(attacker, defender)
	}
}

// ReapWith reaps with a creature, ability-driven (ignoring active player and
// house). A creature can only be used while ready and while its name is under the
// Rule of Six, so an exhausted or six-times-used creature does nothing.
func (g *Game) ReapWith(id LocalID) {
	if g.usableByAbility(id) {
		g.reapWith(id)
	}
}

// UseActionOf fires a card's "Action:" ability on behalf of actor, ability-driven
// (ignoring active player and house). A card can only be used while ready and
// while its name is under the Rule of Six, so an exhausted or six-times-used card
// does nothing.
func (g *Game) UseActionOf(actor int, id LocalID) {
	if g.usableByAbility(id) {
		g.useActionOf(actor, id)
	}
}

// TriggerAbilityOf resolves a card's abilities under one trigger on behalf of
// actor. Unlike UseActionOf this does not use the card, so readiness is beside
// the point: an exhausted creature's reap effect still triggers, and the creature
// neither exhausts nor counts as used.
func (g *Game) TriggerAbilityOf(actor int, id LocalID, trigger Trigger) {
	g.triggerAbilitiesAs(actor, id, trigger, 0, false)
}

// TriggerAbilityOfRooted is TriggerAbilityOf carrying the Rule-of-Six cascade
// root, so every ability the chain resolves charges root's name pool rather than
// each resolving card's own.
func (g *Game) TriggerAbilityOfRooted(actor int, id LocalID, trigger Trigger, root LocalID) {
	g.triggerAbilitiesRooted(actor, id, trigger, 0, false, root, true)
}

// AtRuleOfSix reports whether id's card name has been used six times this turn.
func (g *Game) AtRuleOfSix(id LocalID) bool { return g.atRuleOfSix(id) }

// RecordUsage counts one usage of id toward its card name's Rule-of-Six pool.
func (g *Game) RecordUsage(id LocalID) { g.recordUsage(id) }

// HasTrigger reports whether a card has an ability under the trigger.
func (g *Game) HasTrigger(id LocalID, trigger Trigger) bool {
	return g.hasTrigger(id, trigger)
}

// Record appends one narrated outcome to the game log (ADR 0011).
func (g *Game) Record(e LogEntry) { g.record(e) }
