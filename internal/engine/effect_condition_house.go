package engine

// ItIsNotOfHouse is met when the card in context (ctx.It — a revealed, discarded,
// or triggering card) does NOT belong to a named house: the "non-<house> card"
// idiom (Book of leQ acts on a non-Star Alliance revealed card). With no card in
// context the condition is not met, so it holds only when a card is present and
// its house differs from House.
type ItIsNotOfHouse struct {
	House House
}

// CondText renders the condition, e.g. "if it is a non-Star Alliance card".
func (e ItIsNotOfHouse) CondText() string {
	return "if it is a non-" + e.House.String() + " card"
}

// Met reports whether a card is in context and does not belong to the named house.
func (e ItIsNotOfHouse) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.House(ctx.It) != e.House
}
