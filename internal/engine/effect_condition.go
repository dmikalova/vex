package engine

import (
	"fmt"
	"strings"
)

// A conditional gates an effect behind a check on the current game state — the
// "If ..." clause a card opens with, e.g. "If your opponent has 7 or more Æmber,
// they lose 4 Æmber." The effect resolves only when the condition is met. Unlike
// a result gate (A -> B), which turns on an action succeeding, a conditional turns
// on a fact about the board.
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
	// comparisonUnset is the invalid zero value; a condition must name a real
	// comparison, so an unset one is caught at init.
	comparisonUnset Comparison = iota
	// AtLeast is met when the quantity is at least Amount.
	AtLeast
	// AtMost is met when the quantity is at most Amount.
	AtMost
	// Exactly is met when the quantity is exactly Amount (which is why the
	// comparison is named separately from the amount: Exactly with Amount 0 is a
	// real check that a bare integer field could not tell from "unset").
	Exactly
	// MoreThanYou is met when the opponent's pool holds strictly more Æmber than the
	// controller's; it ignores Amount and applies only to PoolAember{Player: Opponent}.
	MoreThanYou
	// MoreThanOpponent is met when the controller's pool holds strictly more Æmber
	// than the opponent's; it ignores Amount and applies only to
	// PoolAember{Player: Controller}.
	MoreThanOpponent
	// Even is met when the quantity is even (including zero); it ignores Amount —
	// Even Ivan steals while the opponent's pool is even.
	Even
	// Odd is met when the quantity is odd; it ignores Amount — Odd Clawde steals
	// while the opponent's pool is odd.
	Odd
)

// validateCondition returns any configuration error a condition reports (an unset
// PoolAember comparison, say). Conditions that cannot be misconfigured
// implement no validator and pass.
func validateCondition(c Condition) error {
	if v, ok := c.(validator); ok {
		return v.validate()
	}
	return nil
}

// HouseChoice names one house by reference — the house chosen this turn, the
// active house, the house of the card in context, or no house at all — for an
// effect that must resolve it to a single concrete house: ChangeActiveHouse
// switches the active house to it, CardsInHand counts a hand by it. It is the
// deliberate sibling of HouseMatcher (ADR 0038), which is a predicate over a set
// of houses (a named house, a non-<house>): a HouseChoice resolves to one house,
// so its Except/named-set notions would be meaningless and it stays separate.
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

// phrase renders the trailing "of the … house" fragment an effect welds onto a
// noun ("creatures of the chosen house"). The no-scope choices (unset, AnyHouse)
// render nothing, so a caller reads an empty phrase as "no house scope".
func (h HouseChoice) phrase() string {
	switch h {
	case TheChosenHouse:
		return "of the chosen house"
	case TheActiveHouse:
		return "of the active house"
	case TheContextualHouse:
		return "of that card's house"
	default:
		return ""
	}
}

// typeNoun renders the bare card-type noun — "creature", "artifact", or "card"
// when the type is unset — without any house qualifier.
func typeNoun(typ CardType) string {
	switch typ {
	case Creature:
		return "creature"
	case Artifact:
		return "artifact"
	case Tactic:
		return "tactic"
	case Upgrade:
		return "upgrade"
	}
	return "card"
}

// Or is met when any one of its Conditions is met, composing conditions instead of
// baking each combination into a bespoke one — Guji Dinosaur Hunter boosts against
// a Dinosaur creature or a creature with Æmber on it.
type Or struct {
	Conditions []Condition
}

// validate requires at least two conditions and rejects any invalid one.
func (o Or) validate() error {
	if len(o.Conditions) < 2 {
		return fmt.Errorf("Or: needs at least two conditions")
	}
	for _, c := range o.Conditions {
		if err := validateCondition(c); err != nil {
			return err
		}
	}
	return nil
}

// AlwaysMet is a Condition that is always met, rendering no clause of its own. It
// is the non-nil sentinel a nil-means-absent condition field uses to say "on,
// unconditionally" — WithAemberCannotBeStolen() sets it so the pool is protected
// with no "while …" qualifier.
type AlwaysMet struct{}

// CondText renders nothing: an always-true condition adds no "while …" clause.
func (AlwaysMet) CondText() string { return "" }

// Met is always true.
func (AlwaysMet) Met(*EffectContext) bool { return true }

// negatable is a Condition that can render its own negation — the clause a Not// wrapper prints. Negated wording is not uniform across conditions (a flank reads
// "is not on a flank", a threshold flips to "fewer than"), so each negatable
// condition owns its negated text rather than a wrapper deriving it by string
// surgery.
type negatable interface {
	Condition
	negatedText() string
}

// Not is met when its Condition is not met, composing negation instead of a
// per-condition Not flag. Its Condition must be negatable (render its own negated
// clause): Not{OnFlank{}} reads "if it is not on a flank", Not{ForgedKey{...}}
// "if you have not forged a key this turn".
type Not struct {
	Cond Condition
}

// validate rejects a Condition that cannot render its own negation, and rejects an
// invalid inner condition.
func (n Not) validate() error {
	if _, ok := n.Cond.(negatable); !ok {
		return fmt.Errorf("Not: %T cannot be negated", n.Cond)
	}
	if c, ok := n.Cond.(CountIs); ok && c.Is != AtLeast {
		return fmt.Errorf("Not: CountIs must be AtLeast to negate")
	}
	return validateCondition(n.Cond)
}

// CondText renders the inner condition's negated clause.
func (n Not) CondText() string { return n.Cond.(negatable).negatedText() }

// Met reports whether the inner condition is not met.
func (n Not) Met(ctx *EffectContext) bool { return !n.Cond.Met(ctx) }

// CondText joins the sub-clauses with "or", e.g. "if it is a Dinosaur creature or
// it has Æmber on it". Each condition renders "if <clause>" (the shared
// convention), so the leading "if " is dropped before the clauses are joined.
func (o Or) CondText() string {
	if s, ok := o.combinedHouses(); ok {
		return s
	}
	clauses := make([]string, len(o.Conditions))
	for i, c := range o.Conditions {
		clauses[i] = strings.TrimPrefix(c.CondText(), "if ")
	}
	return "if " + strings.Join(clauses, " or ")
}

// combinedHouses renders an Or of ItIs clauses that differ only in a single named
// house as one phrase — "if it is a Dis or Shadows card" — rather than repeating
// "it is" once per house (Ambassador Liu). It returns false unless every clause is
// such an ItIs and they share the same type, subject, and other flag.
func (o Or) combinedHouses() (string, bool) {
	houses := make([]string, 0, len(o.Conditions))
	var shape ItIs
	for i, c := range o.Conditions {
		it, ok := c.(ItIs)
		if !ok {
			return "", false
		}
		h, sh, ok := it.asNamedHouseAlt()
		if !ok {
			return "", false
		}
		if i == 0 {
			shape = sh
		} else if sh != shape {
			return "", false
		}
		houses = append(houses, h.String())
	}
	noun := strings.Join(houses, " or ") + " " + typeNoun(shape.Type)
	return "if " + shape.Noun.noun() + " is " + indefinite(noun), true
}

// Met reports whether any of the conditions is met.
func (o Or) Met(ctx *EffectContext) bool {
	for _, c := range o.Conditions {
		if c.Met(ctx) {
			return true
		}
	}
	return false
}

// And is met when every one of its Conditions is met, composing conditions
// instead of baking each combination into a bespoke one — Dark Æmber Vault draws
// when a creature you played is a Mutant (it is friendly and it is a Mutant).
type And struct {
	Conditions []Condition
}

// validate requires at least two conditions and rejects any invalid one.
func (a And) validate() error {
	if len(a.Conditions) < 2 {
		return fmt.Errorf("And: needs at least two conditions")
	}
	for _, c := range a.Conditions {
		if err := validateCondition(c); err != nil {
			return err
		}
	}
	return nil
}

// CondText joins the sub-clauses with "and", e.g. "if it is a friendly creature
// and it is used to reap". Each condition renders "if <clause>" (the shared
// convention), so the leading "if " is dropped before the clauses are joined.
// Clauses that all describe the shape of the card in context collapse into one
// noun phrase first — see collapsedItShape.
func (a And) CondText() string {
	if phrase, ok := collapsedItShape(a.Conditions); ok {
		return phrase
	}
	clauses := make([]string, len(a.Conditions))
	for i, c := range a.Conditions {
		clauses[i] = strings.TrimPrefix(c.CondText(), "if ")
	}
	return "if " + strings.Join(clauses, " and ")
}

// itShaped is a condition that constrains the shape of the card in context and
// can render as an adjective on a shared noun rather than as a whole clause.
type itShaped interface {
	// itAdjective is the word the clause contributes before the noun — "friendly",
	// "Cat", "Mars". An empty string declines the collapse.
	itAdjective() string
	// itNoun is the noun the clause describes. Clauses collapse only when they
	// agree on it.
	itNoun() string
}

// collapsedItShape renders conditions that all describe the shape of the card in
// context as a single noun phrase — "if it is a friendly Cat creature" (Mercy,
// Malkin Queen) rather than "if it is a friendly creature and it is a Cat
// creature". It declines unless every condition is an itShaped that offers an
// adjective and they agree on the noun. Pinned by TestAndCollapsesItShapeClauses.
func collapsedItShape(conds []Condition) (string, bool) {
	adjectives := make([]string, 0, len(conds))
	noun := ""
	for _, c := range conds {
		shaped, ok := c.(itShaped)
		if !ok {
			return "", false
		}
		adj := shaped.itAdjective()
		if adj == "" || (noun != "" && shaped.itNoun() != noun) {
			return "", false
		}
		noun = shaped.itNoun()
		adjectives = append(adjectives, adj)
	}
	return "if it is " + indefinite(strings.Join(adjectives, " ")+" "+noun), true
}

// Met reports whether every condition is met.
func (a And) Met(ctx *EffectContext) bool {
	for _, c := range a.Conditions {
		if !c.Met(ctx) {
			return false
		}
	}
	return true
}

// Conditional resolves Then only when Cond is met. It renders as "<cond>, <then>",
// e.g. "if your opponent has 7 Æmber or more, your opponent loses 4 Æmber".
//
// Else, when set, resolves when Cond is not met and renders as the second
// sentence the cards use for a two-way branch: "<cond>, <then>. Otherwise,
// <else>." (Vespilon Theorist archives the revealed card or discards it).
type Conditional struct {
	Cond Condition
	Then Effect
	Else Effect
}

// Text joins the condition and the gated effect.
func (e Conditional) Text() string {
	body := e.Cond.CondText() + ", " + e.Then.Text()
	if e.Else == nil {
		return body
	}
	return body + ". Otherwise, " + e.Else.Text()
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

// OrAmount is an alternate scalar an effect switches to when When holds, so a card
// that only changes a number by a fact prints the linear "<base> …, or <alt> if
// <cond>" form instead of a two-armed "If <cond>, … Otherwise, …" branch
// (card-wording rule 22). The effect owns how the numbers read — a bare "2" for a
// steal, a "+2" surcharge for a forge — so OrAmount supplies only the alternate
// value, the guard, and the shared ", or … if <cond>" tail. Its zero value (When
// nil) means "no alternate", so an effect treats an unset OrAmount as absent.
type OrAmount struct {
	Amount int
	When   Condition
}

// set reports whether an alternate is configured.
func (o OrAmount) set() bool { return o.When != nil }

// pick returns the alternate amount when the guard holds, else base.
func (o OrAmount) pick(base int, ctx *EffectContext) int {
	if o.set() && o.When.Met(ctx) {
		return o.Amount
	}
	return base
}

// tail renders ", or <alt> if <cond>", with alt already formatted by the effect.
func (o OrAmount) tail(alt string) string {
	return ", or " + alt + " " + o.When.CondText()
}

// validate rejects an alternate whose guard is missing or invalid.
func (o OrAmount) validate() error {
	if !o.set() {
		return nil
	}
	return validateCondition(o.When)
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
	case AtLeast, AtMost, Exactly, Even, Odd:
		return nil
	default:
		return fmt.Errorf("CountIs: Is must be AtLeast, AtMost, Exactly, Even, or Odd")
	}
}

// CondText renders the condition, e.g. "if you used 3 or more creatures this
// turn", asking the Count for the clause that reads naturally after "if". Parity
// says "number" where a threshold says a figure, because a Count counts things;
// the mass-noun Æmber pool says "amount" on PoolAember for the same reason.
func (c CountIs) CondText() string {
	var quantity string
	switch c.Is {
	case AtMost:
		quantity = fmt.Sprintf("%d or fewer", c.Amount)
	case Exactly:
		quantity = fmt.Sprintf("exactly %d", c.Amount)
	case Even:
		quantity = "an even number of"
	case Odd:
		quantity = "an odd number of"
	default:
		quantity = fmt.Sprintf("%d or more", c.Amount)
	}
	return "if " + c.Count.(countClauser).CountClause(quantity, c.Amount != 1 || c.Is != Exactly)
}

// Met compares the count's current value against the threshold. Even and Odd
// ignore Amount: parity is a property of the value, not a comparison to a figure.
func (c CountIs) Met(ctx *EffectContext) bool {
	switch c.Is {
	case AtMost:
		return c.Count.Value(ctx) <= c.Amount
	case Exactly:
		return c.Count.Value(ctx) == c.Amount
	case Even:
		return c.Count.Value(ctx)%2 == 0
	case Odd:
		return c.Count.Value(ctx)%2 != 0
	default:
		return c.Count.Value(ctx) >= c.Amount
	}
}

// negatedText renders the "fewer than" clause a Not wrapper prints — "at least N"
// flips to "fewer than N". Only an AtLeast CountIs negates; Not.validate holds the
// guard against negating an Exactly count.
func (c CountIs) negatedText() string {
	quantity := fmt.Sprintf("fewer than %d", c.Amount)
	return "if " + c.Count.(countClauser).CountClause(quantity, c.Amount != 1)
}

// countClauser is the optional capability a Count implements to render the "if
// ..." clause CountIs needs. A Count's CountText is a noun ("card destroyed this
// way") that reads well after "for each" but not after "if"; a Count that has a
// natural verb phrase supplies it here, given the rendered quantity ("3 or more",
// "exactly 1") and whether that quantity takes a plural noun.
type countClauser interface {
	CountClause(quantity string, plural bool) string
}
