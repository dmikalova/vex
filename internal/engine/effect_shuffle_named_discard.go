package engine

// ShuffleNamedFromDiscardIntoDeck shuffles one card of the given printed name from
// the controller's discard pile back into their deck — Chain Gang returns a Subtle
// Chain from the discard so it can be played again. Where a ShuffleFromDiscard with
// an Each selection moves every House/Type match, this moves a single named card;
// it stays its own node because its printed article ("a Subtle Chain") belongs to
// the verb rather than to the Selection (a Named renders bare).
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
	name := e.Name
	cards := Named{Name: name}.pick(ctx, ctx.Resolver.Discard(ctx.Controller))
	if len(cards) == 0 {
		return
	}
	ctx.Resolver.BeginShuffleBatch()
	ctx.Resolver.ShuffleFromDiscardIntoDeck(cards[0])
	ctx.Resolver.EndShuffleBatch()
}
