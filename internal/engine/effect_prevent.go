package engine

import "fmt"

// PreventDamage marks each creature its Target selects immune to damage for the
// remainder of the turn — Shield of Justice protects each friendly creature. The
// mark is cleared by EndTurn.
type PreventDamage struct {
	Target Target
}

// validate requires an explicit target.
func (e PreventDamage) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PreventDamage")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, each friendly
// creature cannot be dealt damage".
func (e PreventDamage) Text() string {
	return fmt.Sprintf("for the remainder of the turn, %s cannot be dealt damage", e.Target.Text())
}

// Resolve marks each selected creature damage-immune for the rest of the turn.
func (e PreventDamage) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.PreventDamage(id)
	}
}
