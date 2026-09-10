package engine

// ShuffleNamedFromDiscardIntoDeck shuffles one card of the given printed name from
// the controller's discard pile back into their deck — Chain Gang returns a Subtle
// Chain from the discard so it can be played again. Where
// ShuffleMatchingFromDiscardIntoDeck moves every House/Type match, this moves a
// single named card.
type ShuffleNamedFromDiscardIntoDeck struct {
	Name string
}

// Text renders the effect, e.g. "shuffle a Subtle Chain from your discard pile
// into your deck".
func (e ShuffleNamedFromDiscardIntoDeck) Text() string {
	return "shuffle a " + e.Name + " from your discard pile into your deck"
}

// Resolve moves the first discard-pile card of the given name into the deck.
func (e ShuffleNamedFromDiscardIntoDeck) Resolve(ctx *EffectContext) {
	cards := discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		return ctx.Resolver.Name(id) == e.Name
	})
	if len(cards) == 0 {
		return
	}
	ctx.Resolver.BeginShuffleBatch()
	ctx.Resolver.ShuffleFromDiscardIntoDeck(cards[0])
	ctx.Resolver.EndShuffleBatch(ctx.Source)
}
