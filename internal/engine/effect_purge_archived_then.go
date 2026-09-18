package engine

// PurgeArchivedCardThen lets the controller optionally purge one card from their
// own archives to pay for Then — Yzphyz Knowdrone's "you may purge an archived
// card to stun a creature". The purge is the cost: Then resolves only if a card
// is actually purged, and the controller may decline (or have no card to purge),
// in which case nothing happens. The optional-purge gate is intrinsic — no
// standalone purge-from-archives gate exists to compose with a general cost-then
// combinator — so it stays one node rather than May + a gate + Then.
type PurgeArchivedCardThen struct {
	// Then is the benefit that resolves once a card is purged.
	Then Effect
}

// validate requires a valid Then.
func (e PurgeArchivedCardThen) validate() error {
	return validateEffect(e.Then)
}

// Text renders the effect, e.g. "you may purge a card from your archives to stun
// a creature".
func (e PurgeArchivedCardThen) Text() string {
	return "you may purge a card from your archives to " + e.Then.Text()
}

// Resolve offers to purge one archived card; if the controller purges one, Then
// resolves. Declining, or an empty archive, does nothing.
func (e PurgeArchivedCardThen) Resolve(ctx *EffectContext) {
	cands := ctx.Resolver.Archives(ctx.Controller)
	if len(cands) == 0 {
		return
	}
	chosen, ok := ctx.ChooseCardOptional(
		"Choose a card to purge from your archives", cands)
	if !ok {
		return
	}
	purgeFrom(ctx, Archives, ctx.Controller, chosen)
	e.Then.Resolve(ctx)
}
