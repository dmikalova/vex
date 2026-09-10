package engine

// ReturnItToHand puts the creature in context ("it") into its owner's hand,
// recovering it from the discard pile when it has already been destroyed. It is
// the reaction counterpart to PutFromPlay, which only reaches a creature still in
// play: Nizak, The Forgotten returns an enemy destroyed fighting it, and that
// creature already sits in the discard by the time the reaction resolves.
type ReturnItToHand struct{}

// Text renders the effect, e.g. "put it into its owner's hand".
func (e ReturnItToHand) Text() string {
	return ToHand.clause(Target{Kind: TargetTriggeringCreature}.Text(), false)
}

// Resolve moves the contextual creature to its owner's hand — from play if it is
// still there, or recovered from the discard pile if it was already destroyed. It
// does nothing when no card is in context.
func (e ReturnItToHand) Resolve(ctx *EffectContext) {
	if !ctx.HasIt {
		return
	}
	if resolverInPlay(ctx, ctx.It) {
		ctx.Resolver.PutIntoHand(ctx.It)
		return
	}
	ctx.Resolver.PutFromDiscardIntoHand(ctx.It)
}
