package engine

// ActiveHouseMatchesNoCardsInPlay is met when no card in play — across both
// players and every card type (creatures, their upgrades, and artifacts) —
// belongs to the currently chosen active house. Sci. Officer Qincan steals only
// when a player turns to a house that nothing on the board shares.
type ActiveHouseMatchesNoCardsInPlay struct{}

// CondText renders the condition clause.
func (ActiveHouseMatchesNoCardsInPlay) CondText() string {
	return "which matches no cards in play"
}

// Met reports whether the active house matches no card in play. The first card
// of the active house makes the condition false. Upgrades count as cards in play,
// including the ones on an artifact, which the hand-rolled scan this replaced
// missed.
func (ActiveHouseMatchesNoCardsInPlay) Met(ctx *EffectContext) bool {
	house := ctx.Resolver.ActiveHouse()
	for _, p := range []int{ctx.Controller, ctx.Opponent()} {
		for _, id := range resolverCardsInPlay(ctx, p) {
			if ctx.Resolver.House(id) == house {
				return false
			}
		}
	}
	return true
}
