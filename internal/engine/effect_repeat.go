package engine

import "errors"

// Repeat resolves an effect once for each of a running count, so every repetition
// makes its own choices — Mothership Support deals 2 damage per friendly ready
// Mars creature and may pick a different creature each time. It is the
// choose-again counterpart to a Per clause, which multiplies one effect's amount
// against a single target.
type Repeat struct {
	// Times is how many repetitions to run.
	Times Count
	// Do is the effect resolved once per repetition.
	Do Effect
}

// Text renders the repetition as a leading "for each" clause over the effect's
// own phrase (rule 9).
func (e Repeat) Text() string { return forEach(e.Times, e.Do.Text()) }

// validate requires both halves of the repetition.
func (e Repeat) validate() error {
	if e.Times == nil {
		return errors.New("Repeat needs a Times count")
	}
	if e.Do == nil {
		return errors.New("Repeat needs an effect to Do")
	}
	return validateEffect(e.Do)
}

// Resolve runs the effect once per repetition, counting first so an effect that
// changes the board does not change how many times it runs.
func (e Repeat) Resolve(ctx *EffectContext) {
	for range e.Times.Value(ctx) {
		e.Do.Resolve(ctx)
	}
}
