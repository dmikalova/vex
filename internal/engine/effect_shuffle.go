package engine

// ShuffleDiscard shuffles the controller's discard pile into their deck (Help from
// Future Self, after tutoring).
//
//rulebook:effect Shuffle Discard
type ShuffleDiscard struct{}

// Text renders the effect, "shuffle your discard pile into your deck".
func (ShuffleDiscard) Text() string {
	return "shuffle your discard pile into your deck"
}

// Resolve moves the controller's discard pile into their deck and shuffles it.
func (ShuffleDiscard) Resolve(ctx *EffectContext) {
	ctx.Resolver.ShuffleDiscardIntoDeck(ctx.Controller)
}

// ShuffleHandAndDiscard shuffles both the controller's hand and their discard pile
// into their deck (Screaming Cave).
type ShuffleHandAndDiscard struct{}

// Text renders the effect, "shuffle your hand and discard pile into your deck".
func (ShuffleHandAndDiscard) Text() string {
	return "shuffle your hand and discard pile into your deck"
}

// Resolve moves the controller's hand and discard pile into their deck and shuffles.
func (ShuffleHandAndDiscard) Resolve(ctx *EffectContext) {
	ctx.Resolver.ShuffleHandAndDiscardIntoDeck(ctx.Controller)
}
