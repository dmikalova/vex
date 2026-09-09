package engine

// SaveFromDestruction is a creature's own "Destroyed:" replacement: instead of
// being destroyed, the creature stays in play and Do resolves on it (Reassembling
// Automaton fully heals, exhausts, and moves itself to a flank). It marks itself
// saved so the destruction batch's discard step leaves it in play, then runs Do
// with the creature in context ("it"). Do must keep the creature out of a
// destroyable state — fully healing it — or the state-based sweep destroys it
// again.
type SaveFromDestruction struct {
	Do Effect
}

// validate requires a well-formed replacement effect.
func (e SaveFromDestruction) validate() error { return validateEffect(e.Do) }

// Text renders the effect, e.g. "instead of destroying {self}, fully heal it,
// exhaust it, and move it to a flank".
func (e SaveFromDestruction) Text() string {
	return "instead of destroying " + SelfName + ", " + e.Do.Text()
}

// Resolve marks the source saved from the current destruction batch, then runs Do
// with the source as the creature in context ("it").
func (e SaveFromDestruction) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate performs the replacement and reports true, so a Then can hang a
// follow-up on the save ("If you do, give it two +1 power counters"). It always
// succeeds once reached — the caller gates whether it is reached (Self-Bolstering
// Automata's "if you have any other creatures in play").
func (e SaveFromDestruction) resolveGate(ctx *EffectContext) bool {
	ctx.Resolver.SaveFromDestruction(ctx.Source)
	ctx.It = ctx.Source
	ctx.HasIt = true
	e.Do.Resolve(ctx)
	return true
}
