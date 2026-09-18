package engine

import (
	"fmt"
	"strings"
)

// This file holds the authoring effects for lasting "for the remainder of the
// turn" behavior: ForRemainderOfTurn installs a reaction (do something after an
// event), Instead installs a replacement (change an event's own outcome). Both are
// thin: they translate a composed effect into a flat lasting record via AddLasting
// (see game_lasting.go for how the records are fired and queried).

// ForRemainderOfTurn installs a reaction that runs for the rest of the controller's
// turn each time On occurs — Full Moon (gain Æmber when you play a creature),
// Charge! (deal damage when you play a creature), Crystal Hive (gain Æmber after a
// creature reaps), Library Access (draw a card each time you play another card).
// Do is the effect that runs each time; the supported effects are GainAember
// (crediting the controller), DealDamage to an enemy creature, CaptureAember,
// GiveAember, and Draw. The card that installs the reaction never triggers it
// itself, so an event phrased as "another card" excludes the play that armed it.
type ForRemainderOfTurn struct {
	On Event
	Do Effect
}

// validate rejects a non-reaction event or a Do the reaction cannot carry.
func (e ForRemainderOfTurn) validate() error {
	return validateLastingReaction("ForRemainderOfTurn", e.On, e.Do)
}

// Text renders the effect, e.g. "for the remainder of the turn, each time you play
// a creature, gain 1 Æmber".
func (e ForRemainderOfTurn) Text() string {
	return durationClause(RemainderOfPlayerTurn, "") + ", " + e.On.clause() + ", " + e.Do.Text()
}

// Resolve registers the reaction on the controller for the rest of their turn. On
// EventCardPlayed the installing card is excepted, so "each time you play another
// card" does not count the play that armed it.
func (e ForRemainderOfTurn) Resolve(ctx *EffectContext) {
	installLastingReaction(ctx, e.On, e.Do, ctx.Controller)
}

// ForOpponentNextTurn installs a reaction that lies dormant for the rest of this
// turn and fires only during the opponent's next turn, clearing at the end of that
// turn — Interdimensional Graft (after the opponent forges a key, they give the
// controller all their Æmber). It carries the same reaction Do as ForRemainderOfTurn
// but arms it on the opponent, so the window is the opponent's next turn rather than
// the rest of this one.
type ForOpponentNextTurn struct {
	On Event
	Do Effect
}

// validate rejects a non-reaction event or a Do the reaction cannot carry.
func (e ForOpponentNextTurn) validate() error {
	return validateLastingReaction("ForOpponentNextTurn", e.On, e.Do)
}

// Text renders the effect, e.g. "during your opponent's next turn, after forging a
// key, your opponent gives you all their Æmber".
func (e ForOpponentNextTurn) Text() string {
	return durationClause(OpponentNextTurn, "") + ", " +
		e.On.clauseOnOpponentTurn() + ", " + e.Do.Text()
}

// Resolve arms the reaction on the opponent, so it lies dormant this turn and fires
// on the opponent's next turn.
func (e ForOpponentNextTurn) Resolve(ctx *EffectContext) {
	installLastingReaction(ctx, e.On, e.Do, ctx.Opponent())
}

// validateLastingReaction is the shared gate for the lasting-reaction nodes: the
// event must be a reaction, and Do must be an effect the flat registry can carry
// (and, for DealDamage, must target an enemy creature). name labels the error.
func validateLastingReaction(name string, on Event, do Effect) error {
	if !on.isReaction() {
		return fmt.Errorf("%s: On must be a reaction event", name)
	}
	if _, _, ok := lastingActionOf(do); !ok {
		return fmt.Errorf("%s: unsupported Do %T", name, do)
	}
	if d, ok := do.(DealDamage); ok && d.Target != (Target{Kind: TargetChosenEnemyCreature}) {
		return fmt.Errorf("%s: DealDamage must target an enemy creature", name)
	}
	return validateEffect(do)
}

// installLastingReaction registers a reaction owned by owner (whose ready phase
// clears it). On EventCardPlayed the installing card is excepted, so an event
// phrased as "another card" does not count the play that armed it.
func installLastingReaction(ctx *EffectContext, on Event, do Effect, owner int) {
	action, amount, _ := lastingActionOf(do)
	ctx.Resolver.AddLasting(LastingEffect{
		On:         on,
		Do:         action,
		Controller: int8(owner),
		Amount:     int8(amount),
		Except:     ctx.Source,
		HasExcept:  on == EventCardPlayed,
		Source:     ctx.Source,
		HasSource:  true,
	})
}

// lastingActionOf maps a reaction's Do effect to the flat action and amount stored
// in the registry, reporting whether the effect is one the registry can carry.
func lastingActionOf(e Effect) (lastingAction, int, bool) {
	switch d := e.(type) {
	case GainAember:
		return actGainAember, d.Amount, true
	case LoseAember:
		return actLoseAember, d.Amount, true
	case StealAember:
		return actSteal, d.Amount, true
	case DealDamage:
		return actDealDamage, d.Amount, true
	case CaptureAember:
		if d.Target.Kind == TargetChosenFriendlyCreature {
			return actCaptureChosen, d.Amount, true
		}
		return actCapture, d.Amount, true
	case GiveAember:
		return actGiveRemainingAember, d.Amount, true
	case Draw:
		return actDraw, d.Amount, true
	case Ready:
		return actReadyPlayed, 0, true
	case Exalt:
		return actExalt, d.Amount, true
	case Stun:
		return actStun, 0, true
	}
	return 0, 0, false
}

// reactionEventOf maps a triggered ability's trigger to the reaction event that
// fires it, reporting whether the trigger is one a lasting per-creature grant can
// hang on. Reap (Spectral Tunneler) and Fight (Into the Fray) are supported, as is
// Before Fight (Diplomacy); it lives here so GainAbility and the registry agree on
// the mapping.
func reactionEventOf(t Trigger) (Event, bool) {
	switch t {
	case TriggerAfterReap:
		return EventReap, true
	case TriggerAfterFight:
		return EventFight, true
	case TriggerBeforeFight:
		return EventBeforeFight, true
	}
	return eventUnset, false
}

// GainAbility grants each creature its Target selects a triggered ability for the
// remainder of the turn — Spectral Tunneler grants a chosen creature "Reap: Draw a
// card". Unlike ForRemainderOfTurn's controller-wide reaction, this one is scoped
// to the single granted creature through the registry's Subject, so only that
// creature's own trigger fires it. The ability is stored flat, so only triggers
// and effects the registry can carry (a Reap, Fight, or Before Fight reaction whose
// effect is a Draw, damage, capture, Ready, or Exalt) are allowed.
//
// Duration widens the window past the current turn. Unset (the zero value) is the
// remainder of the controller's turn. StartOfPlayerNextTurn lasts until the
// start of the controller's next turn — through the opponent's whole turn — so an
// enemy creature fires the grant on the opponent's turn too (Diplomacy). Because
// the registry clears a player's own entries at their ready phase, a next-turn
// grant is owned by the opponent, whose ready phase falls just before the
// controller's next turn, and it fires only for a Before Fight ability, whose
// firing (fireLastingBeforeFight) matches on the subject alone rather than the
// acting player.
type GainAbility struct {
	Target   Target
	Ability  Ability
	Duration Duration
}

// validate requires a target and an ability the flat registry can carry. A
// next-turn grant is limited to a Before Fight ability, the only trigger whose
// firing ignores which player is acting.
func (e GainAbility) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("GainAbility")
	}
	if _, ok := reactionEventOf(e.Ability.Trigger); !ok {
		return fmt.Errorf("GainAbility: unsupported trigger %v", e.Ability.Trigger)
	}
	if _, _, ok := lastingActionOf(e.Ability.Effect); !ok {
		return fmt.Errorf("GainAbility: unsupported ability effect %T", e.Ability.Effect)
	}
	if e.Duration == StartOfPlayerNextTurn && e.Ability.Trigger != TriggerBeforeFight {
		return fmt.Errorf(
			"GainAbility: StartOfPlayerNextTurn grant supports only a Before Fight ability, got %v",
			e.Ability.Trigger,
		)
	}
	return validateEffect(e.Ability.Effect)
}

// Text renders the effect, e.g. `it gains, "Reap: Draw a card."` — a granted
// ability takes a comma and quotes (card-wording rule 2), the period inside. A
// self-reference in the granted ability names the creature that gains it, so it
// renders "this creature" rather than the source card's name. A next-turn grant
// names its window first (Diplomacy); a RemainderOfPlayerTurn grant names its
// window first too when the card renders the duration itself (Adaptoid), rather
// than leaning on a sibling effect to carry the shared clause (Spectral Tunneler).
func (e GainAbility) Text() string {
	granted := strings.ReplaceAll(RenderAbility(e.Ability), SelfName, "this creature")
	text := e.Target.Text() + ` gains, "` + granted + `"`
	switch e.Duration {
	case StartOfPlayerNextTurn:
		return durationClause(e.Duration, "") + ", " + text
	case RemainderOfPlayerTurn:
		return durationClause(e.Duration, "") + ", " + text
	}
	return text
}

// Resolve registers the ability as a per-creature reaction on each selected
// creature for the rest of the controller's turn. A next-turn grant is owned by
// the opponent so it clears at their ready phase — the start of the controller's
// next turn — rather than at the end of this one.
func (e GainAbility) Resolve(ctx *EffectContext) {
	event, _ := reactionEventOf(e.Ability.Trigger)
	action, amount, _ := lastingActionOf(e.Ability.Effect)
	owner := ctx.Controller
	if e.Duration == StartOfPlayerNextTurn {
		owner = ctx.Opponent()
	}
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddLasting(LastingEffect{
			On:         event,
			Do:         action,
			Controller: int8(owner),
			Amount:     int8(amount),
			Subject:    id,
			HasSubject: true,
			Source:     ctx.Source,
			HasSource:  true,
		})
	}
}

// TakesExtraDamage makes each creature its Target selects take an additional
// Amount damage whenever it takes damage, for the rest of the controller's turn —
// Lethal Distraction's "for the remainder of the turn, whenever this creature takes
// damage, it takes an additional 2 damage". It is the MODIFIER flavor of the
// lasting spine (see game_lasting.go): subject-scoped to the chosen creature and
// summed at the damage site rather than swapping an outcome like Instead.
type TakesExtraDamage struct {
	Target Target
	Amount int
}

// validate requires a target and a positive Amount.
func (e TakesExtraDamage) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("TakesExtraDamage")
	}
	if e.Amount <= 0 {
		return fmt.Errorf("TakesExtraDamage: Amount must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, whenever it takes
// damage, it takes an additional 2 damage".
func (e TakesExtraDamage) Text() string {
	return fmt.Sprintf(
		"%s, whenever %s takes damage, it takes an additional %d damage",
		durationClause(RemainderOfPlayerTurn, ""), e.Target.Text(), e.Amount)
}

// Resolve registers the modifier on each selected creature for the rest of the
// controller's turn.
func (e TakesExtraDamage) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddLasting(LastingEffect{
			On:         EventCreatureTakesDamage,
			Do:         actTakeExtraDamage,
			Controller: int8(ctx.Controller),
			Amount:     int8(e.Amount),
			Subject:    id,
			HasSubject: true,
		})
	}
}

// Replacement is a lasting change to an event's own outcome, used by Instead.
type Replacement uint8

const (
	// replacementUnset is the invalid zero value: an Instead must name its	// replacement rather than leave it unset.
	replacementUnset Replacement = iota
	// Steal replaces gaining Æmber with stealing that much from the opponent.
	Steal
	// Capture replaces adding Æmber to a pool with the source creature capturing it
	// (Ether Spider). It is applied continuously at the add-to-pool site, not through
	// the turn-scoped lasting registry, so it has no lastingAction.
	Capture
	// FromCommonSupply replaces the source of a steal or capture: the Æmber is drawn
	// from the common supply instead of the target's pool (Po's Pixies). Like Capture
	// it is applied continuously at the take site, never through the turn-scoped
	// registry, so it has no lastingAction.
	FromCommonSupply
)

// valid reports whether r names a real replacement (not the unset zero value).
func (r Replacement) valid() bool { return r != replacementUnset }

// action maps the replacement to the flat action stored in the registry.
func (Replacement) action() lastingAction { return actSteal }

// text renders the replacement clause, e.g. "steal the same amount".
func (Replacement) text() string { return "steal the same amount" }

// Replace is a continuous replacement a card applies to a game event while in
// play: when the event When would happen, the effect With resolves in its place.
// An Upgrade carries it to replace the event for its host (Armageddon Cloak
// replaces its host's destruction with "fully heal it and destroy Armageddon
// Cloak"); a creature carries it to replace the event for itself (Reassembling
// Automaton replaces its own destruction with "fully heal it, exhaust it, and move
// it to a flank"). Unlike the turn-scoped Instead — a flat outcome swap kept in
// the pointerless game state — a Replace lives in the card definition, so its With
// is a full effect tree. When Cond is set, the replacement applies only while that
// condition holds; a nil Cond always applies.
type Replace struct {
	When Event
	Cond Condition
	With Effect
}

// valid reports whether a replacement is set, distinguishing a StaticModifier that
// carries a Replace from the zero value that carries none.
func (r Replace) valid() bool { return r.When != eventUnset }

// validate surfaces a configuration error in the replacement effect (and its
// condition, if any), ignoring the zero value (a StaticModifier with no
// replacement).
func (r Replace) validate() error {
	if !r.valid() {
		return nil
	}
	if r.Cond != nil {
		if err := validateCondition(r.Cond); err != nil {
			return err
		}
	}
	return validateEffect(r.With)
}

// Instead installs a replacement that, for the rest of the controller's turn,
// changes the outcome of the event Of before it happens — Dimension Door replaces
// gaining Æmber from reaping with stealing it. As a plain {Of, With} value it also
// describes a continuous replacement a card applies while in play (CardDefinition.Replaces,
// Ether Spider capturing Æmber added to its opponent's pool); in that use it is read,
// never resolved. Player scopes an event that names a pool — which player's pool the
// replacement watches (Ether Spider watches its Opponent's).
type Instead struct {
	Of     Event
	With   Replacement
	Player Player
}

// valid reports whether a replacement is set, distinguishing a card that carries a
// continuous Instead from the zero value that carries none.
func (e Instead) valid() bool { return e.Of != eventUnset }

// validate rejects an Of that is not a replacement event or an unset With, and
// requires the pool-scoping Player when Of names a pool.
func (e Instead) validate() error {
	if e.Of.isReaction() {
		return fmt.Errorf("Instead: Of must be a replacement event")
	}
	// Text() names the replaced event by its gerund, so an event without one cannot
	// be rendered and would print another event's wording.
	if e.Of.gerund() == "" {
		return fmt.Errorf("Instead: %d has no gerund to render", e.Of)
	}
	if !e.With.valid() {
		return fmt.Errorf("Instead: replacement must be set")
	}
	if (e.Of == EventAemberAddedToPool || e.Of == EventAemberTakenFromPool) &&
		!e.Player.valid() {
		return fmt.Errorf("Instead: a pool event needs a Player to scope it")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, instead of gaining
// Æmber from reaping, steal the same amount".
func (e Instead) Text() string {
	return durationClause(RemainderOfPlayerTurn, "") + ", instead of " +
		e.Of.gerund() + ", " + e.With.text()
}

// Resolve registers the replacement on the controller for the rest of the turn.
func (e Instead) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLasting(LastingEffect{
		On:         e.Of,
		Do:         e.With.action(),
		Controller: int8(ctx.Controller),
	})
}

// DamageOthersAfterUsingTrait installs a reaction that runs for the rest of the
// controller's turn: each time they use a creature carrying Trait, it deals Amount
// damage to each creature that lacks Trait, on both battlelines — Legion's March
// ("after you use a Dinosaur creature, deal 1 damage to each non-Dinosaur
// creature").
type DamageOthersAfterUsingTrait struct {
	Trait  Trait
	Amount int
}

// validate requires a trait to gate the use and a positive damage amount.
func (e DamageOthersAfterUsingTrait) validate() error {
	if e.Trait == traitUnset {
		return fmt.Errorf("DamageOthersAfterUsingTrait: trait must be set")
	}
	if e.Amount <= 0 {
		return fmt.Errorf("DamageOthersAfterUsingTrait: amount must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, after you use a
// Dinosaur creature, deal 1 damage to each non-Dinosaur creature".
func (e DamageOthersAfterUsingTrait) Text() string {
	damage := DealDamage{
		Amount: e.Amount,
		Target: Target{Kind: TargetEachCreature}.ExceptTrait(e.Trait),
	}
	return durationClause(RemainderOfPlayerTurn, "") + ", after you use a " +
		e.Trait.String() + " creature, " + damage.Text()
}

// Resolve registers the reaction on the controller for the rest of their turn.
func (e DamageOthersAfterUsingTrait) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLasting(LastingEffect{
		On:         EventUsed,
		Do:         actDamageOthersOfTrait,
		Controller: int8(ctx.Controller),
		Amount:     int8(e.Amount),
		Trait:      e.Trait,
		Source:     ctx.Source,
		HasSource:  true,
	})
}

// PutNextTacticIntoHand makes the next action card its controller resolves this
// turn return to their hand instead of their discard pile — High Priest Torvus,
// once exalted, sends its controller's next action back to hand. It registers a
// one-shot arming the action-play path consumes; the turn's end clears it if no
// action card resolves.
type PutNextTacticIntoHand struct{}

// Text renders the effect.
func (PutNextTacticIntoHand) Text() string {
	return "after you resolve your next tactic this turn, put it into your " +
		"hand instead of your discard pile"
}

// Resolve registers the one-shot redirect on the controller.
func (PutNextTacticIntoHand) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLasting(LastingEffect{
		On:         EventNextTacticIntoHand,
		Do:         actPutIntoHand,
		Controller: int8(ctx.Controller),
		Once:       true,
	})
}

// NextPlayed makes the next card its controller plays this turn that matches its
// filters enter play with EntersPlay applied to it — Blypyp readying the next Mars
// creature, Soft Landing readying the next creature or artifact. It registers a
// one-shot reaction that applies the effect to that card and then removes itself;
// the turn's end clears it if no matching card is played. EntersPlay must be an
// enter-play effect the flat registry can carry (Ready).
//
// Of narrows the card to one house, and is optional. Type narrows it to one card
// type; AnyType means "creature or artifact", the two types that stay in play.
type NextPlayed struct {
	Of         HouseMatcher
	Type       CardType
	EntersPlay Effect
}

// validate requires a card type and rejects an EntersPlay effect the registry
// cannot carry. The house matcher must be context-free (named or any), since the
// registry outlives the resolution that armed it.
func (e NextPlayed) validate() error {
	if e.Type == TypeUnset {
		return fmt.Errorf("NextPlayed: type must be set")
	}
	switch e.Of.Kind {
	case MatchAnyHouse, MatchNamedHouse, MatchExceptHouse:
	default:
		return fmt.Errorf("NextPlayed: house matcher %v needs resolution context", e.Of.Kind)
	}
	if err := e.Of.validate(); err != nil {
		return err
	}
	if _, ok := enterActionOf(e.EntersPlay); !ok {
		return fmt.Errorf("NextPlayed: unsupported EntersPlay %T", e.EntersPlay)
	}
	return validateEffect(e.EntersPlay)
}

// Text renders the effect, e.g. "the next Mars creature you play this turn enters
// play ready" or "the next creature or artifact you play this turn enters play
// ready".
func (e NextPlayed) Text() string {
	noun := "creature or artifact"
	if e.Type != AnyType {
		noun = typeWord(e.Type)
	}
	noun = e.Of.qualifyNoun(noun)
	return fmt.Sprintf(
		"the next %s you play this turn enters play %s",
		noun,
		enterStateWord(e.EntersPlay),
	)
}

// Resolve registers the one-shot enter-play reaction on the controller.
func (e NextPlayed) Resolve(ctx *EffectContext) {
	action, _ := enterActionOf(e.EntersPlay)
	ctx.Resolver.AddLasting(LastingEffect{
		On:         EventCardEntersPlay,
		Do:         action,
		Controller: int8(ctx.Controller),
		House:      e.Of,
		Type:       e.Type,
		Once:       true,
		Source:     ctx.Source,
		HasSource:  true,
	})
}

// enterActionOf maps an enter-play effect to the flat lasting action that applies
// it to the next played creature, reporting whether the registry can carry it.
func enterActionOf(e Effect) (lastingAction, bool) {
	if _, ok := e.(Ready); ok {
		return actReadyPlayed, true
	}
	return 0, false
}
