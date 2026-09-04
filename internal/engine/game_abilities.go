package engine

import "slices"

// This file holds how a card is USED — the player actions that reap, fight, or
// activate an "Action:" ability — the checks that gate them, and the machinery
// that fires a card's triggered abilities. Combat resolution lives in
// game_combat.go and destruction in game_destroy.go.

// usableInActiveHouse reports whether the in-play card id may be USED (reaped,
// fought, or its action ability activated) under the active house. It holds when
// the card is in the active house or is Versatile — printed on it or granted by
// an attached upgrade — letting it be used as if it belonged to the active house.
// Versatile only relaxes using: a Versatile card is still played from hand only
// when its own house is active.
func (g *Game) usableInActiveHouse(id LocalID) bool {
	return g.manual ||
		g.State.ActiveHouse == HouseNone ||
		g.House(id) == g.State.ActiveHouse ||
		g.hasKeyword(id, Versatile) ||
		(g.State.MayUseHouse[g.controller(id)] != HouseNone && g.House(id) == g.State.MayUseHouse[g.controller(id)])
}

// usable runs the checks shared by reaping, fighting, and using an action
// ability: the card must be controlled by the active player, in play, unexhausted,
// and usable under the active house. It does not restrict by card type — callers
// add that.
func (g *Game) usable(player int, id LocalID) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	if g.controller(id) != player || !g.inPlay(id) {
		return ErrWrongType
	}
	if g.State.CannotUse[player].Value {
		return ErrCannotUse
	}
	if g.State.Cards[id].Exhausted {
		return ErrCardExhausted
	}
	if !g.usableInActiveHouse(id) {
		return ErrWrongHouse
	}
	def := g.cat.def(id)
	if def.Restricts.UseCondition != nil {
		ctx := &EffectContext{Resolver: g, Source: id, Controller: player}
		if !def.Restricts.UseCondition.Met(ctx) {
			return ErrCannotUse
		}
	}
	return nil
}

// canUse validates that a creature may be used to reap or fight right now.
func (g *Game) canUse(player int, id LocalID) error {
	if err := g.usable(player, id); err != nil {
		return err
	}
	if g.cat.def(id).Type != Creature {
		return ErrWrongType
	}
	if !g.hasAnyUse(player, id) {
		return ErrCannotUse
	}
	return nil
}

// hasAnyUse reports whether at least one of the three ways to use a creature is
// open to it, so a card barred from some of them (Tireless Crocag cannot reap) is
// still offered while another way remains. It keeps CanUse's promise honest: a
// Crocag with nothing to fight has no use at all this turn.
func (g *Game) hasAnyUse(player int, id LocalID) bool {
	if !g.cannotBeUsedTo(id, ReapUse) {
		return true
	}
	if !g.cannotBeUsedTo(id, FightUse) &&
		!g.cannotFight(player) &&
		len(g.Battleline(1-player)) > 0 {
		return true
	}
	return !g.cannotBeUsedTo(id, ActionUse) && g.hasTrigger(id, TriggerAction)
}

// canUseTo is canUse for one specific way of using the creature, so a card that
// bars only one of them (Tireless Crocag cannot reap) stays usable the other ways.
func (g *Game) canUseTo(player int, id LocalID, kind UseKind) error {
	if err := g.canUse(player, id); err != nil {
		return err
	}
	if g.cannotBeUsedTo(id, kind) {
		return ErrCannotUse
	}
	return nil
}

// CanUse reports whether a creature may currently be used (reap/fight/action) by
// the player: nil if usable, otherwise the reason. UIs can call it to reject an
// action before prompting for a target.
func (g *Game) CanUse(player int, id LocalID) error { return g.canUse(player, id) }

// CanUseTo reports whether a creature may currently be used one specific way, so
// a UI can offer Fight but not Reap on a creature that cannot reap.
func (g *Game) CanUseTo(player int, id LocalID, kind UseKind) error {
	return g.canUseTo(player, id, kind)
}

// CanUseArtifact reports whether an artifact's "Action:" ability may currently be
// used by the player: nil if usable, otherwise the reason. It shares usable's
// house check, so a Versatile artifact (Lifeward) is correctly offered out of
// the active house rather than a UI reimplementing the house check on its own.
func (g *Game) CanUseArtifact(player int, id LocalID) error {
	if err := g.usable(player, id); err != nil {
		return err
	}
	if g.cat.def(id).Type != Artifact {
		return ErrWrongType
	}
	if !g.hasTrigger(id, TriggerAction) {
		return ErrCannotUse
	}
	return nil
}

// Reap uses a creature to reap, gaining 1 Æmber and firing "Reap:" abilities.
func (g *Game) Reap(player int, id LocalID) error {
	if err := g.canUseTo(player, id, ReapUse); err != nil {
		return err
	}
	g.reapWith(id)
	return nil
}

// Unstun spends a stunned creature's use for the turn shaking off the stun
// instead — the only thing a stunned creature can do with a use, so it is
// offered wherever Reap/Fight/an action would be, under the same checks
// (including the active-house one). A creature that is *forced* into an
// action rather than choosing one skips this check entirely; reapWith,
// fight, and useActionOf each absorb that forced use into the stun recovery
// on their own, with no active-player or house check of their own.
func (g *Game) Unstun(player int, id LocalID) error {
	if err := g.usable(player, id); err != nil {
		return err
	}
	if g.cat.def(id).Type != Creature || !g.State.Cards[id].Stunned {
		return ErrCannotUse
	}
	g.reapWith(id)
	return nil
}

// recordUse marks one successful use of a creature this turn. "Used" is the
// rulebook umbrella for reaping, fighting, or using an Action: ability; the count
// advances before that use resolves so Fight:/Reap: abilities see their first use
// as count 1.
func (g *Game) recordUse(id LocalID) {
	if g.cat.def(id).Type == Creature {
		g.State.Cards[id].TimesUsedThisTurn++
	}
}

// reapWith performs a reap driven by a rule or ability, with no active-player or
// active-house checks: a stunned creature recovers instead; otherwise it
// exhausts, its controller gains 1 Æmber, and its "Reap:" abilities fire. A card
// that cannot reap is not made to by an ability either.
func (g *Game) reapWith(id LocalID) {
	if g.cannotBeUsedTo(id, ReapUse) {
		return
	}
	g.recordUse(id)
	if g.recoverFromStun(id) {
		return
	}
	p := g.controller(id)
	g.State.Cards[id].Exhausted = true
	g.gainReapAember(p, id)
	g.triggerAbilities(id, TriggerAfterReap, 0, false)
	g.emitCardUsed(p, id)
	g.emitLasting(EventReap, p, id)
}

// gainReapAember pays out the Æmber a reap grants to player p, applying any lasting
// replacement of that payout (Dimension Door makes it steal instead of gain).
func (g *Game) gainReapAember(p int, source LocalID) {
	if act, ok := g.lastingReplacement(p, EventReapAember); ok && act == actSteal {
		stolen := min(1, g.State.Aember[1-p])
		g.State.Aember[1-p] -= stolen
		g.State.Aember[p] += stolen
		g.record(ReapedStealing{Player: p, Card: source, Amount: stolen})
		return
	}
	if capturer, ok := g.gainAember(p, 1); ok {
		g.record(ReapedCaptured{Player: p, Card: source, Creature: capturer})
		return
	}
	g.record(Reaped{Player: p, Card: source})
}

// UseAction uses a creature's or artifact's "Action:" ability.
func (g *Game) UseAction(player int, id LocalID) error {
	if err := g.usable(player, id); err != nil {
		return err
	}
	if g.cannotBeUsedTo(id, ActionUse) {
		return ErrCannotUse
	}
	if !g.hasTrigger(id, TriggerAction) {
		return ErrWrongType
	}
	if g.cat.def(id).Type == Artifact {
		if err := g.chargeToll(player, TollUseArtifact); err != nil {
			return err
		}
	}
	g.useActionOf(player, id)
	return nil
}

// useActionOf fires a card's "Action:" ability on behalf of actor, driven by a
// rule or ability: a stunned card recovers instead; otherwise it exhausts and its
// "Action:" fires. actor is normally the card's own controller, but an ability
// that uses a card "as if it were yours" (Remote Access) passes itself instead,
// so the Action resolves for the user rather than the card's controller.
func (g *Game) useActionOf(actor int, id LocalID) {
	g.recordUse(id)
	if g.recoverFromStun(id) {
		return
	}
	g.State.Cards[id].Exhausted = true
	g.record(ActionAbilityUsed{Player: actor, Card: id})
	g.triggerAbilitiesAs(actor, id, TriggerAction, 0, false)
	g.emitCardUsed(actor, id)
}

// Fight uses attacker to fight the enemy creature defender.
// fightErrorForgiven reports whether a canUse error may be forgiven for a fight.
// A grant such as Brothers in Battle lets a creature of a chosen house fight out
// of the active house, so a wrong-house error is excused when the grant covers
// this attacker; every other check still applies.
func (g *Game) fightErrorForgiven(err error, attacker LocalID) bool {
	return err == ErrWrongHouse && g.mayFightOutOfHouse(attacker)
}

// Fight uses a creature to fight an enemy creature: it validates the attacker and
// the target, then resolves simultaneous combat. It returns an error when the
// attacker cannot be used (readiness, or wrong house unless a grant forgives it)
// or the target is not a legal enemy creature.
func (g *Game) Fight(player int, attacker, defender LocalID) error {
	if g.cannotFight(player) {
		return ErrCannotFight
	}
	if err := g.canUseTo(player, attacker, FightUse); err != nil &&
		!g.fightErrorForgiven(err, attacker) {
		return err
	}
	if g.cat.def(defender).Type != Creature ||
		g.controller(defender) == player ||
		!g.inPlay(defender) ||
		g.protectedByTaunt(attacker, defender) {
		return ErrNoTarget
	}
	if fr := g.cat.def(attacker).FightRestriction; fr != (Target{}) &&
		!fr.allows(&EffectContext{Resolver: g, Source: attacker, Controller: player}, defender) {
		return ErrNoTarget
	}
	g.fight(attacker, defender)
	return nil
}

// hasTrigger reports whether an in-play card has the trigger itself, from an
// attached upgrade, or from a constant ability affecting it.
func (g *Game) hasTrigger(id LocalID, trigger Trigger) bool {
	if g.cat.def(id).hasTrigger(trigger) {
		return true
	}
	for _, upgrade := range g.Upgrades(id) {
		for _, ab := range g.cat.def(upgrade).Static.Granted {
			if ab.Trigger == trigger {
				return true
			}
		}
	}
	for player := 0; player < 2; player++ {
		for _, grantor := range g.allInPlay(player) {
			for _, c := range g.cat.def(grantor).ConstantAbilities {
				if !g.constantAffects(grantor, c, id) {
					continue
				}
				for _, ab := range c.Granted {
					if ab.Trigger == trigger {
						return true
					}
				}
			}
		}
	}
	return false
}

// mayFightOutOfHouse reports whether a fight grant (Brothers in Battle) lets the
// attacker fight this turn despite not being in the active house.
func (g *Game) mayFightOutOfHouse(attacker LocalID) bool {
	p := g.controller(attacker)
	if g.State.MayFightAny[p] {
		return true
	}
	h := g.State.MayFightHouse[p]
	return h != HouseNone && g.House(attacker) == h
}

// FightTargets returns the enemy creatures the attacker may legally fight right
// now, mirroring the checks in Fight. It is empty when the player is barred from
// fighting, the attacker cannot be used (readiness, or wrong house without a
// fight grant), or no enemy creature satisfies the attacker's fight restriction.
// A UI offers Fight only when this is non-empty.
func (g *Game) FightTargets(player int, attacker LocalID) []LocalID {
	if g.cannotFight(player) {
		return nil
	}
	if err := g.canUse(player, attacker); err != nil && !g.fightErrorForgiven(err, attacker) {
		return nil
	}
	fr := g.cat.def(attacker).FightRestriction
	var targets []LocalID
	for _, def := range g.State.Battleline[1-player].slice() {
		if g.protectedByTaunt(attacker, def) {
			continue
		}
		if fr != (Target{}) &&
			!fr.allows(&EffectContext{Resolver: g, Source: attacker, Controller: player}, def) {
			continue
		}
		targets = append(targets, def)
	}
	return targets
}

// recoverFromStun handles using a stunned creature: it exhausts and clears the
// stun instead of performing the reap/fight/action. It reports whether the
// creature was stunned, in which case the caller should stop (the use is spent
// removing the stun and nothing else happens).
func (g *Game) recoverFromStun(id LocalID) bool {
	core := &g.State.Cards[id]
	if !core.Stunned {
		return false
	}
	core.Stunned = false
	core.Exhausted = true
	g.record(StunRecovered{Creature: id})
	return true
}

// readyToUse reports whether a creature may be used by an ability right now. A
// creature can only be used while ready: an ability may still target an exhausted
// creature, but using it then does nothing. Unlike canUse this ignores the active
// player and active house — ability-driven use cares only about readiness.
func (g *Game) readyToUse(id LocalID) bool {
	if g.State.Cards[id].Exhausted {
		g.record(CardCannotBeUsed{Card: id})
		return false
	}
	return true
}

// resolveUpgradePlay resolves an upgrade's own "Play:" abilities when it is attached,
// acting on its host creature. The upgrade is the card that was played, but its
// ability speaks of "this creature" — the host — so the host is the effect source
// that self-references resolve to.
func (g *Game) resolveUpgradePlay(host, upgrade LocalID, up *CardDefinition) {
	for _, ab := range up.Abilities {
		if ab.Trigger != TriggerAfterPlay {
			continue
		}
		// The ability's own text is not a log line — it is attribution (ADR 0011), so
		// the outcomes the ability produces carry it and narrate themselves.
		closeFrame := g.openFrame(Frame{
			Actor:      g.owner(upgrade),
			Source:     host,
			HasSource:  true,
			Trigger:    TriggerAfterPlay,
			Grantor:    upgrade,
			HasGrantor: true,
		})
		ab.Effect.Resolve(&EffectContext{
			Resolver:   g,
			Source:     host,
			Upgrade:    upgrade,
			Controller: g.owner(upgrade),
		})
		closeFrame()
	}
}

// emitCreatureEnters is the enter-play event for a creature. It first resolves the
// entering creature's own "enters play" abilities (Chuff Ape entering stunned),
// then fires "after a creature enters play" abilities on every other in-play card,
// with the entering creature as the trigger target ("it").
func (g *Game) emitCreatureEnters(entered LocalID) {
	g.triggerAbilities(entered, TriggerEntersPlay, 0, false)
	for player := 0; player < 2; player++ {
		for _, id := range g.allInPlay(player) {
			if id == entered {
				continue
			}
			g.triggerAbilities(id, TriggerAfterCreatureEnters, entered, true)
		}
	}
}

// emitEnemyDestroyed fires the persistent "after an enemy creature is destroyed
// during your turn" reaction (Pile of Skulls) on the active player's in-play
// cards. It fires only for the active player, and only when the destroyed
// creature is one of their enemies, so the reaction is naturally limited to your
// own turn and to enemy creatures.
func (g *Game) emitEnemyDestroyed(destroyed LocalID) {
	active := g.State.ActivePlayer
	if g.controller(destroyed) == active {
		return
	}
	for _, id := range g.allInPlay(active) {
		g.triggerAbilities(id, TriggerAfterEnemyCreatureDestroyed, destroyed, true)
	}
}

// emitCardPlayed fires "after you play a card" abilities on the playing player's
// other in-play cards, with the played card as "it", and the EventCardPlayed
// lasting reactions the player has armed (Library Access). Only an actual play from
// hand fires it; a card put into play by another effect enters (emitCreatureEnters)
// but is not played.
func (g *Game) emitCardPlayed(player int, played LocalID) {
	for _, id := range g.allInPlay(player) {
		if id == played {
			continue
		}
		g.triggerAbilities(id, TriggerAfterCardPlayed, played, true)
	}
	for _, id := range g.allInPlay(1 - player) {
		g.triggerAbilities(id, TriggerAfterEnemyCardPlayed, played, true)
	}
	g.emitLasting(EventCardPlayed, player, played)
}

// emitCardUsed fires "after you use a card" abilities on the user's other in-play
// cards, with the used card as "it". Every route to using a card — reaping,
// fighting, or an "Action:" — ends here, so a reaction to using never has to be
// wired into each verb separately.
func (g *Game) emitCardUsed(player int, used LocalID) {
	for _, id := range g.allInPlay(player) {
		if id == used {
			continue
		}
		g.triggerAbilities(id, TriggerAfterUse, used, true)
	}
}

// triggerAbilities resolves every ability matching the trigger that the card
// carries itself, is granted by an attached upgrade, or is granted by an in-play
// card's constant ability. Destroyed abilities do not come through here: they
// share one simultaneous window across every dying creature, so destroyTogether
// owns them (ADR 0013), including the rule that a creature taken out of play
// mid-window drops its remaining abilities.
func (g *Game) triggerAbilities(src LocalID, trigger Trigger, it LocalID, hasIt bool) {
	g.triggerAbilitiesAs(g.controller(src), src, trigger, it, hasIt)
}

// triggerAbilitiesAs is triggerAbilities with the resolving controller named
// explicitly, for an ability used on behalf of someone other than its controller.
func (g *Game) triggerAbilitiesAs(
	actor int,
	src LocalID,
	trigger Trigger,
	it LocalID,
	hasIt bool,
) {
	for _, t := range g.orderTriggered(actor, g.triggeredBy(src, trigger)) {
		closeFrame := g.openFrame(Frame{
			Actor:      actor,
			Source:     src,
			HasSource:  true,
			Trigger:    trigger,
			Grantor:    t.grantor,
			HasGrantor: t.grantor != src,
		})
		t.ability.Effect.Resolve(&EffectContext{
			Resolver:   g,
			Source:     src,
			Controller: actor,
			It:         it,
			HasIt:      hasIt,
		})
		closeFrame()
	}
}

// triggeredBy collects every ability matching the trigger that src carries itself,
// is granted by an attached upgrade, or is granted by an in-play card's constant
// ability. Collection happens before any of them resolves so the whole window can
// be ordered (ADR 0013).
func (g *Game) triggeredBy(src LocalID, trigger Trigger) []triggeredAbility {
	var pending []triggeredAbility
	keep := func(grantor LocalID, ab Ability) {
		if ab.Trigger == trigger {
			pending = append(pending, triggeredAbility{source: src, grantor: grantor, ability: ab})
		}
	}
	for _, ab := range g.cat.def(src).Abilities {
		keep(src, ab)
	}
	for up, ok := g.firstUpgrade(src); ok; up, ok = g.nextUpgrade(up) {
		for _, ab := range g.cat.def(up).Static.Granted {
			keep(up, ab)
		}
	}
	for player := 0; player < 2; player++ {
		for _, grantor := range g.allInPlay(player) {
			for _, c := range g.cat.def(grantor).ConstantAbilities {
				if len(c.Granted) == 0 || !g.constantAffects(grantor, c, src) {
					continue
				}
				for _, ab := range c.Granted {
					keep(grantor, ab)
				}
			}
		}
	}
	return pending
}

// orderTriggerPrompt is the prompt shown when abilities on several cards trigger
// at once and the active player must say which card's resolves next.
const orderTriggerPrompt = "Choose which card's ability resolves next"

// orderGrantorPrompt is the prompt shown when one card has several abilities in
// the same trigger window (its own text plus one an upgrade or a constant ability
// granted it) and the active player must say which resolves next.
const orderGrantorPrompt = "Choose which ability resolves next"

// orderTriggered lets the active player arrange a trigger window's abilities into
// a resolution order (ADR 0013). Ordering runs at two levels, because the Chooser
// port speaks in cards: first the distinct triggering cards, then — within one
// card — the distinct cards whose text granted it each ability. Either level with
// a single entry is forced and never prompts, so an event on one card with one
// ability is silent, as is a batch of creatures that each carry one.
func (g *Game) orderTriggered(actor int, pending []triggeredAbility) []triggeredAbility {
	if len(pending) <= 1 {
		return pending
	}
	ordered := make([]triggeredAbility, 0, len(pending))
	for _, src := range g.orderByChoice(actor, orderTriggerPrompt, distinctBy(pending, sourceOf)) {
		of := filterBy(pending, sourceOf, src)
		for _, grantor := range g.orderByChoice(actor, orderGrantorPrompt, distinctBy(of, grantorOf)) {
			ordered = append(ordered, filterBy(of, grantorOf, grantor)...)
		}
	}
	return ordered
}

// sourceOf and grantorOf name the two card fields orderTriggered orders by.
func sourceOf(t triggeredAbility) LocalID  { return t.source }
func grantorOf(t triggeredAbility) LocalID { return t.grantor }

// distinctBy lists the distinct values of key across pending, first-seen order.
func distinctBy(pending []triggeredAbility, key func(triggeredAbility) LocalID) []LocalID {
	var out []LocalID
	for _, t := range pending {
		if !slices.Contains(out, key(t)) {
			out = append(out, key(t))
		}
	}
	return out
}

// filterBy keeps the pending abilities whose key equals want, in collection order.
func filterBy(
	pending []triggeredAbility,
	key func(triggeredAbility) LocalID,
	want LocalID,
) []triggeredAbility {
	var out []triggeredAbility
	for _, t := range pending {
		if key(t) == want {
			out = append(out, t)
		}
	}
	return out
}

// triggeredAbility is an ability waiting to resolve in a trigger window. It
// retains the triggering card as source so a granted "purge this creature"
// resolves against the creature that gained it rather than the card that granted
// it, and the granting card as grantor for ordering and attribution.
type triggeredAbility struct {
	source  LocalID
	grantor LocalID
	ability Ability
}

// destroyedAbilities collects every Destroyed ability the creatures about to be
// destroyed carry: printed, upgrade-granted, and constant-granted. The collection
// happens before any resolves, then destroyTogether lets the active player order
// the whole set as KeyForge requires.
func (g *Game) destroyedAbilities(ids []LocalID) []triggeredAbility {
	var pending []triggeredAbility
	for _, id := range ids {
		pending = append(pending, g.triggeredBy(id, TriggerDestroyed)...)
	}
	return pending
}
