package engine

import "fmt"

// A conditional gates an effect behind a check on the current game state — the
// "If ..." clause a card opens with, e.g. "If your opponent has 7 or more Aember,
// they lose 4 Aember." The effect resolves only when the condition is met. Unlike
// a result gate (A -> B), which turns on an action succeeding, a conditional turns
// on a fact about the board.
//
//rulebook:effect Conditional

// Condition is a boolean predicate on the live game, used by Conditional.
type Condition interface {
	CondText() string
	Met(ctx *EffectContext) bool
}

// Comparison selects how a condition compares a quantity to a threshold. It has
// no valid zero value: a condition must name one, so an unset comparison is
// caught at init rather than silently reading as "0 or more".
type Comparison int

const (
	comparisonUnset Comparison = iota
	// AtLeast is met when the quantity is at least Amount.
	AtLeast
	// Exactly is met when the quantity is exactly Amount (which is why the
	// comparison is named separately from the amount: Exactly with Amount 0 is a
	// real check that a bare integer field could not tell from "unset").
	Exactly
	// MoreThanYou is met when the opponent's pool holds strictly more Æmber than the
	// controller's; it ignores Amount and applies only to OpponentAember.
	MoreThanYou
)

// OpponentAember gates on the opponent's Æmber pool: Is names the comparison and
// Amount the threshold it compares against (unused by MoreThanYou). It replaces
// the separate at-least / exactly / more-than-you conditions with one node.
type OpponentAember struct {
	Is     Comparison
	Amount int
}

// validate requires a comparison to be named (its zero value is invalid).
func (c OpponentAember) validate() error {
	switch c.Is {
	case AtLeast, Exactly, MoreThanYou:
		return nil
	default:
		return fmt.Errorf("OpponentAember: Is must be AtLeast, Exactly, or MoreThanYou")
	}
}

// CondText renders the condition, e.g. "if your opponent has 7 Æmber or more".
func (c OpponentAember) CondText() string {
	switch {
	case c.Is == Exactly && c.Amount == 0:
		return "if your opponent has no Æmber"
	case c.Is == Exactly:
		return fmt.Sprintf("if your opponent has exactly %d Æmber", c.Amount)
	case c.Is == MoreThanYou:
		return "if your opponent has more Æmber than you"
	default:
		return fmt.Sprintf("if your opponent has %d Æmber or more", c.Amount)
	}
}

// Met reports whether the opponent's pool satisfies the comparison.
func (c OpponentAember) Met(ctx *EffectContext) bool {
	opp := ctx.Resolver.Aember(ctx.Opponent())
	switch c.Is {
	case Exactly:
		return opp == c.Amount
	case MoreThanYou:
		return opp > ctx.Resolver.Aember(ctx.Controller)
	default:
		return opp >= c.Amount
	}
}

// validateCondition returns any configuration error a condition reports (an unset
// OpponentAember comparison, say). Conditions that cannot be misconfigured
// implement no validator and pass.
func validateCondition(c Condition) error {
	if v, ok := c.(validator); ok {
		return v.validate()
	}
	return nil
}

// ControlsMoreCreatures is met while the controller has more creatures in play
// than the opponent.
type ControlsMoreCreatures struct{}

// CondText renders the condition.
func (ControlsMoreCreatures) CondText() string {
	return "if you control more creatures than your opponent"
}

// Met reports whether the controller has more creatures in play than the opponent.
func (ControlsMoreCreatures) Met(ctx *EffectContext) bool {
	return len(
		ctx.Resolver.Battleline(ctx.Controller),
	) > len(
		ctx.Resolver.Battleline(ctx.Opponent()),
	)
}

// FirstCreaturePlayedThisTurn is met when the card in context (ctx.It, the
// creature that fired the trigger) is the first creature its player played this
// turn — Speed Sigil readies it. It is a once-per-turn charge that needs no state
// of its own: the turn's play record already says whether the charge is spent, and
// the record is cleared when the next turn begins.
//
// A creature put into play by an effect rather than played never matches, so it
// neither benefits nor spends the charge.
type FirstCreaturePlayedThisTurn struct{}

// CondText renders the condition.
func (FirstCreaturePlayedThisTurn) CondText() string {
	return "if it is the first creature played this turn"
}

// Met reports whether the context card is the earliest creature in the active
// player's plays this turn.
func (FirstCreaturePlayedThisTurn) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	for _, id := range ctx.Resolver.PlayedThisTurn(ctx.Resolver.ActivePlayer()) {
		if ctx.Resolver.TypeOf(id) == Creature {
			return id == ctx.It
		}
	}
	return false
}

// CardsDestroyedFewerThan is met when fewer than Amount cards were destroyed this
// way — the tally a preceding effect records on the context. Bonkers Killing
// Machine destroys itself when its house-driven destruction removed fewer than two.
type CardsDestroyedFewerThan struct {
	Amount int
}

// CondText renders the condition, e.g. "if fewer than 2 cards are destroyed this
// way".
func (c CardsDestroyedFewerThan) CondText() string {
	return fmt.Sprintf("if fewer than %d cards are destroyed this way", c.Amount)
}

// Met reports whether fewer than Amount cards were destroyed this way.
func (c CardsDestroyedFewerThan) Met(ctx *EffectContext) bool {
	return ctx.Produced.TotalDestroyed() < c.Amount
}

// HouseChoice names a house a condition compares against by reference rather than
// by a fixed value — the house chosen this turn, or the active house — so one
// condition works wherever such a house is meaningful.
type HouseChoice uint8

const (
	// houseChoiceUnset is the invalid zero value.
	houseChoiceUnset HouseChoice = iota
	// TheChosenHouse is the house picked by an enclosing ChooseHouseThen.
	TheChosenHouse
	// TheActiveHouse is the player's active house this turn.
	TheActiveHouse
	// TheContextualHouse is the house of the card in context (ctx.It) — e.g. a
	// just-discarded deck card, for "of the discarded card's house".
	TheContextualHouse
	// AnyHouse names no house at all, so a count or condition using it applies no
	// house filter — Key Abduction counts every card in hand, whatever its house.
	AnyHouse
)

// resolveHouse turns a HouseChoice into the concrete house it names in the current
// context (HouseNone when it names a card in context and none is present).
func (h HouseChoice) resolveHouse(ctx *EffectContext) House {
	switch h {
	case TheActiveHouse:
		return ctx.Resolver.ActiveHouse()
	case TheContextualHouse:
		if ctx.HasIt {
			return ctx.Resolver.House(ctx.It)
		}
		return HouseNone
	default:
		return ctx.ChosenHouse
	}
}

// ItIsOfHouse is met when the card in context (ctx.It — a revealed, discarded, or
// triggering card) belongs to a referenced house. It replaces the one-off
// "revealed card of the chosen house" (Chaos Portal) and "discarded card of the
// active house" (Evasion Sigil) with a single filter on the contextual card.
type ItIsOfHouse struct {
	House HouseChoice
}

// CondText renders the condition, e.g. "if it is of the chosen house".
func (e ItIsOfHouse) CondText() string {
	if e.House == TheActiveHouse {
		return "if it is of the active house"
	}
	return "if it is of the chosen house"
}

// Met reports whether a card is in context and belongs to the referenced house.
func (e ItIsOfHouse) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	house := ctx.Resolver.House(ctx.It)
	if e.House == TheActiveHouse {
		return house == ctx.Resolver.ActiveHouse()
	}
	return house == ctx.ChosenHouse
}

// ItIs is met when the card in context (ctx.It — a just-played, revealed, or
// discarded card) matches a concrete House and/or Type filter, e.g. "if it is a
// Mars creature" (Brain Stem Antenna reacting to a played card) or "if it is an
// artifact" (Carlo Phantom). Either filter may be left unset to match any. It is
// the concrete-value counterpart to ItIsOfHouse, which names the house by
// reference (the chosen or active house).
type ItIs struct {
	House House
	Type  CardType
}

// CondText renders the condition, e.g. "if it is a Mars creature" or "if it is an
// artifact".
func (e ItIs) CondText() string {
	return "if it is " + indefinite(houseTypeNoun(e.House, e.Type))
}

// Met reports whether a card is in context and matches the house and type filters.
func (e ItIs) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	if e.House != HouseNone && ctx.Resolver.House(ctx.It) != e.House {
		return false
	}
	if e.Type != TypeUnset && ctx.Resolver.TypeOf(ctx.It) != e.Type {
		return false
	}
	return true
}

// ItIsOffIdentity is met when the card in context (ctx.It) belongs to none of the
// controller's identity houses — the three houses of their deck. Sneklifter uses
// it to reassign a seized enemy artifact to Shadows only when it is off your
// identity; KeyForge phrases the check negatively ("if it does not belong to one
// of your three houses").
type ItIsOffIdentity struct{}

// CondText renders the condition.
func (ItIsOffIdentity) CondText() string {
	return "if it does not belong to a house on your identity"
}

// Met reports whether a card is in context and belongs to none of the controller's
// identity houses.
func (ItIsOffIdentity) Met(ctx *EffectContext) bool {
	return ctx.HasIt && !ctx.Resolver.PlayerHasHouse(ctx.Controller, ctx.Resolver.House(ctx.It))
}

// houseTypeNoun renders a card filtered by house and type as a noun, e.g. "Mars
// creature", "artifact", or the bare "card" when neither is set.
func houseTypeNoun(house House, typ CardType) string {
	n := "card"
	switch typ {
	case Creature:
		n = "creature"
	case Artifact:
		n = "artifact"
	}
	if house != HouseNone {
		n = house.String() + " " + n
	}
	return n
}

// Conditional resolves Then only when Cond is met. It renders as "<cond>, <then>",
// e.g. "if your opponent has 7 Æmber or more, your opponent loses 4 Æmber".
//
// Else, when set, resolves instead when Cond is not met, and leads the sentence so
// the common case reads first: "<else>. <Cond>, <then> instead." Key of Darkness
// forges at +6, or at +2 when the opponent has no Æmber.
type Conditional struct {
	Cond Condition
	Then Effect
	Else Effect
}

// Text joins the condition and the gated effect.
func (e Conditional) Text() string {
	if e.Else != nil {
		return e.Else.Text() + ". " + capitalizeFirst(e.Cond.CondText()) + ", " +
			e.Then.Text() + " instead"
	}
	return e.Cond.CondText() + ", " + e.Then.Text()
}

// Resolve runs Then when Cond is met, otherwise Else if one is set.
func (e Conditional) Resolve(ctx *EffectContext) {
	if e.Cond.Met(ctx) {
		e.Then.Resolve(ctx)
		return
	}
	if e.Else != nil {
		e.Else.Resolve(ctx)
	}
}

// validate checks the gated effect for configuration errors.
func (e Conditional) validate() error {
	if err := validateCondition(e.Cond); err != nil {
		return err
	}
	if err := validateEffect(e.Then); err != nil {
		return err
	}
	if e.Else == nil {
		return nil
	}
	return validateEffect(e.Else)
}

// RuleOfSix is the most times a card can be played, used, or made to resolve
// again in one turn. A self-repeating effect is bounded by it: Bait and Switch's
// "steal 1 Æmber -> repeat this effect" resolves the initial steal plus at most
// five repeats, so it steals six at most however far ahead the opponent is.
const RuleOfSix = 6

// RepeatWhile resolves Do again and again for as long as Cond holds, re-checking
// after each pass — a self-looping effect such as "if your opponent has more
// Æmber than you, steal 1 Æmber -> repeat this effect". The loop also stops the
// moment Do makes no progress (its gate reports it did nothing), so an action that
// is prevented — a steal against a protected pool — ends the loop instead of
// spinning even though Cond still holds. Do is a GatingEffect for exactly that
// reason: the repeat gates on the action completing, not only on Cond.
type RepeatWhile struct {
	Cond Condition
	Do   GatingEffect
}

// Text renders the loop, leading with the condition and closing with the
// self-repeat gate.
func (e RepeatWhile) Text() string {
	return e.Cond.CondText() + ", " + e.Do.Text() + " -> repeat this effect"
}

// Resolve runs Do while Cond is met, stopping as soon as Do does nothing or the
// Rule of Six is reached.
func (e RepeatWhile) Resolve(ctx *EffectContext) {
	for range RuleOfSix {
		if !e.Cond.Met(ctx) || !e.Do.resolveGate(ctx) {
			return
		}
	}
}

// validate checks the looped effect for configuration errors.
func (e RepeatWhile) validate() error {
	return validateEffect(e.Do)
}

// Overwhelmed reports whether the controller is overwhelmed — their opponent
// controls more creatures than they do. "Overwhelmed" is the keyword form of that
// board state; Numquid the Fair repeats its destruction while overwhelmed.
type Overwhelmed struct{}

// CondText renders the condition.
func (Overwhelmed) CondText() string { return "if you are overwhelmed" }

// Met reports whether the opponent controls more creatures than the controller.
func (Overwhelmed) Met(ctx *EffectContext) bool {
	return len(
		ctx.Resolver.Battleline(ctx.Opponent()),
	) > len(
		ctx.Resolver.Battleline(ctx.Controller),
	)
}

// RepeatOnCondition performs a gating effect and repeats it while the effect keeps
// succeeding and a condition holds — Numquid the Fair's "destroy an enemy creature
// -> if you are overwhelmed, repeat this effect." Do runs at least once; the loop
// stops as soon as Do does nothing (its gate is false) or Cond is not met. Do is a
// GatingEffect so an empty board ends the loop instead of spinning.
type RepeatOnCondition struct {
	Do   GatingEffect
	Cond Condition
}

// validate checks the repeated effect.
func (e RepeatOnCondition) validate() error {
	return validateEffect(e.Do)
}

// Text renders the effect, e.g. "destroy an enemy creature -> if you are
// overwhelmed, repeat this effect".
func (e RepeatOnCondition) Text() string {
	return e.Do.Text() + " -> " + e.Cond.CondText() + ", repeat this effect"
}

// Resolve runs Do, repeating while it keeps doing something, Cond holds, and the
// Rule of Six allows another pass.
func (e RepeatOnCondition) Resolve(ctx *EffectContext) {
	for range RuleOfSix {
		if !e.Do.resolveGate(ctx) || !e.Cond.Met(ctx) {
			return
		}
	}
}

// MayRepeat resolves Do once, then offers the controller the choice to resolve it
// again for as long as Cond holds and they keep accepting — the optional
// counterpart to RepeatWhile, modelling "<do>. If <cond>, you may repeat this
// effect." Do should make progress toward failing Cond so the loop can end.
type MayRepeat struct {
	Cond Condition
	Do   Effect
}

// Text renders the effect, closing with the optional self-repeat gate.
func (e MayRepeat) Text() string {
	return e.Do.Text() + " -> " + e.Cond.CondText() + ", you may repeat this effect"
}

// Resolve runs Do once, then repeats it while Cond holds, the controller keeps
// choosing to repeat, and the Rule of Six allows another pass.
func (e MayRepeat) Resolve(ctx *EffectContext) {
	e.Do.Resolve(ctx)
	for range RuleOfSix - 1 {
		if !e.Cond.Met(ctx) ||
			ctx.ChooseOption("Repeat this effect?", []string{"Yes", "No"}) != 0 {
			return
		}
		e.Do.Resolve(ctx)
	}
}

// validate checks the repeated effect for configuration errors.
func (e MayRepeat) validate() error {
	return validateEffect(e.Do)
}

// ChoseHouse is met when the controller's active house is House. It is the
// condition behind an "After you choose <House> as your active house, ..."
// ability (Jehu the Bureaucrat): the AfterChooseHouse trigger fires for the
// active player as they pick their house, and this checks whether they picked
// the house the ability watches for.
type ChoseHouse struct {
	House House
}

// CondText renders the condition clause.
func (c ChoseHouse) CondText() string {
	return "you choose " + c.House.String() + " as your active house"
}

// Met reports whether the active house is the one the ability watches for.
func (c ChoseHouse) Met(ctx *EffectContext) bool {
	return ctx.Resolver.ActiveHouse() == c.House
}

// CountIs gates on any Count: Is names the comparison and Amount the threshold
// it compares the count's value against. It is the general "if <something>
// happened N times" condition, so a card that checks a quantity reuses the Count
// vocabulary instead of adding a bespoke condition — Stampede checks the
// creatures used this turn, Vigor checks the damage it just healed.
type CountIs struct {
	Count  Count
	Is     Comparison
	Amount int
}

// validate requires a Count that can render its own "if" clause and a comparison
// the count can answer (MoreThanYou compares two Æmber pools and means nothing
// here).
func (c CountIs) validate() error {
	if _, ok := c.Count.(countClauser); !ok {
		return fmt.Errorf("CountIs: Count must be set and render a CountClause")
	}
	switch c.Is {
	case AtLeast, Exactly:
		return nil
	default:
		return fmt.Errorf("CountIs: Is must be AtLeast or Exactly")
	}
}

// CondText renders the condition, e.g. "if you used 3 or more creatures this
// turn", asking the Count for the clause that reads naturally after "if".
func (c CountIs) CondText() string {
	quantity := fmt.Sprintf("%d or more", c.Amount)
	if c.Is == Exactly {
		quantity = fmt.Sprintf("exactly %d", c.Amount)
	}
	return "if " + c.Count.(countClauser).CountClause(quantity, c.Amount != 1 || c.Is != Exactly)
}

// Met compares the count's current value against the threshold.
func (c CountIs) Met(ctx *EffectContext) bool {
	if c.Is == Exactly {
		return c.Count.Value(ctx) == c.Amount
	}
	return c.Count.Value(ctx) >= c.Amount
}

// countClauser is the optional capability a Count implements to render the "if
// ..." clause CountIs needs. A Count's CountText is a noun ("card destroyed this
// way") that reads well after "for each" but not after "if"; a Count that has a
// natural verb phrase supplies it here, given the rendered quantity ("3 or more",
// "exactly 1") and whether that quantity takes a plural noun.
type countClauser interface {
	CountClause(quantity string, plural bool) string
}

// ForgedKey is the condition on whether a player forged a key in a given window —
// this turn (Smiling Ruth) or on their own previous turn (Tendrils of Pain, Key
// Hammer). It reads the turn history rather than the running key total, so a key
// forged several turns ago does not keep the condition true.
type ForgedKey struct {
	Player   Player
	Previous bool
}

// validate requires the condition to name whose key it asks about.
func (c ForgedKey) validate() error {
	if !c.Player.valid() {
		return fmt.Errorf("ForgedKey: Player must be set")
	}
	return nil
}

// CondText renders the clause, e.g. "if your opponent forged a key on their
// previous turn".
func (c ForgedKey) CondText() string {
	subject, possessive := "you", "your"
	if c.Player == Opponent {
		subject, possessive = "your opponent", "their"
	}
	when := "this turn"
	if c.Previous {
		when = "on " + possessive + " previous turn"
	}
	return fmt.Sprintf("if %s forged a key %s", subject, when)
}

// Met reports whether the named player forged at least one key in the window.
func (c ForgedKey) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.PlayerFor(c.Player), c.stat()) > 0
}

// stat picks the tally the window corresponds to.
func (c ForgedKey) stat() TurnStat {
	if c.Previous {
		return KeysForgedLastTurn
	}
	return KeysForgedThisTurn
}
