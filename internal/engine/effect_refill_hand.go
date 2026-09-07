package engine

// EachPlayerDiscardsAndRefillsHand makes both players discard their whole hand
// and then draw a fresh one as if it were the end of their turn — honoring each
// player's chains and draw modifiers. It is Punctuated Equilibrium: "Each player
// discards their hand, then refills their hand as if it were the end of their
// turn." The controller resolves first.
type EachPlayerDiscardsAndRefillsHand struct{}

// validate has no required fields.
func (e EachPlayerDiscardsAndRefillsHand) validate() error { return nil }

// Text renders the printed sentence.
func (e EachPlayerDiscardsAndRefillsHand) Text() string {
	return "each player discards their hand, then refills their hand as if it " +
		"were the end of their turn"
}

// Resolve discards each player's hand, then refills each — the controller first,
// then the opponent — so every discard happens before any refill.
func (e EachPlayerDiscardsAndRefillsHand) Resolve(ctx *EffectContext) {
	players := []int{ctx.Controller, ctx.Opponent()}
	for _, p := range players {
		for _, id := range ctx.Resolver.Hand(p) {
			ctx.Resolver.DiscardCardFromHand(p, id)
		}
	}
	for _, p := range players {
		ctx.Resolver.RefillHand(p)
	}
}
