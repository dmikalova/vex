package engine

import (
	"strconv"
)

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
	return g.usableInActiveHouseWithoutPermit(id) ||
		g.offHouseUsePermit(g.controller(id), id) >= 0
}

// usableInActiveHouseWithoutPermit is usableInActiveHouse minus the this-turn
// off-house use permit (CXO Taber). It is the base house eligibility, so callers
// can tell whether a use spent a permit — the creature was usable only because of
// one — and charge it exactly once. A player at No House — locked out of every
// house and resolved to no active house for their play phase (ADR 0035) — matches
// nothing, so only Versatile cards and explicit off-house grants are usable.
// Before the choice resolves (HouseNone outside the play phase) nothing is gated.
func (g *Game) usableInActiveHouseWithoutPermit(id LocalID) bool {
	return g.manual ||
		(g.State.ActiveHouse == HouseNone && g.State.Phase != PhasePlay) ||
		g.House(id) == g.State.ActiveHouse ||
		g.hasKeyword(id, Versatile) ||
		(g.State.MayUseArtifactsAnyHouse[g.controller(id)] && g.TypeOf(id) == Artifact) ||
		(g.State.MayUseHouse[g.controller(id)] != HouseNone && g.House(id) == g.State.MayUseHouse[g.controller(id)]) ||
		(g.State.MayUseTrait[g.controller(id)] != traitUnset && g.HasTrait(id, g.State.MayUseTrait[g.controller(id)]))
}

// spendOffHouseUse charges one use of a this-turn off-house use permit when the
// creature was usable only because of it — a card in the active house, or freed by
// Versatile or a house grant, spends nothing (CXO Taber lets its controller use
// one non-Star Alliance card).
func (g *Game) spendOffHouseUse(player int, id LocalID) {
	if g.usableInActiveHouseWithoutPermit(id) {
		return
	}
	if i := g.offHouseUsePermit(player, id); i >= 0 {
		g.consumeOffHousePermit(player, i)
	}
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
	if g.atRuleOfSix(id) {
		return ErrRuleOfSix
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
	// An attached Upgrade may also bar the host's use (Earthbind).
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		if c := g.cat.def(up).Restricts.UseCondition; c != nil {
			ctx := &EffectContext{Resolver: g, Source: id, Controller: player}
			if !c.Met(ctx) {
				return ErrCannotUse
			}
		}
	}
	return nil
}

// canUse validates that a creature may be used to reap or fight right now.
func (g *Game) canUse(player int, id LocalID) error {
	if err := g.usable(player, id); err != nil {
		return err
	}
	if g.TypeOf(id) != Creature {
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
	if !g.cannotBeUsedTo(id, ReapUse) && !g.cannotReap(player) &&
		!g.cannotReapHouse(player, id) {
		return true
	}
	if g.canFight(player, id) {
		return true
	}
	return !g.cannotBeUsedTo(id, ActionUse) && g.hasTrigger(id, TriggerAction)
}

// canFight reports whether the creature is genuinely able to fight right now: it
// is not barred from fighting — player-wide (Fogbank, a Fighting restriction) or
// per-card — and has a legal target. It is the fight branch of hasAnyUse and what
// "must fight if able" (Little Rapscal) tests, so a creature that cannot fight is
// free to reap rather than stranded.
func (g *Game) canFight(player int, id LocalID) bool {
	return !g.cannotBeUsedTo(id, FightUse) &&
		!g.cannotFight(player) &&
		g.canFightSomeEnemy(player, id)
}

// canUseTo is canUse for one specific way of using the creature, so a card that
// bars only one of them (Tireless Crocag cannot reap) stays usable the other ways.
// A wrong-house error is forgiven for FightUse when a fight grant (Brothers in
// Battle) covers the creature, so the UI offers Fight on exactly the creatures
// Fight itself would let through.
func (g *Game) canUseTo(player int, id LocalID, kind UseKind) error {
	if err := g.canUse(player, id); err != nil &&
		(kind != FightUse || !g.fightErrorForgiven(err, id)) {
		return err
	}
	if kind == ReapUse && g.cannotReap(player) {
		return ErrCannotUse
	}
	if kind == ReapUse && g.cannotReapHouse(player, id) {
		return ErrCannotUse
	}
	if (kind == ReapUse || kind == ActionUse) && g.mustFightWhenUsed(player, id) {
		return ErrCannotUse
	}
	if g.cannotBeUsedTo(id, kind) {
		return ErrCannotUse
	}
	return nil
}

// mustFightWhenUsed reports whether this creature is currently barred from reaping
// or using an Action ability because it must fight instead — true when either the
// global "creatures must fight when used, if able" rule (Little Rapscal) is in play
// or the creature is enraged, and the creature is able to fight. A creature with
// nothing to fight, or barred from fighting (Fogbank), may still reap or act.
func (g *Game) mustFightWhenUsed(player int, id LocalID) bool {
	return (g.mustFightIfAble() || g.Enraged(id)) && g.canFight(player, id)
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
	if g.TypeOf(id) != Artifact {
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
	g.spendOffHouseUse(player, id)
	g.reapWith(id)
	// Boundary: the reap and any lasting reaction it fired can change power
	// anywhere on the board (ADR 0029).
	g.settleDestroyed(player)
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
	if g.TypeOf(id) != Creature || !g.State.Cards[id].Stunned {
		return ErrCannotUse
	}
	g.reapWith(id)
	// Boundary: the reap and any lasting reaction it fired can change power
	// anywhere on the board (ADR 0029).
	g.settleDestroyed(player)
	return nil
}

// recordUse marks one successful use of a creature this turn. "Used" is the
// rulebook umbrella for reaping, fighting, or using an Action: ability; the count
// advances before that use resolves so Fight:/Reap: abilities see their first use
// as count 1. Using any card also counts a usage of its name toward the Rule of
// Six, so an artifact's Action: is bounded even though it is not a creature.
func (g *Game) recordUse(id LocalID) {
	g.recordUsage(id)
	if g.TypeOf(id) == Creature {
		g.State.Cards[id].TimesUsedThisTurn++
		g.State.TurnHistory[g.controller(id)][CreaturesUsedThisTurn]++
	}
}

// reapWith performs a reap driven by a rule or ability, with no active-player or
// active-house checks: a stunned creature recovers instead; otherwise it
// exhausts, its controller gains 1 Æmber, and its "Reap:" abilities fire. A card
// that cannot reap is not made to by an ability either.
func (g *Game) reapWith(id LocalID) {
	if g.cannotBeUsedTo(id, ReapUse) || g.cannotReap(g.controller(id)) {
		return
	}
	g.recordUse(id)
	if g.recoverFromStun(id) {
		return
	}
	p := g.controller(id)
	g.State.Cards[id].Exhausted = true
	g.gainReapAember(p, id)
	// A reaction that asks "is this the first reap this turn?" (Aember Conduction
	// Unit) must count this reap, so the tally is raised before the window opens.
	g.State.TurnHistory[p][CreaturesReapedThisTurn]++
	g.emitReapWindow(p, id)
}

// emitReapWindow resolves everything that reacts to reaper's creature reaping as
// one ordered window: the creature's own "Reap:" and "after I am used" abilities,
// the "after you use a card" reactions on the user's other cards, and the
// board-wide "after a creature reaps" / "after an enemy creature reaps" reactions
// (Orb of Invidius, Pip Pip). They are gathered before any resolves so the active
// player orders the whole set (ADR 0013), collapsing identical abilities. The
// lasting-registry reactions (Crystal Hive) are folded into the same window, so a
// capable client interleaves them with the card abilities.
func (g *Game) emitReapWindow(reaper int, reaped LocalID) {
	pending := g.reapReactions(reaper, reaped)
	pending = append(pending, g.lastingReactions(EventReap, reaper, reaped)...)
	pending = append(pending, g.lastingReactions(EventUsed, reaper, reaped)...)
	g.resolveWindow(g.orderTriggered(reaper, pending))
}

// abilityWindow accumulates the triggered abilities that fire in one trigger
// window. The several gatherers (reap, action, creature-play, other-play) share it
// instead of each re-declaring the "collect the abilities, stamp their actor and
// it, append" closure.
type abilityWindow struct {
	g       *Game
	pending []triggeredAbility
}

// window starts an empty gather.
func (g *Game) window() *abilityWindow { return &abilityWindow{g: g} }

// add collects src's abilities under trigger, each resolving for src's own
// controller against it — the common case where a card's reaction resolves for
// whoever controls it.
func (w *abilityWindow) add(src LocalID, trigger Trigger, it LocalID, hasIt bool) {
	w.addAs(src, trigger, w.g.controller(src), it, hasIt)
}

// addAs is add with the resolving actor named explicitly, for an ability used on
// behalf of someone other than its controller (a borrowed Tactic's "Play:").
func (w *abilityWindow) addAs(src LocalID, trigger Trigger, actor int, it LocalID, hasIt bool) {
	got := w.g.triggeredBy(src, trigger)
	kept := got[:0]
	for i := range got {
		got[i].actor = int8(actor)
		got[i].it, got[i].hasIt = it, hasIt
		if w.g.firesForSubject(got[i]) {
			kept = append(kept, got[i])
		}
	}
	w.pending = append(w.pending, kept...)
}

// firesForSubject reports whether a reaction narrowed by a predicate fixed for the
// whole window fires for the card and turn in context. A reaction that renders
// "after you play an artifact" (a Conditional{ItIs} over its whole effect) does not
// fire on a creature, so it never joins the ordering window on a play it does not
// narrow to — Harmonia, narrowed to creatures, does not order with a played Tactic.
// Likewise "after an enemy creature is destroyed during your turn"
// (Conditional{And{ItIsEnemy, ItIsYourTurn}}) never joins the window on the
// opponent's turn. Only a fixed-context predicate gates firing this way; a volatile
// board "if" (Overwhelmed) is not one, so that reaction still fires and is checked
// at resolution.
func (g *Game) firesForSubject(t triggeredAbility) bool {
	if !t.hasIt {
		return true
	}
	cond, ok := fixedNarrowing(t.ability.Effect)
	if !ok {
		return true
	}
	ctx := &EffectContext{
		Resolver:   g,
		Source:     t.source,
		Controller: int(t.actor),
		It:         t.it,
		HasIt:      true,
	}
	return cond.Met(ctx)
}

// fixedNarrowing returns the condition a reaction uses to narrow itself out of the
// window before firing — a predicate fixed for the whole window, wrapping the whole
// effect with no Else. It is exactly the shape afterYouActOnText folds into "after
// you play an artifact", so a reaction that reads as narrowed is the reaction that
// fires only when narrowed.
func fixedNarrowing(e Effect) (Condition, bool) {
	c, ok := e.(Conditional)
	if !ok || c.Else != nil {
		return nil, false
	}
	if isFixedPredicate(c.Cond) {
		return c.Cond, true
	}
	return nil, false
}

// isFixedPredicate reports whether a condition is decided by facts fixed for the
// whole reaction window — the card in context (its house, type, name, trait,
// friendliness, or flank position) or whose turn it is — with no dependence on
// volatile board state, so a reaction gated on it can be narrowed before resolution
// instead of joining every ordering window. An And of fixed predicates is itself
// fixed (Dark Æmber Vault's friendly Mutant creature; Pile of Skulls' enemy
// creature during your turn); a board "if" (Overwhelmed) is not, because an earlier
// reaction in the window can change the count, so that reaction still fires and is
// rechecked at resolution.
func isFixedPredicate(c Condition) bool {
	switch cc := c.(type) {
	case ItIs, ItIsNamed, ItIsOfTrait, ItIsFriendly, ItIsEnemy, ItIsYourTurn:
		return true
	case OnFlank:
		return cc.OfIt
	case And:
		if len(cc.Conditions) == 0 {
			return false
		}
		for _, sub := range cc.Conditions {
			if !isFixedPredicate(sub) {
				return false
			}
		}
		return true
	}
	return false
}

// reapReactions gathers, in default resolution order, every card-sourced ability
// that reacts to reaped being used to reap: its own "Reap:" and "after I am used"
// abilities (bound to no "it"), the "after you use a card" reactions on the user's
// other in-play cards, and the board-wide "after a creature reaps" reactions —
// including those narrowed to an enemy reaper (Pip Pip) by an ItIsEnemy condition,
// which fire only when the reaper is the reacting card's enemy — the last bound to
// the reaping creature as "it". Every entry carries its own actor and "it" so the
// whole set can order as one window while each bystander reaction resolves for its
// owner.
func (g *Game) reapReactions(reaper int, reaped LocalID) []triggeredAbility {
	w := g.window()
	w.add(reaped, TriggerAfterReap, 0, false)
	w.add(reaped, TriggerAfterUsedSelf, 0, false)
	for _, id := range g.allInPlay(reaper) {
		if id != reaped {
			w.add(id, TriggerAfterUse, reaped, true)
		}
	}
	for player := 0; player < 2; player++ {
		for _, id := range g.allInPlay(player) {
			w.add(id, TriggerAfterCreatureReaps, reaped, true)
		}
	}
	return w.pending
}

// gainReapAember pays out the Æmber a reap grants to player p, applying any lasting
// replacement of that payout (Dimension Door makes it steal instead of gain).
func (g *Game) gainReapAember(p int, source LocalID) {
	if le, ok := g.lastingReplacement(p, EventReapAember); ok && le.Do == actSteal {
		stolen := min(1, g.State.Aember[1-p])
		g.SetAember(1-p, g.State.Aember[1-p]-stolen)
		g.SetAember(p, g.State.Aember[p]+stolen)
		g.record(ReapedStealing{
			Player: p,
			Card:   source,
			Amount: stolen,
			Cause:  le.Source,
		})
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
	if g.TypeOf(id) == Creature && g.mustFightWhenUsed(player, id) {
		return ErrCannotUse
	}
	if !g.hasTrigger(id, TriggerAction) {
		return ErrWrongType
	}
	if g.TypeOf(id) == Artifact {
		if err := g.chargeToll(player, TollUseArtifact); err != nil {
			return err
		}
	}
	g.spendOffHouseUse(player, id)
	g.useActionOf(player, id)
	// Boundary: the action and any lasting reaction it fired can change power
	// anywhere on the board (ADR 0029).
	g.settleDestroyed(player)
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
	pending := g.actionReactions(actor, id)
	pending = append(pending, g.lastingReactions(EventUsed, actor, id)...)
	g.resolveWindow(g.orderTriggered(actor, pending))
}

// actionReactions gathers the "Action:" ability and the reactions to using the
// card as one ordered window: the card's own "Action:" (resolving for actor, who
// may be borrowing the card via Remote Access), its "after I am used" (resolving
// for its own controller), and the "after you use a card" reactions on the acting
// player's other in-play cards, with the used card as "it".
func (g *Game) actionReactions(actor int, id LocalID) []triggeredAbility {
	w := g.window()
	w.addAs(id, TriggerAction, actor, 0, false)
	w.add(id, TriggerAfterUsedSelf, 0, false)
	for _, o := range g.allInPlay(actor) {
		if o != id {
			w.add(o, TriggerAfterUse, id, true)
		}
	}
	return w.pending
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
	if g.TypeOf(defender) != Creature ||
		g.controller(defender) == player ||
		!g.inPlay(defender) ||
		g.protectedByTaunt(attacker, defender) {
		return ErrNoTarget
	}
	// Camouflage bars an attacker that is not on a flank from fighting its host.
	if g.protectedFromNonFlank(defender) && !g.onFlankOf(attacker) {
		return ErrNoTarget
	}
	if fr := g.cat.def(attacker).FightRestriction; fr != (Target{}) &&
		!fr.allows(&EffectContext{Resolver: g, Source: attacker, Controller: player}, defender) {
		return ErrNoTarget
	}
	g.spendOffHouseUse(player, attacker)
	g.fight(attacker, defender)
	// Boundary: combat resolution and any lasting reaction it fired can change
	// power anywhere on the board (ADR 0029).
	g.settleDestroyed(player)
	return nil
}

// hasTrigger reports whether an in-play card has the trigger itself, from an
// attached upgrade, or from a constant ability affecting it. A blanked card's own
// printed trigger is ignored, so its "Action:"/"Omni:" is not offered while blank.
func (g *Game) hasTrigger(id LocalID, trigger Trigger) bool {
	if !g.textBlanked(id) && g.cat.def(id).hasTrigger(trigger) {
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
				if !g.constantActive(grantor, c) || !g.constantAffects(grantor, c, id) {
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
	var targets []LocalID
	for _, def := range g.State.Battleline[1-player].slice() {
		if g.fightAllows(player, attacker, def) {
			targets = append(targets, def)
		}
	}
	return targets
}

// fightAllows reports whether attacker may legally fight the enemy def right now,
// accounting for taunt and the attacker's fight restriction (Bigtwig fights only
// stunned creatures). It is the per-target predicate shared by FightTargets and
// canFightSomeEnemy; it does not repeat the player-wide and readiness checks its
// callers make.
func (g *Game) fightAllows(player int, attacker, def LocalID) bool {
	if g.protectedByTaunt(attacker, def) {
		return false
	}
	// Camouflage bars an attacker that is not on a flank from fighting its host.
	if g.protectedFromNonFlank(def) && !g.onFlankOf(attacker) {
		return false
	}
	fr := g.cat.def(attacker).FightRestriction
	return fr == (Target{}) ||
		fr.allows(&EffectContext{Resolver: g, Source: attacker, Controller: player}, def)
}

// canFightSomeEnemy reports whether attacker has at least one legal fight target,
// so hasAnyUse does not count fighting as an available use for a creature every
// enemy is protected from (taunt) or that its restriction forbids (Bigtwig with no
// stunned enemy). It does not call canUse, so hasAnyUse can call it without
// recursing.
func (g *Game) canFightSomeEnemy(player int, attacker LocalID) bool {
	for _, def := range g.State.Battleline[1-player].slice() {
		if g.fightAllows(player, attacker, def) {
			return true
		}
	}
	return false
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
	g.record(StunRecovered{Player: g.controller(id), Creature: id})
	// Spending a use to recover from a stun still counts as using the creature, so
	// the "after you use a creature" reactions fire (Legion's March). The other use
	// paths fold these into their own window; a stun recovery has no other window.
	user := g.controller(id)
	g.resolveWindow(g.orderTriggered(user, g.lastingReactions(EventUsed, user, id)))
	return true
}

// usableByAbility reports whether a card may be used by an ability right now. A
// card can only be used while ready, and only while its name is under the Rule of
// Six: an ability that would use a creature whose name has already been used six
// times this turn does nothing, the same cap the active player's own uses obey, so
// two Legatus Raptors readying and using each other cannot fight past six between
// them. Unlike canUse this ignores the active player and active house — ability-
// driven use cares only about readiness and the Rule of Six.
func (g *Game) usableByAbility(id LocalID) bool {
	if g.State.Cards[id].Exhausted || g.atRuleOfSix(id) {
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

// emitCreaturePlayed fires the "after a creature is played" reaction on every
// in-play card of both players except the played creature, with the played
// creature as "it" (The Big One). It fires only on an actual play from hand — the
// play path calls it — so a creature put into play by another effect does not, and
// it reaches the whole board, not only the played creature's neighbours.
func (g *Game) emitCreaturePlayed(played LocalID) {
	for player := 0; player < 2; player++ {
		for _, id := range g.allInPlay(player) {
			if id == played {
				continue
			}
			g.triggerAbilities(id, TriggerAfterCreaturePlayed, played, true)
		}
	}
}

// emitUpgradeEntered fires the "after an upgrade enters play" reaction on every
// in-play card of both players, with the entering upgrade as "it" (Armory Officer
// Nel). An upgrade is attached to its host, not itself in the battleline, so no
// card is skipped.
func (g *Game) emitUpgradeEntered(upgrade LocalID) {
	for player := 0; player < 2; player++ {
		for _, id := range g.allInPlay(player) {
			g.triggerAbilities(id, TriggerAfterUpgradeEnters, upgrade, true)
		}
	}
}

// afterDestroyedReactions gathers, as one ordered window, every reaction to the
// creatures of a finished Destroyed window now reaching their discard piles: the
// board-wide "after a creature is destroyed" reactions (Neffru) on every card —
// including those narrowed to an enemy death during your turn (Pile of Skulls, by
// And{ItIsEnemy, ItIsYourTurn}) or to a friendly death (Spartasaur, by ItIsFriendly)
// — and the lasting "each time an enemy creature is destroyed" reactions (Loot the
// Bodies) from the registry. Every entry carries its own actor and the destroyed
// creature as "it", so the active player orders the whole set while each reaction
// resolves for its owner (ADR 0013). It runs once the creatures have reached the
// discard pile, so a card destroyed in the same window is out of play and neither
// reacts nor can be chosen.
func (g *Game) afterDestroyedReactions(members []LocalID) []triggeredAbility {
	w := g.window()
	for _, id := range members {
		if g.TypeOf(id) != Creature {
			continue
		}
		for p := 0; p < 2; p++ {
			for _, c := range g.allInPlay(p) {
				w.add(c, TriggerAfterCreatureDestroyed, id, true)
			}
		}
		// Loot the Bodies' "each time an enemy creature is destroyed: gain Æmber" is a
		// registry reaction the enemy of the dead creature's controller owns; fold it
		// into the same window so it orders with the card reactions (ADR 0013).
		w.pending = append(
			w.pending,
			g.lastingReactions(EventEnemyCreatureDestroyed, 1-g.controller(id), id)...,
		)
	}
	return w.pending
}

// afterPlayReactions gathers, as one ordered window, the "Play:" ability of a
// just-played non-creature card and the "after you play a card" reactions on the
// board. The card's own "Play:" resolves for player, who may be borrowing it (a
// Tactic played from the opponent's discard via Mimicry); the after-card-played
// reactions resolve for their own controllers, with the played card as "it". The
// played card is one of its own targets for a friendly "after you play a card"
// reaction: an artifact's abilities are live the moment it enters play, so it
// counts as a card you played (Harmonia, Hunting Witch; rulebook Harmonia 357 FAQ).
// The caller appends the EventCardPlayed lasting reactions (Library Access) to this
// window, so they order together with the card abilities rather than in a trailing
// window of their own (ADR 0013).
func (g *Game) afterPlayReactions(player int, played LocalID) []triggeredAbility {
	w := g.window()
	w.addAs(played, TriggerAfterPlay, player, 0, false)
	for _, id := range g.allInPlay(player) {
		w.add(id, TriggerAfterCardPlayed, played, true)
	}
	for _, id := range g.allInPlay(1 - player) {
		w.add(id, TriggerAfterEnemyCardPlayed, played, true)
	}
	return w.pending
}

// playCreatureReactions gathers every ability that triggers when a creature is
// played into one ordered window: the creature's own "Play:" and "enters play",
// every other card's "after a creature enters play" and "after a creature is
// played", each neighbour's "after a creature is played adjacent", and the
// playing player's "after you play a card" (with the opponent's "after an enemy
// plays a card"). The played creature is "it" for the bystander reactions. Every
// entry resolves for its own controller — a creature is always played by its
// controller, so there is no borrower asymmetry here.
func (g *Game) playCreatureReactions(player int, played LocalID) []triggeredAbility {
	w := g.window()
	w.add(played, TriggerAfterPlay, 0, false)
	w.add(played, TriggerEntersPlay, 0, false)
	for p := 0; p < 2; p++ {
		for _, id := range g.allInPlay(p) {
			if id != played {
				w.add(id, TriggerAfterCreatureEnters, played, true)
			}
		}
	}
	for _, neighbor := range neighbors(&EffectContext{Resolver: g}, played) {
		w.add(neighbor, TriggerAfterCreaturePlayedAdjacent, played, true)
	}
	for p := 0; p < 2; p++ {
		for _, id := range g.allInPlay(p) {
			if id != played {
				w.add(id, TriggerAfterCreaturePlayed, played, true)
			}
		}
	}
	for _, id := range g.allInPlay(player) {
		w.add(id, TriggerAfterCardPlayed, played, true)
	}
	for _, id := range g.allInPlay(1 - player) {
		w.add(id, TriggerAfterEnemyCardPlayed, played, true)
	}
	return w.pending
}

// emitActionPlayedBeforeResolve fires "after a Tactic is played but before it
// resolves" abilities on every in-play card of both players, with the played
// Tactic as "it". It runs before the Tactic's own Play: effect resolves, so a
// reaction (Encounter Suit warding its host) acts on the board the Tactic is about
// to affect. A Tactic either player plays fires it.
func (g *Game) emitActionPlayedBeforeResolve(player int, played LocalID) {
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerAfterTacticPlayedBeforeResolve, played, true)
	}
	for _, id := range g.allInPlay(1 - player) {
		g.triggerAbilities(id, TriggerAfterTacticPlayedBeforeResolve, played, true)
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
	g.triggerAbilitiesRooted(actor, src, trigger, it, hasIt, 0, false)
}

// triggerAbilitiesRooted is triggerAbilitiesAs carrying the Rule-of-Six cascade
// root, so a chain of Replicator-style triggers charges the card that started it
// rather than each card whose ability resolves along the way.
func (g *Game) triggerAbilitiesRooted(
	actor int,
	src LocalID,
	trigger Trigger,
	it LocalID,
	hasIt bool,
	root LocalID,
	hasRoot bool,
) {
	pending := g.triggeredBy(src, trigger)
	for i := range pending {
		pending[i].actor = int8(actor)
		pending[i].it, pending[i].hasIt = it, hasIt
		pending[i].root, pending[i].hasRoot = root, hasRoot
	}
	g.resolveWindow(g.orderTriggered(actor, pending))
}

// resolveWindow resolves an already-ordered window of triggered abilities in
// order, each for its own actor and against its own source, grantor, and "it"
// binding. It is the single resolve loop every trigger window flows through: a
// lone trigger from triggerAbilitiesAs, and the unified use/play windows that
// gather one card's own printed ability together with the bystander "after ..."
// reactions so the active player orders the whole set as one (ADR 0013), while each
// bystander reaction still resolves for its own controller. A window may also carry
// duration reactions from the lasting registry, which resolve through
// resolveReaction rather than as card abilities.
func (g *Game) resolveWindow(ordered []triggeredAbility) {
	for _, t := range ordered {
		if t.lasting {
			// A duration reaction has no source card and never checks in-play: its
			// subject may have left mid-window (a played creature destroyed by an
			// earlier reaction) yet the controller's economy reaction (Full Moon's
			// Æmber) still fires. It opens no card Frame — there is no card to
			// attribute it to — and settles after like every other entry (ADR 0029).
			g.resolveReaction(t.le, int(t.actor), t.it)
			if t.le.Once {
				g.removeLasting(t.le)
			}
			g.settleDestroyed(int(t.actor))
			continue
		}
		if g.resolveTriggered(t) {
			// Boundary: a resolved ability can change power anywhere on the board (a
			// buff left, a counter moved, Æmber spent off a creature that draws power
			// from it), so settle here rather than making each effect settle itself.
			// New mechanics rely on this boundary and must not hand-call
			// settleDestroyed (ADR 0029). The settling flag collapses nested sweeps,
			// so a batch still leaves together.
			g.settleDestroyed(int(t.actor))
		}
	}
}

// resolveTriggered resolves one triggered ability against its own source, grantor,
// and "it" binding, inside its own Frame, and reports whether it resolved. It is
// the single per-entry step every trigger window shares — the up-front-ordered
// windows (resolveWindow) and the Destroyed window that re-gathers as it resolves
// (resolveDestroyedWindow) — so ordering is the only thing that differs between
// them, and the source guard and Frame are identical everywhere.
//
// The source-in-play guard is applied here (RAW §190, ADR 0030): an earlier ability
// in the same window can take src out of play (Strange Gizmo destroys friendly
// artifacts as its house is chosen; a forge-key reactor destroys itself), and a
// card that has left play resolves nothing more, so it is skipped and false
// returned. A tactic is the one source that resolves its own ability while not in
// play — it never enters play, so "source in play" does not apply to it; its Play:
// still resolves. A TriggersFromDiscard card is the other exception: it resolves
// its choose-house ability from its owner's discard pile (Relentless Creeper
// returns itself to hand), so a source that is discard-active is not skipped.
func (g *Game) resolveTriggered(t triggeredAbility) bool {
	src := t.source
	actor := int(t.actor)
	if !g.inPlay(src) && g.cat.def(src).Type != Tactic && !g.activeInDiscard(src) {
		return false
	}
	closeFrame := g.openFrame(Frame{
		Actor:      actor,
		Source:     src,
		HasSource:  true,
		Trigger:    t.ability.Trigger,
		Grantor:    t.grantor,
		HasGrantor: t.grantor != src,
	})
	ec := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: actor,
		It:         t.it,
		HasIt:      t.hasIt,
		Root:       t.root,
		HasRoot:    t.hasRoot,
	}
	// A granted ability still knows the card that granted it (Source is the
	// creature the ability now lives on), so an effect on the granted ability can
	// move or anchor to that grantor — a "blaster" upgrade attaching itself to its
	// named creature, or Uncharted Lands' reap moving Æmber off that one artifact.
	if t.grantor != src {
		ec.Grantor = t.grantor
		ec.HasGrantor = true
		if g.cat.def(t.grantor).Type == Upgrade {
			ec.Upgrade = t.grantor
		}
	}
	t.ability.Effect.Resolve(ec)
	closeFrame()
	return true
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
		if g.textBlanked(src) {
			break
		}
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
				if len(c.Granted) == 0 || !g.constantActive(grantor, c) ||
					!g.constantAffects(grantor, c, src) {
					continue
				}
				for _, ab := range c.Granted {
					keep(grantor, ab)
				}
			}
		}
	}
	// An also-triggers-on rule (Kompsos Haruspex's constant, Livia the Elder's lasting fuse)
	// makes src's own abilities under one trigger also fire on this one — its play
	// effect on reap, its fight and reap effects on each other. Gather those printed
	// abilities here so an also-fired ability orders in the same window as a natural one.
	if !g.textBlanked(src) {
		for _, from := range g.additionalTriggers(src, trigger) {
			for _, ab := range g.cat.def(src).Abilities {
				if ab.Trigger == from {
					pending = append(pending, triggeredAbility{
						source:  src,
						grantor: src,
						ability: ab,
					})
				}
			}
		}
	}
	// A creature that has GAINED another card's text box (Mimic Gel copies a chosen
	// creature; Creed of Nurture lends one for the turn) also fires that card's
	// printed abilities as if they were its own. They keep src as their source, so
	// self-referential text rebinds to the gaining creature.
	if !g.textBlanked(src) {
		for _, textSource := range g.grantedTextBoxSources(src) {
			for _, ab := range g.cat.def(textSource).Abilities {
				keep(src, ab)
			}
		}
	}
	return pending
}

// orderTriggerPrompt is the prompt shown when abilities on several cards trigger
// at once and the active player must say which card's resolves next.
const orderTriggerPrompt = "Choose which card's ability resolves next"

// orderDestroyedPrompt names the destroyed-ability window plainly, since the
// player is arranging the Destroyed abilities of several cards leaving play at
// once rather than picking one card off a generic list. The Destroyed window
// orders itself through pickNextReaction as it re-gathers (a Destroyed ability can
// destroy more creatures, folding their abilities into the same window), so it
// cannot pre-order the whole window up front the way orderTriggered does — but it
// resolves each picked entry through the same resolveTriggered step every other
// window uses, so the source guard and Frame stay identical.
const orderDestroyedPrompt = "Resolve destroyed abilities"

// orderTriggered lets the active player arrange a trigger window's abilities into
// a resolution order (ADR 0013). The window is one flat list: every pending
// ability — a card's printed or granted ability and a duration reaction from the
// lasting registry alike — is one entry the player orders through the flat
// ReactionChooser port, picking the next to resolve until one remains. Identical
// abilities (same trigger and same rendered text, even from different cards)
// resolve the same in any order, so the player is prompted only while the pending
// entries are not all identical: a window whose entries are all identical is
// auto-ordered and never prompts, and once a series of picks leaves only identical
// entries the remainder auto-resolves too. (Identity is compared, not card, because
// the same card can resolve differently as the board changes.) A Chooser without
// the ReactionChooser port keeps the gathered order, so the AI and the simulator
// resolve a window in the order it was collected.
func (g *Game) orderTriggered(
	actor int,
	pending []triggeredAbility,
) []triggeredAbility {
	if len(pending) <= 1 || g.allIdentical(pending) {
		return pending
	}
	return g.orderReactionsFlat(actor, orderTriggerPrompt, pending)
}

// allIdentical reports whether every pending ability shares one identity, so the
// window resolves the same in any order and never needs a prompt.
func (g *Game) allIdentical(pending []triggeredAbility) bool {
	first := identityOf(pending[0])
	for _, t := range pending[1:] {
		if identityOf(t) != first {
			return false
		}
	}
	return true
}

// orderReactionsFlat asks the active player to arrange the whole window by picking
// the next reaction to resolve, dropping it, and repeating until one remains (the
// last is forced). Card abilities and duration reactions are offered together in
// one labeled list. A Chooser without the ReactionChooser port, or one that
// returns an out-of-range index, keeps the reactions in their gathered order.
func (g *Game) orderReactionsFlat(
	actor int,
	prompt string,
	pending []triggeredAbility,
) []triggeredAbility {
	remaining := append([]triggeredAbility(nil), pending...)
	ordered := make([]triggeredAbility, 0, len(pending))
	for len(remaining) > 1 {
		i := g.pickNextReaction(actor, prompt, remaining)
		ordered = append(ordered, remaining[i])
		remaining = append(remaining[:i], remaining[i+1:]...)
	}
	return append(ordered, remaining...)
}

// pickNextReaction returns the index of the next reaction to resolve from pending.
// It answers without a prompt — index 0, the gathered front — when one entry
// remains, every entry is identical (they resolve the same in any order), the
// Chooser lacks the ReactionChooser port, or the Chooser returns an out-of-range
// index. Otherwise the active player picks from the flat labeled list. It is the
// one-step primitive behind both a whole window ordered up front (orderReactionsFlat)
// and the Destroyed window that re-gathers as it resolves (resolveDestroyedWindow).
func (g *Game) pickNextReaction(actor int, prompt string, pending []triggeredAbility) int {
	if len(pending) <= 1 || g.allIdentical(pending) {
		return 0
	}
	rc, ok := g.chooserFor(actor).(ReactionChooser)
	if !ok {
		return 0
	}
	reactions := make([]OrderableReaction, len(pending))
	for i, t := range pending {
		reactions[i] = g.orderableFor(t)
	}
	i := rc.ChooseReaction(prompt, reactions)
	if i < 0 || i >= len(pending) {
		return 0
	}
	return i
}

// orderableFor renders a pending ability for the flat reaction list: a card
// ability shows its source card and rendered ability text, a duration reaction its
// rendered effect. HasCard/Card let a client that highlights the board point at
// the source card as well.
func (g *Game) orderableFor(r triggeredAbility) OrderableReaction {
	if r.lasting {
		return OrderableReaction{Label: r.le.Do.describe()}
	}
	return OrderableReaction{
		Card:    r.source,
		HasCard: true,
		Label:   g.cat.def(r.source).Name + ": " + r.ability.Effect.Text(),
	}
}

// abilityIdentity keys an ability by what it will do — its trigger and rendered
// text — so two abilities that resolve identically share a key. It is what
// allIdentical compares, so a window only auto-orders when every entry would
// resolve the same.
func abilityIdentity(a Ability) string {
	return strconv.Itoa(int(a.Trigger)) + "\x00" + a.Effect.Text()
}

// identityOf keys any pending ability — a card ability by abilityIdentity, a
// duration reaction by its event, amount, and rendered effect. The duration key is
// prefixed so it can never collide with a card ability's, keeping the two kinds
// from ever being treated as identical.
func identityOf(t triggeredAbility) string {
	if t.lasting {
		return "L\x00" + strconv.Itoa(int(t.le.On)) + "\x00" +
			strconv.Itoa(int(t.le.Amount)) + "\x00" + t.le.Do.describe()
	}
	return abilityIdentity(t.ability)
}

// triggeredAbility is an ability waiting to resolve in a trigger window. It
// retains the triggering card as source so a granted "purge this creature"
// resolves against the creature that gained it rather than the card that granted
// it, and the granting card as grantor for ordering and attribution. it/hasIt is
// the "it" the ability resolves against — unset for a card's own "after I ..."
// ability, and the acting card for a bystander's "after a creature ..." reaction,
// so a unified window can carry both self and bystander reactions at once. actor
// is the player the ability resolves for: its own controller for a natural
// trigger, or someone acting on its behalf (Remote Access). A window is ordered by
// the active player but each entry still resolves for its own actor, so a
// bystander's "after an enemy creature reaps: gain Æmber" credits its owner.
type triggeredAbility struct {
	source  LocalID
	grantor LocalID
	ability Ability
	actor   int8
	it      LocalID
	hasIt   bool
	// root carries the Rule-of-Six cascade root through a chain of Replicator-style
	// triggers, so every ability the chain resolves charges the card that started
	// it rather than each resolving card's own name. hasRoot reports whether one is
	// set; only a chained trigger window carries it.
	root    LocalID
	hasRoot bool
	// lasting marks a duration reaction from the lasting registry (Full Moon,
	// Charge!, Crystal Hive) rather than a card's printed ability. It has no source
	// card in play, so it carries its LastingEffect in le and resolves through
	// resolveReaction; its subject is it and its controller is actor. Carrying it in
	// the same slice lets a window order duration reactions together with the card
	// abilities that fire on the same event (ADR 0013).
	lasting bool
	le      LastingEffect
}

// destroyedAbilities collects every Destroyed ability the creatures about to be
// destroyed carry: printed, upgrade-granted, and constant-granted. Each ability
// resolves for its creature's controller, so the actor is set at gather time (the
// creatures are all still in play). The collection happens before any resolves,
// then destroyTogether lets the active player order the whole set as KeyForge
// requires.
func (g *Game) destroyedAbilities(ids []LocalID) []triggeredAbility {
	if g.triggerDisabled(TriggerDestroyed) {
		return nil
	}
	var pending []triggeredAbility
	for _, id := range ids {
		got := g.triggeredBy(id, TriggerDestroyed)
		for i := range got {
			got[i].actor = int8(g.controller(id))
		}
		pending = append(pending, got...)
	}
	return pending
}

// triggerDisabled reports whether any active constant ability in play disables the
// given trigger, so no card's ability with that trigger fires — Purifier of Souls
// disables every Destroyed ability while it stays in play. It reads board-wide,
// independent of the disabling ability's Target.
func (g *Game) triggerDisabled(t Trigger) bool {
	for p := 0; p < 2; p++ {
		for _, src := range g.allInPlay(p) {
			for _, c := range g.cat.def(src).ConstantAbilities {
				if len(c.DisableTriggers) == 0 || !g.constantActive(src, c) {
					continue
				}
				for _, dt := range c.DisableTriggers {
					if dt == t {
						return true
					}
				}
			}
		}
	}
	return false
}
