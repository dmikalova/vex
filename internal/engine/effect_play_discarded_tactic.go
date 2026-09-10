package engine

// PlayDiscardedTacticFromOpponent is Fidgit's reap: the controller discards a
// card from their opponent's archives or the top of their opponent's deck, and if
// that discarded card is a Tactic, plays it from the opponent's discard pile as if
// it were their own. The controller picks the source; a random card leaves the
// facedown archives, while the deck contributes its named top card. An empty
// chosen source discards nothing.
type PlayDiscardedTacticFromOpponent struct{}

// validate accepts the effect; it has no configuration.
func (PlayDiscardedTacticFromOpponent) validate() error { return nil }

// Text renders the two printed sentences as one unit.
func (PlayDiscardedTacticFromOpponent) Text() string {
	return "discard a random card from your opponent's archives or the top card " +
		"of their deck. If that card is a Tactic, play it as if it were yours"
}

// Resolve discards from the chosen source, binds the discarded card in context,
// and plays it from the opponent's discard pile when it is a Tactic.
func (PlayDiscardedTacticFromOpponent) Resolve(ctx *EffectContext) {
	opp := ctx.Opponent()
	var (
		id LocalID
		ok bool
	)
	if ctx.ChooseOption(
		"Which card to discard?",
		[]string{"your opponent's archives", "the top card of their deck"},
	) == 0 {
		id, ok = discardRandomArchivesCard(ctx, opp)
	} else {
		id, ok = ctx.Resolver.DiscardTopOfDeck(opp)
	}
	if !ok {
		return
	}
	ctx.It, ctx.HasIt = id, true
	if ctx.Resolver.TypeOf(id) != Tactic {
		return
	}
	ctx.Resolver.PlayFromOpponentDiscard(ctx.Controller, id)
}

// discardRandomArchivesCard picks a random card from player's archives, discards
// it, and reports it. Empty archives discard nothing.
func discardRandomArchivesCard(ctx *EffectContext, player int) (LocalID, bool) {
	id, ok := ctx.Resolver.ChooseRandom(ctx.Resolver.Archives(player))
	if !ok {
		return 0, false
	}
	ctx.Resolver.DiscardCardFromArchives(player, id)
	return id, true
}
