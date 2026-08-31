package engine

// This file holds LASTING EFFECTS: "for the remainder of the turn" effects that
// attach to a later game event. Two flavors share the same flat registry:
//
//   - a REACTION runs after an event (Full Moon gains Æmber after you play a
//     creature, Charge! deals damage after you play a creature, Crystal Hive gains
//     Æmber after a creature reaps). Reactions are fired by emitLasting, which
//     gathers every reaction responding to an event and lets the controller order
//     them when several fire at once.
//   - a REPLACEMENT changes an event's own outcome before it happens (Dimension
//     Door makes reaping steal Æmber instead of gaining it). The event site queries
//     for a replacement (lastingReplacement) and applies it in place.
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
	// EventForgeKey fires after a player forges a key (a reaction point). A reaction
	// owned by that player fires during their turn; because the registry clears a
	// player's own entries at their EndTurn, a reaction armed on an opponent (owned
	// by the forger) survives the arming turn and fires during the forger's next
	// turn — the "during your opponent's next turn" window (Interdimensional Graft).
	EventForgeKey
)

// isReaction reports whether the event is a reaction point (fired after) rather
// than a replacement point (queried during).
func (e Event) isReaction() bool {
	return e == EventCreaturePlayed || e == EventReap || e == EventFight ||
		e == EventEnemyCreatureDestroyed || e == EventForgeKey
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
	default:
		return "each time you play a creature"
	}
}

// gerund renders the "instead of ..." phrase for a replacement, e.g. "gaining
// Æmber from reaping".
func (e Event) gerund() string {
	return "gaining Æmber from reaping"
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
	default:
		return "gain Æmber"
	}
}

// LastingEffect is one registered lasting effect: the event it attaches to, what it
// does, the player it belongs to, and its magnitude. It is a plain value so the
// flat GameState can hold a fixed array of them.
type LastingEffect struct {
	On         Event
	Do         lastingAction
	Controller int8
	Amount     int8
	// House, when set, limits a reaction to a subject of that house (Blypyp readies
	// only Mars creatures). HouseNone (the zero value) reacts to any subject.
	House House
	// Once removes the record after it fires a single time — "the next" rather than
	// "each time".
	Once bool
}

// maxLasting bounds how many lasting effects can be active at once — generous for
// the handful of "remainder of the turn" cards, and a fixed size keeps the state a
// flat value.
const maxLasting = 8

// AddLasting registers a lasting effect owned by controller for the rest of their
// turn. It is the single seam a "for the remainder of the turn" effect uses instead
// of hardcoding itself into the play or reap path.
func (g *Game) AddLasting(on Event, do lastingAction, controller, amount int) {
	g.addLasting(LastingEffect{On: on, Do: do, Controller: int8(controller), Amount: int8(amount)})
}

// AddLastingOnce registers a one-shot reaction that fires the next time its event
// occurs — and, when house is set, only for a subject of that house — then removes
// itself (Blypyp readying the next Mars creature you play).
func (g *Game) AddLastingOnce(on Event, do lastingAction, controller, amount int, house House) {
	g.addLasting(
		LastingEffect{
			On:         on,
			Do:         do,
			Controller: int8(controller),
			Amount:     int8(amount),
			House:      house,
			Once:       true,
		},
	)
}

// addLasting appends a lasting record, dropping it silently when the registry is
// full.
func (g *Game) addLasting(le LastingEffect) {
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
}

// emitLasting resolves every reaction actor owns that responds to event. When
// several fire at once the controller chooses the order (KeyForge lets the active
// player order simultaneous triggers). subject is the card that caused the event —
// the played creature — for reactions that need a source.
func (g *Game) emitLasting(event Event, actor int, subject LocalID) {
	var pending []LastingEffect
	for i := 0; i < int(g.State.LastingCount); i++ {
		le := g.State.Lasting[i]
		if int(le.Controller) != actor || le.On != event {
			continue
		}
		if le.House != HouseNone && g.cat.def(subject).House != le.House {
			continue
		}
		pending = append(pending, le)
	}
	for len(pending) > 0 {
		idx := 0
		if len(pending) > 1 {
			labels := make([]string, len(pending))
			for i, le := range pending {
				labels[i] = le.Do.describe()
			}
			idx = g.chooseOption(actor, "", "Choose the next effect to resolve", labels)
		}
		le := pending[idx]
		g.resolveReaction(le, actor, subject)
		if le.Once {
			g.removeLasting(le)
		}
		pending = append(pending[:idx], pending[idx+1:]...)
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
		g.logf("%s is readied", g.Name(subject))
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
	case actGiveRemainingAember:
		beneficiary := 1 - actor
		amount := g.State.Aember[actor]
		g.State.Aember[actor] = 0
		g.State.Aember[beneficiary] += amount
		g.logf(
			"%s gives %d Æmber to %s after forging a key",
			g.names[actor],
			amount,
			g.names[beneficiary],
		)
	default: // actGainAember
		if capturer, ok := g.gainAember(actor, int(le.Amount)); ok {
			g.logf(
				"%s captures %d Æmber instead of %s gaining it (%s)",
				g.Name(capturer),
				le.Amount,
				g.names[actor],
				le.On.clause(),
			)
			return
		}
		g.logf("%s gains %d Æmber (%s)", g.names[actor], le.Amount, le.On.clause())
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
