package engine

// DiscardOpponentArchivesOrDeckTop discards one card from a source the controller
// picks between their opponent's facedown archives (a uniformly random card) and
// the top of their opponent's deck, and binds the discarded card as ctx.It so a
// following effect can react to it — Fidgit plays it when it is a Tactic. An empty
// chosen source discards nothing and binds nothing.
type DiscardOpponentArchivesOrDeckTop struct{}

// validate accepts the effect; it has no configuration.
func (DiscardOpponentArchivesOrDeckTop) validate() error { return nil }

// Text renders the source choice as one clause.
func (DiscardOpponentArchivesOrDeckTop) Text() string {
	return "discard a random card from your opponent's archives or the top card " +
		"of their deck"
}

// Resolve discards from the chosen source and binds the discarded card in context.
func (DiscardOpponentArchivesOrDeckTop) Resolve(ctx *EffectContext) {
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
