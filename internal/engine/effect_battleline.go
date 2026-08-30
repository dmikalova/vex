package engine

// Swap exchanges this creature's battleline position with the creature its With
// target selects, then puts that creature in context (ctx.It) so a following
// effect can act on it — Transposition Sandals swaps with another friendly
// creature and then uses it. Only positions move: no damage, upgrades, status,
// control, or other card state travels between the two creatures, since battleline
// order is all that matters for flanks and neighbors. A With target that selects
// nothing leaves the battleline unchanged. With is a full Target, so a later card
// can swap with an enemy creature rather than a friendly one.
type Swap struct {
	With Target
}

// validate requires the creature to swap with.
func (e Swap) validate() error {
	if !e.With.valid() {
		return errUnsetTarget("Swap")
	}
	return nil
}

// Text renders the effect, e.g. "swap this creature with another friendly creature
// in your battleline".
func (e Swap) Text() string {
	return "swap this creature with " + e.With.Text() + " in your battleline"
}

// Resolve swaps this creature's position with the selected creature and puts that
// creature in context.
func (e Swap) Resolve(ctx *EffectContext) {
	for _, other := range e.With.Select(ctx) {
		ctx.Resolver.SwapBattlelinePositions(ctx.Source, other)
		ctx.It, ctx.HasIt = other, true
	}
}
