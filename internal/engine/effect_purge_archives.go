package engine

// PurgeArchives lets the controller purge any number of cards from their own
// archives, one at a time, and records how many were purged for a following
// effect to scale by — the producer half of Destructive Analysis's "purge any
// number ... to deal an additional 2 damage for each card purged this way".
// Purging none is allowed and records nothing.
type PurgeArchives struct{}

// validate has nothing to reject: the effect takes no fields.
func (PurgeArchives) validate() error { return nil }

// Text renders the effect.
func (PurgeArchives) Text() string {
	return "purge any number of cards from your archives"
}

// Resolve purges cards from the controller's archives one at a time until they
// decline, recording the per-player tally a following CardsPurged reads.
func (PurgeArchives) Resolve(ctx *EffectContext) {
	chosen := pickCards(
		ctx,
		"Choose a card to purge from your archives",
		0,
		true,
		func() []LocalID { return ctx.Resolver.Archives(ctx.Controller) },
	)
	for _, id := range chosen {
		purgeFrom(ctx, Archives, ctx.Controller, id)
		ctx.Produced.Purged[ctx.Controller]++
	}
}
