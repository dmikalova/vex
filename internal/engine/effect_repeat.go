package engine

import "errors"

// ForEach resolves an effect once for each of a running count, so every repetition
// makes its own choices — Mothership Support deals 2 damage per friendly ready
// Mars creature and may pick a different creature each time. It is the
// choose-again counterpart to a Per clause, which multiplies one effect's amount
// against a single target.
type ForEach struct {
	// Times is how many repetitions to run.
	Times Count
	// Do is the effect resolved once per repetition.
	Do Effect
}

// Text renders the repetition as a leading "for each" clause over the effect's
// own phrase (rule 9).
func (e ForEach) Text() string { return forEach(e.Times, e.Do.Text()) }

// validate requires both halves of the repetition.
func (e ForEach) validate() error {
	if e.Times == nil {
		return errors.New("ForEach needs a Times count")
	}
	if e.Do == nil {
		return errors.New("ForEach needs an effect to Do")
	}
	return validateEffect(e.Do)
}

// Resolve runs the effect once per repetition, counting first so an effect that
// changes the board does not change how many times it runs.
func (e ForEach) Resolve(ctx *EffectContext) {
	for range e.Times.Value(ctx) {
		e.Do.Resolve(ctx)
	}
}

// RuleOfSix is the most times a card can be played, used, or made to resolve
// again in one turn. A self-repeating effect is bounded by it: Bait and Switch's
// "steal 1 Æmber -> repeat this effect" resolves the initial steal plus at most
// five repeats, so it steals six at most however far ahead the opponent is.
const RuleOfSix = 6

// RepeatOnCondition performs an effect and repeats it while the effect keeps
// succeeding and a condition holds — Numquid the Fair's "destroy an enemy creature
// -> if you are overwhelmed, repeat this effect." Do runs at least once; the loop
// stops as soon as Do does nothing (its gate is false) or Cond is not met. When Do
// cannot report progress the Rule of Six alone bounds the loop.
type RepeatOnCondition struct {
	Do   Effect
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
		if !resolveGateOf(ctx, e.Do) || !e.Cond.Met(ctx) {
			return
		}
	}
}

// MayRepeat resolves Do once, then offers the controller the choice to resolve it
// again for as long as Cond holds and they keep accepting — the optional
// counterpart to RepeatOnCondition, modelling "<do>. If <cond>, you may repeat
// this effect." Do should make progress toward failing Cond so the loop can end.
type MayRepeat struct {
	Cond Condition
	Do   Effect
}

// Text renders the effect, closing with the optional self-repeat gate.
func (e MayRepeat) Text() string {
	return e.Do.Text() + " -> " + e.Cond.CondText() + ", you may repeat this effect"
}

// Resolve runs Do once, then repeats it while Cond holds and the Rule of Six
// allows another pass. When Do leads with a single clickable choice, each repeat
// is driven by that choice — the controller keeps picking to repeat, or passes
// with Done — rather than a separate Yes/No question.
func (e MayRepeat) Resolve(ctx *EffectContext) {
	e.Do.Resolve(ctx)
	d, byChoice := e.Do.(declinableEffect)
	byChoice = byChoice && d.declinable()
	for range RuleOfSix - 1 {
		if !e.Cond.Met(ctx) {
			return
		}
		if byChoice {
			if !d.resolveOptional(ctx) {
				return
			}
			continue
		}
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
