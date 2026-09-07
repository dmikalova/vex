package engine

import (
	"errors"
	"fmt"
)

// To exalt a creature is to place 1 Æmber from the common supply onto a chosen
// friendly or enemy creature. The Æmber sits on the creature, belonging to no
// pool, until it leaves play, then goes to the owner's opponent's pool. Exalting
// N times places N Æmber.
type Exalt struct {
	Target Target
	Amount int
}

// validate requires an explicit target.
func (e Exalt) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Exalt")
	}
	return nil
}

// Text renders the effect, e.g. "exalt an enemy creature 2 times". A single
// exalt drops the count so it reads naturally.
func (e Exalt) Text() string {
	if e.Amount == 1 {
		return "exalt " + e.Target.Text()
	}
	return fmt.Sprintf("exalt %s %d times", e.Target.Text(), e.Amount)
}

// Resolve chooses a creature (through the Target) and places Amount Æmber on it.
func (e Exalt) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddAmberOn(id, e.Amount)
		ctx.Resolver.Record(AemberExalted{Creature: id, Amount: e.Amount})
	}
}

// ExaltToRepeat resolves Do once, then lets the controller exalt a creature to
// repeat it, for as long as they keep paying — the exalt is the cost of another
// pass, bounded by the Rule of Six. It models "<do>. You may exalt <a creature>
// to repeat the preceding effect." (Phalanx Strike).
type ExaltToRepeat struct {
	// Do is the preceding effect each exalt repeats.
	Do Effect
	// Exalt names the creature exalted to pay for a repeat.
	Exalt Target
}

// validate requires an effect to repeat and an explicit exalt target.
func (e ExaltToRepeat) validate() error {
	if e.Do == nil {
		return errors.New("ExaltToRepeat needs an effect to Do")
	}
	if !e.Exalt.valid() {
		return errUnsetTarget("ExaltToRepeat")
	}
	return validateEffect(e.Do)
}

// Text renders Do as its own sentence, then the optional exalt-to-repeat gate.
func (e ExaltToRepeat) Text() string {
	return punctuate(e.Do.Text()) +
		" You may exalt " + e.Exalt.Text() + " to repeat the preceding effect"
}

// Resolve runs Do once, then repeats it each time the controller exalts a chosen
// creature (a declinable choice), up to the Rule of Six. Declining the exalt ends
// the loop.
func (e ExaltToRepeat) Resolve(ctx *EffectContext) {
	e.Do.Resolve(ctx)
	for range RuleOfSix - 1 {
		ids := e.Exalt.SelectOptional(ctx)
		if len(ids) == 0 {
			return
		}
		for _, id := range ids {
			ctx.Resolver.AddAmberOn(id, 1)
			ctx.Resolver.Record(AemberExalted{Creature: id, Amount: 1})
		}
		e.Do.Resolve(ctx)
	}
}
