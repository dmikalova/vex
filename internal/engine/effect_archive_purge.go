package engine

// ArchivePurgedCard lets the controller take one card out of their own purge
// pile and set it in their archives — Universal Recycle Bin recovering a card
// the game had set aside for good. It does nothing when the controller's purge
// pile is empty.
type ArchivePurgedCard struct{}

// Text renders the effect, e.g. "archive a purged card you own".
func (e ArchivePurgedCard) Text() string { return "archive a purged card you own" }

// Resolve has the controller choose one card from their purge pile and move it
// to their archives, doing nothing if the purge pile is empty.
func (e ArchivePurgedCard) Resolve(ctx *EffectContext) {
	cards := ctx.Resolver.Purge(ctx.Controller)
	if len(cards) == 0 {
		return
	}
	id, _ := ctx.ChooseCard("Choose a purged card to archive", cards)
	ToArchives.moveFrom(ctx, purged, ctx.Controller, id)
}
