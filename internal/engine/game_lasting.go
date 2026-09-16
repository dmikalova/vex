package engine

// This file holds LASTING EFFECTS: "for the remainder of the turn" effects that
// attach to a later game event. Three flavors share the same flat registry:
//
//   - a REACTION runs after an event (Full Moon gains Æmber after you play a
//     creature, Charge! deals damage after you play a creature, Crystal Hive gains
//     Æmber after a creature reaps). Reactions fold into the event's trigger window
//     through lastingReactions, so they order together with the card abilities that
//     fire on the same event when several fire at once (ADR 0013).
//   - a REPLACEMENT changes an event's own outcome before it happens (Dimension
//     Door makes reaping steal Æmber instead of gaining it). The event site queries
//     for a replacement (lastingReplacement) and applies it in place.
//   - a MODIFIER adjusts an event's amount rather than swapping its outcome (Lethal
//     Distraction makes a chosen creature take additional damage whenever it takes
//     damage). Like a replacement it is queried at the event site, but the matches
//     are summed (lastingExtraDamage) so modifiers stack, and it is subject-scoped
//     to the one creature it names.
//
// Because the game state is a flat, pointerless value it cannot hold effect
// closures, so a lasting effect is a small enum-tagged record (LastingEffect)
// resolved by a central switch. Adding a new one is a new Event and/or
// lastingAction plus a case here — never a bespoke branch in the play/reap path.

// Event names a moment a lasting effect attaches to.
type Event uint8

const (
	// eventUnset is the zero value. A real reaction or replacement always names one
	// of the points below, so an unset Event marks "no event" — used to tell a
	// StaticModifier that carries no replacement from one that does.
	eventUnset Event = iota
	// EventCreaturePlayed fires after the controller plays a creature (a reaction
	// point). Full Moon and Charge! attach here.
	EventCreaturePlayed
	// EventReap fires after one of the controller's creatures reaps (a reaction
	// point). Crystal Hive attaches here.
	EventReap
	// EventFight fires after one of the controller's creatures fights (a reaction
	// point). Warsong attaches here.
	EventFight
	// EventEnemyCreatureDestroyed fires, for a player, each time a creature their
	// opponent controls is destroyed (a reaction point). Loot the Bodies attaches
	// here; it is dispatched to the opponent of the destroyed creature's controller.
	EventEnemyCreatureDestroyed
	// EventReapAember is the Æmber a reap grants (a replacement point). Dimension
	// Door replaces it.
	EventReapAember
	// EventCreatureDestroyed is a creature about to be destroyed (a replacement
	// point). An attached Upgrade that replaces its host's destruction — Armageddon
	// Cloak — is queried here.
	EventCreatureDestroyed
	// EventAemberAddedToPool is Æmber about to be added to a player's pool from the
	// common supply (a replacement point). A creature that captures such Æmber before
	// it lands — Ether Spider, scoped to its opponent's pool — is queried here. It
	// intercepts the incoming gain, so the pool's owner keeps the Æmber they already
	// have.
	EventAemberAddedToPool
	// EventAemberTakenFromPool is Æmber about to be taken from a player's pool by a
	// steal or capture (a replacement point). A card that redirects the take's source
	// — Po's Pixies, scoped to its controller's own pool — is queried here, drawing
	// the Æmber from the common supply instead so the pool's owner keeps what they
	// have while the taker still gains it. It is the source end of the same Æmber-flow
	// replacement spine EventAemberAddedToPool is the destination end of.
	EventAemberTakenFromPool
	// EventAemberStolen is Æmber a steal has just removed from a pool, about to be
	// added to the thief's pool (a replacement point). A card that redirects the
	// theft's destination — Gargantodon, capturing it onto a creature the active
	// player controls — is queried here. Unlike the pool events it is global rather
	// than pool-scoped: the redirect applies to every steal regardless of who
	// controls the card.
	EventAemberStolen
	// EventForgeKey fires after a player forges a key (a reaction point). A reaction
	// owned by that player fires during their turn; because the registry clears a
	// player's own entries at their ready phase, a reaction armed on an opponent
	// (owned by the forger) survives the arming turn and fires during the forger's
	// next turn — the "during your opponent's next turn" window (Interdimensional
	// Graft).
	EventForgeKey
	// EventCardEntersPlay fires after the controller plays a card that stays in play
	// under its own name — a creature or an artifact (a reaction point). It is the
	// window "the next creature/artifact you play" effects arm, so they can reach an
	// artifact too, which EventCreaturePlayed cannot.
	EventCardEntersPlay
	// EventCardPlayed fires after the controller plays a card of any type, including
	// one that never enters play (a reaction point). Library Access attaches here,
	// excepting itself so it draws only for another card.
	EventCardPlayed
	// EventCreatureTakesDamage is an amount of damage about to land on a creature (a
	// modifier point, always subject-scoped). Lethal Distraction attaches a "this
	// creature takes an additional N damage whenever it takes damage" modifier to one
	// chosen creature; applyRawDamage queries the registry and sums the bonuses.
	EventCreatureTakesDamage
	// EventBeforeFight fires just before a creature is used to fight, keyed to the
	// attacker and always subject-scoped. Unlike the other reaction points it is fired
	// by fireLastingBeforeFight, which matches on the subject alone rather than the
	// acting player, so a grant that lasts into the opponent's turn (Diplomacy's "each
	// creature gains 'Before Fight: Exalt this creature'") fires for an enemy attacker
	// on the opponent's turn too.
	EventBeforeFight // EventBonusIconBoost is a one-shot arming, not a reaction or replacement: it is
	// neither fired in a window nor queried at an outcome, only consumed by the play
	// path the next time its owner plays a card (Wild Bounty). The play-path icon loop
	// consumes it via consumeBonusIconBoost and resolves each icon an additional time.
	EventBonusIconBoost
)

// isReaction reports whether the event is a reaction point (fired after) rather
// than a replacement point (queried during).
func (e Event) isReaction() bool {
	return e == EventCreaturePlayed || e == EventReap || e == EventFight ||
		e == EventEnemyCreatureDestroyed || e == EventForgeKey ||
		e == EventCardEntersPlay || e == EventCardPlayed
}

// clause renders the "when" phrase for a reaction, e.g. "each time you play a
// creature".
func (e Event) clause() string {
	switch e {
	case EventReap:
		return "after a creature reaps"
	case EventFight:
		return "each time a friendly creature fights"
	case EventEnemyCreatureDestroyed:
		return "each time an enemy creature is destroyed"
	case EventForgeKey:
		return "after forging a key"
	case EventCardPlayed:
		return "each time you play another card"
	default:
		return "each time you play a creature"
	}
}

// gerund renders the "instead of ..." phrase for a replacement, e.g. "gaining
// Æmber from reaping".
func (e Event) gerund() string {
	return "gaining Æmber from reaping"
}

// clauseOnOpponentTurn renders the reaction's "when" phrase from the opponent's
// turn, where the acting creature is an enemy one. It differs from clause only for
// events whose phrasing names the acting side; the rest read the same either way.
func (e Event) clauseOnOpponentTurn() string {
	switch e {
	case EventFight:
		return "after an enemy creature is used to fight"
	default:
		return e.clause()
	}
}

// lastingAction is what a lasting effect does when it fires or replaces.
type lastingAction uint8

const (
	actGainAember lastingAction = iota
	actDealDamage
	actSteal
	actReadyPlayed
	actGiveRemainingAember
	actCapture
	actDraw
	actLoseAember
	// actExalt places Amount Æmber on the subject creature — the granted "Before
	// Fight: Exalt this creature" (Diplomacy).
	actExalt
	// actStun stuns the subject creature — Foggify's "after an enemy creature is
	// used to fight, stun it", armed on the opponent for their next turn.
	actStun
	// actTakeExtraDamage is a modifier, not a reaction: it never resolves in a
	// trigger window, only summed at the damage site by lastingExtraDamage.
	actTakeExtraDamage
	// actResolveBonusAgain is a one-shot boost, not a reaction: it never resolves in
	// a window, only consumed by the play path via consumeBonusIconBoost, which
	// resolves each icon of the next played card an additional time (Wild Bounty).
	actResolveBonusAgain
)

// describe is a short label for a reaction, used when the controller orders several
// that fire at once.
func (a lastingAction) describe() string {
	switch a {
	case actDealDamage:
		return "deal damage"
	case actReadyPlayed:
		return "ready the creature"
	case actCapture:
		return "capture Æmber"
	case actGiveRemainingAember:
		return "give remaining Æmber"
	case actDraw:
		return "draw a card"
	case actLoseAember:
		return "opponent loses Æmber"
	case actSteal:
		return "steal Æmber"
	case actExalt:
		return "exalt the creature"
	case actStun:
		return "stun the creature"
	default:
		return "gain Æmber"
	}
}

// LastingEffect is one registered lasting effect: the event it attaches to, what it
// does, the player it belongs to, and its magnitude. It is a plain value so the
// flat GameState can hold a fixed array of them.
type LastingEffect struct {
	// The event this attaches to, what it does, whose it is, and its magnitude.
	On         Event
	Do         lastingAction
	Controller int8
	Amount     int8
	// House, when set, limits a reaction to a subject of that house (Blypyp readies
	// only Mars creatures). HouseNone (the zero value) reacts to any subject.
	House House
	// Type, when set, limits a reaction to a subject of that card type (Soft Landing
	// readies the next creature or artifact, not the next upgrade). TypeUnset (the
	// zero value) reacts to any type, and AnyType means "creature or artifact" — the
	// two types that stay in play under their own name.
	Type CardType
	// Once removes the record after it fires a single time — "the next" rather than
	// "each time".
	Once bool
	// Subject, when HasSubject is set, limits a reaction to that one creature as the
	// subject of the event — Spectral Tunneler grants a single chosen creature
	// "Reap: Draw a card", so only that creature's reap draws. A zero HasSubject
	// reacts to any subject.
	Subject    LocalID
	HasSubject bool
	// Except, when HasExcept is set, is the card that armed the reaction, which it
	// does not fire for: Library Access draws "each time you play another card", and
	// playing Library Access itself is not another card.
	Except    LocalID
	HasExcept bool
	// Source, when HasSource is set, is the card whose ability installed the
	// reaction, so its payout attributes to that card ("Full Moon has Player 1 gain
	// 2 Æmber") rather than to the player, who cannot produce the outcome alone. A
	// zero HasSource — a test-built effect — falls back to the player.
	Source    LocalID
	HasSource bool
}

// maxLasting bounds how many lasting effects can be active at once — generous for
// the handful of "remainder of the turn" cards, and a fixed size keeps the state a
// flat value.
const maxLasting = 8

// AddLasting registers a lasting effect owned by le.Controller for the rest of
// their turn, dropping it silently when the registry is full. It is the single
// seam a "for the remainder of the turn" effect uses instead of hardcoding itself
// into the play or reap path; the record's own fields (Once, House, Type, Except)
// narrow when it fires.
func (g *Game) AddLasting(le LastingEffect) {
	if int(g.State.LastingCount) >= maxLasting {
		return
	}
	g.State.Lasting[g.State.LastingCount] = le
	g.State.LastingCount++
}

// clearLasting drops the lasting effects a player owns, called when their turn ends
// so the effects last only that turn.
func (g *Game) clearLasting(player int) {
	n := 0
	for i := 0; i < int(g.State.LastingCount); i++ {
		if int(g.State.Lasting[i].Controller) != player {
			g.State.Lasting[n] = g.State.Lasting[i]
			n++
		}
	}
	for i := n; i < int(g.State.LastingCount); i++ {
		g.State.Lasting[i] = LastingEffect{}
	}
	g.State.LastingCount = uint8(n)
	g.clearLastingAlsoTriggers(player)
}

// matchingLasting collects every registry reaction actor owns that responds to
// event for subject, applying the house, type, subject, and except filters. It is
// the shared gather behind lastingReactions, which turns the matches into window
// entries that order alongside the card abilities firing on the same event.
func (g *Game) matchingLasting(event Event, actor int, subject LocalID) []LastingEffect {
	var pending []LastingEffect
	for i := 0; i < int(g.State.LastingCount); i++ {
		le := g.State.Lasting[i]
		if int(le.Controller) != actor || le.On != event {
			continue
		}
		if le.House != HouseNone && g.cat.def(subject).House != le.House {
			continue
		}
		if le.Type != TypeUnset && !le.Type.reacts(g.cat.def(subject).Type) {
			continue
		}
		if le.HasSubject && le.Subject != subject {
			continue
		}
		if le.HasExcept && le.Except == subject {
			continue
		}
		pending = append(pending, le)
	}
	return pending
}

// lastingReactions gathers the registry reactions responding to event as window
// entries, so they order together with the card abilities that fire on the same
// event rather than in a trailing window of their own (ADR 0013). Each carries its
// LastingEffect and resolves through resolveWindow's duration branch, crediting
// actor with subject as its subject.
func (g *Game) lastingReactions(event Event, actor int, subject LocalID) []triggeredAbility {
	matched := g.matchingLasting(event, actor, subject)
	pending := make([]triggeredAbility, len(matched))
	for i, le := range matched {
		pending[i] = triggeredAbility{
			actor:   int8(actor),
			it:      subject,
			hasIt:   true,
			lasting: true,
			le:      le,
		}
	}
	return pending
}

// fireLastingBeforeFight resolves every lasting "Before Fight" ability granted to
// the attacker, keyed to the attacker as subject and fired regardless of whose turn
// it is. Diplomacy grants each creature "Before Fight: Exalt this creature" until
// the granting player's next turn, so an enemy creature exalts when it fights on the
// opponent's turn too — which is why this matches on the subject alone rather than
// on the acting player the way the event windows gather their reactions.
func (g *Game) fireLastingBeforeFight(attacker LocalID) {
	actor := g.controller(attacker)
	for i := 0; i < int(g.State.LastingCount); i++ {
		le := g.State.Lasting[i]
		if le.On != EventBeforeFight || !le.HasSubject || le.Subject != attacker {
			continue
		}
		g.resolveReaction(le, actor, attacker)
	}
}

// removeLasting drops the first registry entry equal to target — a one-shot
// reaction removing itself after it fires.
func (g *Game) removeLasting(target LastingEffect) {
	for i := 0; i < int(g.State.LastingCount); i++ {
		if g.State.Lasting[i] != target {
			continue
		}
		for j := i; j < int(g.State.LastingCount)-1; j++ {
			g.State.Lasting[j] = g.State.Lasting[j+1]
		}
		g.State.LastingCount--
		g.State.Lasting[g.State.LastingCount] = LastingEffect{}
		return
	}
}

// resolveReaction resolves a single reaction for actor, using subject as the
// triggering card where one is needed.
func (g *Game) resolveReaction(le LastingEffect, actor int, subject LocalID) {
	if le.HasSource {
		defer g.openFrame(Frame{Actor: actor, Source: le.Source, HasSource: true})()
	}
	switch le.Do {
	case actDealDamage:
		DealDamage{
			Amount: int(le.Amount),
			Target: Target{Kind: TargetChosenEnemyCreature},
		}.Resolve(
			&EffectContext{
				Resolver:   g,
				Source:     subject,
				Controller: actor,
			},
		)
	case actReadyPlayed:
		g.State.Cards[subject].Exhausted = false
		g.record(CreatureReadied{Creature: subject})
	case actCapture:
		CaptureAember{
			Amount: int(le.Amount),
			Target: Target{Kind: TargetTriggeringCreature},
			Source: Opponent,
		}.Resolve(
			&EffectContext{
				Resolver:   g,
				Source:     subject,
				Controller: actor,
				It:         subject,
				HasIt:      true,
			},
		)
	case actExalt:
		g.addAmberOn(subject, int(le.Amount))
		g.record(AemberExalted{Creature: subject, Amount: int(le.Amount)})
	case actStun:
		if g.inPlay(subject) {
			Stun{Target: Target{Kind: TargetTriggeringCreature}}.Resolve(
				&EffectContext{
					Resolver:   g,
					Source:     subject,
					Controller: actor,
					It:         subject,
					HasIt:      true,
				},
			)
		}
	case actDraw:
		g.draw(actor, int(le.Amount))
		g.record(LastingDraw{Player: actor, Amount: int(le.Amount), On: le.On})
	case actLoseAember:
		LoseAember{Player: Opponent, Amount: int(le.Amount)}.Resolve(
			&EffectContext{
				Resolver:   g,
				Source:     subject,
				Controller: actor,
			},
		)
	case actSteal:
		StealAember{Amount: int(le.Amount)}.Resolve(
			&EffectContext{
				Resolver:   g,
				Source:     subject,
				Controller: actor,
			},
		)
	case actGiveRemainingAember:
		beneficiary := 1 - actor
		amount := g.State.Aember[actor]
		g.SetAember(actor, 0)
		g.SetAember(beneficiary, g.State.Aember[beneficiary]+amount)
		g.record(AemberGivenAfterForging{
			Player: actor,
			To:     beneficiary,
			Amount: amount,
		})
	default: // actGainAember
		if capturer, ok := g.gainAember(actor, int(le.Amount)); ok {
			g.record(AemberCapturedInsteadOfGain{
				Creature: capturer,
				Player:   actor,
				Amount:   int(le.Amount),
			})
			return
		}
		g.record(AemberGained{Player: actor, Amount: int(le.Amount)})
	}
}

// lastingReplacement returns the replacement a player has for a replacement event,
// or ok=false when none is active. It is how an event site (reapWith) asks whether
// its outcome is being replaced this turn.
func (g *Game) lastingReplacement(player int, event Event) (lastingAction, bool) {
	for i := 0; i < int(g.State.LastingCount); i++ {
		if le := g.State.Lasting[i]; int(le.Controller) == player && le.On == event {
			return le.Do, true
		}
	}
	return 0, false
}

// consumeBonusIconBoost reports whether player has an armed one-shot bonus-icon
// boost (Wild Bounty) and, if so, removes it — a boost fires for a single played
// card. The play-path icon loop uses it to decide whether to resolve each icon an
// additional time.
func (g *Game) consumeBonusIconBoost(player int) bool {
	for i := 0; i < int(g.State.LastingCount); i++ {
		le := g.State.Lasting[i]
		if int(le.Controller) != player || le.Do != actResolveBonusAgain {
			continue
		}
		g.removeLastingAt(i)
		return true
	}
	return false
}

// removeLastingAt drops the lasting effect at index i, compacting the fixed array
// and clearing the freed tail slot so the state stays a flat value.
func (g *Game) removeLastingAt(i int) {
	last := int(g.State.LastingCount) - 1
	for j := i; j < last; j++ {
		g.State.Lasting[j] = g.State.Lasting[j+1]
	}
	g.State.Lasting[last] = LastingEffect{}
	g.State.LastingCount--
}

// lastingExtraDamage sums the additional damage a creature takes from every
// modifier keyed to it — Lethal Distraction adds 2 to each instance of damage the
// chosen creature takes for the rest of the turn.
func (g *Game) lastingExtraDamage(id LocalID) int {
	total := 0
	for i := 0; i < int(g.State.LastingCount); i++ {
		le := g.State.Lasting[i]
		if le.Do == actTakeExtraDamage && le.HasSubject && le.Subject == id {
			total += int(le.Amount)
		}
	}
	return total
}
