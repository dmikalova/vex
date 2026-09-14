package engine

// ShuffleFriendlyCardsIntoDeck folds every friendly card in play — each creature
// and artifact, and their upgrades — into its owner's deck as one shuffle batch,
// then tallies on ctx.Produced.Moved how many cards went into each owner's deck so
// a following Draw{Per: CardsShuffledIntoDeck} can draw one card per card that
// returned to the controller's own deck (Timequake). It is the one shuffle that
// reaches onto the board. A card the controller played but does not own goes to
// its owner's deck and so is tallied under that owner, not the controller.
type ShuffleFriendlyCardsIntoDeck struct{}

// Text renders the effect.
func (ShuffleFriendlyCardsIntoDeck) Text() string {
	return "shuffle each friendly card in play into your deck"
}

// Resolve shuffles the friendly cards home as one grouped batch and records, per
// owner, how many cards returned to that owner's deck on ctx.Produced.Moved.
func (ShuffleFriendlyCardsIntoDeck) Resolve(ctx *EffectContext) {
	ctx.Resolver.BeginShuffleBatch()
	moved := ctx.Resolver.ShuffleFriendlyCardsInPlayIntoDeck(ctx.Controller)
	ctx.Resolver.EndShuffleBatch()
	ctx.Produced.Moved[0] += moved[0]
	ctx.Produced.Moved[1] += moved[1]
}
