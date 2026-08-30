package engine

import "fmt"

// A Conditional gates an effect behind a check on the current game state — the
// "If ..." clause many cards open with, e.g. "If your opponent has 7 or more
// Æmber, they lose 4 Æmber." The Condition renders its own English (CondText) and
// evaluates itself (Met), so the printed text always matches what is checked.

// Condition is a boolean predicate on the live game, used by Conditional.
type Condition interface {
	CondText() string
	Met(ctx *EffectContext) bool
}

// OpponentAemberAtLeast is met when the opponent's pool holds at least Amount
// Æmber.
type OpponentAemberAtLeast struct {
	Amount int
}

// CondText renders the condition, e.g. "if your opponent has 7 Æmber or more".
func (c OpponentAemberAtLeast) CondText() string {
	return fmt.Sprintf("if your opponent has %d Æmber or more", c.Amount)
}

// Met reports whether the opponent has at least Amount Æmber.
func (c OpponentAemberAtLeast) Met(ctx *EffectContext) bool {
	return ctx.Resolver.Aember(ctx.Opponent()) >= c.Amount
}

// OpponentAemberExactly is met when the opponent's pool holds exactly Amount
// Æmber.
type OpponentAemberExactly struct {
	Amount int
}

// CondText renders the condition, e.g. "if your opponent has exactly 1 Æmber".
func (c OpponentAemberExactly) CondText() string {
	return fmt.Sprintf("if your opponent has exactly %d Æmber", c.Amount)
}

// Met reports whether the opponent has exactly Amount Æmber.
func (c OpponentAemberExactly) Met(ctx *EffectContext) bool {
	return ctx.Resolver.Aember(ctx.Opponent()) == c.Amount
}

// OpponentAemberMoreThanYou is met while the opponent's pool holds strictly more
// Æmber than the controller's.
type OpponentAemberMoreThanYou struct{}

// CondText renders the condition.
func (OpponentAemberMoreThanYou) CondText() string {
	return "if your opponent has more Æmber than you"
}

// Met reports whether the opponent has more Æmber than the controller.
func (OpponentAemberMoreThanYou) Met(ctx *EffectContext) bool {
	return ctx.Resolver.Aember(ctx.Opponent()) > ctx.Resolver.Aember(ctx.Controller)
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
	return ctx.Produced.Destroyed < c.Amount
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
	if e.Type != "" && ctx.Resolver.TypeOf(ctx.It) != e.Type {
		return false
	}
	return true
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
type Conditional struct {
	Cond Condition
	Then Effect
}

// Text joins the condition and the gated effect.
func (e Conditional) Text() string { return e.Cond.CondText() + ", " + e.Then.Text() }

// Resolve runs Then only if Cond is met.
func (e Conditional) Resolve(ctx *EffectContext) {
	if e.Cond.Met(ctx) {
		e.Then.Resolve(ctx)
	}
}

// validate checks the gated effect for configuration errors.
func (e Conditional) validate() error {
	return validateEffect(e.Then)
}

// RepeatWhile resolves Do again and again for as long as Cond holds, re-checking
// after each pass — a self-looping effect such as "if your opponent has more
// Æmber than you, steal 1 Æmber -> repeat this effect". Do must make progress
// toward ending the loop (each pass changes the state Cond checks).
type RepeatWhile struct {
	Cond Condition
	Do   Effect
}

// Text renders the loop, leading with the condition and closing with the
// self-repeat gate.
func (e RepeatWhile) Text() string {
	return e.Cond.CondText() + ", " + e.Do.Text() + " -> repeat this effect"
}

// Resolve runs Do while Cond is met.
func (e RepeatWhile) Resolve(ctx *EffectContext) {
	for e.Cond.Met(ctx) {
		e.Do.Resolve(ctx)
	}
}

// validate checks the looped effect for configuration errors.
func (e RepeatWhile) validate() error {
	return validateEffect(e.Do)
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

// Resolve runs Do once, then repeats it while Cond holds and the controller keeps
// choosing to repeat.
func (e MayRepeat) Resolve(ctx *EffectContext) {
	e.Do.Resolve(ctx)
	for e.Cond.Met(ctx) {
		if ctx.ChooseOption("Repeat this effect?", []string{"Yes", "No"}) != 0 {
			return
		}
		e.Do.Resolve(ctx)
	}
}

// validate checks the repeated effect for configuration errors.
func (e MayRepeat) validate() error {
	return validateEffect(e.Do)
}
