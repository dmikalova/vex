package engine

// ShuffleFriendlyCardsInPlayIntoDeck shuffles every card the controller has in
// play — each creature and artifact, and every upgrade attached to them — into
// their deck, then draws a card for each card shuffled this way (Timequake). The
// two-sentence render pairs the mass shuffle with the payoff draw, so the draw
// count always matches the number of cards that just left play.
type ShuffleFriendlyCardsInPlayIntoDeck struct{}

// Text renders the effect as its two printed sentences.
func (ShuffleFriendlyCardsInPlayIntoDeck) Text() string {
	return "shuffle each friendly card in play into your deck. " +
		"Draw a card for each card shuffled into your deck this way"
}

// Resolve shuffles every friendly in-play card into the controller's deck as one
// grouped batch, then draws that many cards.
func (ShuffleFriendlyCardsInPlayIntoDeck) Resolve(ctx *EffectContext) {
	ctx.Resolver.BeginShuffleBatch()
	n := ctx.Resolver.ShuffleFriendlyCardsInPlayIntoDeck(ctx.Controller)
	ctx.Resolver.EndShuffleBatch(ctx.Source)
	ctx.Resolver.Draw(ctx.Controller, n, ctx.Source)
}
