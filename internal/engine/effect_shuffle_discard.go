package engine

// ShuffleMatchingFromDiscardIntoDeck shuffles every card in the controller's
// discard pile matching House and Type back into their deck — Low Dawn returns
// each Untamed creature from the discard to the deck. Unlike
// ShuffleChosenCreaturesFromDiscard it makes no choice: every match moves. An
// unset House or Type applies no filter on that axis.
type ShuffleMatchingFromDiscardIntoDeck struct {
	House House
	Type  CardType
}

// Text renders the effect, e.g. "shuffle each Untamed creature from your discard
// pile into your deck".
func (e ShuffleMatchingFromDiscardIntoDeck) Text() string {
	return "shuffle each " + houseTypeNoun(e.House, e.Type) +
		" from your discard pile into your deck"
}

// Resolve moves each matching discard-pile card into the deck as one grouped
// shuffle batch, so the several cards narrate one line.
func (e ShuffleMatchingFromDiscardIntoDeck) Resolve(ctx *EffectContext) {
	cards := discardCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
			return false
		}
		return e.Type == TypeUnset || ctx.Resolver.TypeOf(id) == e.Type
	})
	if len(cards) == 0 {
		return
	}
	ctx.Resolver.BeginShuffleBatch()
	for _, id := range cards {
		ctx.Resolver.ShuffleFromDiscardIntoDeck(id)
	}
	ctx.Resolver.EndShuffleBatch(ctx.Source)
}
