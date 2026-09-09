package engine

// Shared candidate-gather helpers for the effects that put cards between zones
// (ADR 0031): an archive, discard, purge, shuffle, or put first selects cards from
// a source zone, and these return that zone's cards matching a per-card filter, in
// zone order.

// handCardsWhere returns the cards in a player's hand that satisfy keep, in hand
// order — the candidate set a hand effect (purge, archive, discard) selects from.
func handCardsWhere(ctx *EffectContext, owner int, keep func(LocalID) bool) []LocalID {
	var out []LocalID
	for _, id := range ctx.Resolver.Hand(owner) {
		if keep(id) {
			out = append(out, id)
		}
	}
	return out
}

// discardCardsWhere returns the cards in a player's discard pile that satisfy keep,
// in pile order — the candidate set a discard-pile effect (shuffle, archive, purge)
// selects from.
func discardCardsWhere(ctx *EffectContext, owner int, keep func(LocalID) bool) []LocalID {
	var out []LocalID
	for _, id := range ctx.Resolver.Discard(owner) {
		if keep(id) {
			out = append(out, id)
		}
	}
	return out
}
