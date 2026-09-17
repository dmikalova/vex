package engine

// NameCard has the controller name any card, recording it on the source permanent
// (Etan's Jar). While the source stays in play a name-based play gate reads the
// stored name and bars every card of that name from being played, by either
// player; the bar lifts when the source leaves play (resetCore clears the name).
//
// The card is named from the cards present in the match rather than a global card
// database: every card that could ever be played is registered in the match
// catalog, so the choice reaches every playable name without the engine importing
// internal/cards (ADR 0003).
type NameCard struct{}

// Text renders the effect. The bar lasts until the source leaves play, so the dash
// binds the naming to its lasting consequence, the way ChooseHouseThen binds a
// house choice to what it does.
func (NameCard) Text() string {
	return "name a card - cards with that name cannot be played until " +
		SelfName + " leaves play"
}

// Resolve asks the controller to name one of the cards present in the match and
// records it on the source. The choice is presented as labelled options, one per
// distinct card name, in a deterministic order so a replay picks the same name.
func (NameCard) Resolve(ctx *EffectContext) {
	cands := ctx.Resolver.NameableCards()
	if len(cands) == 0 {
		return
	}
	labels := make([]string, len(cands))
	for i, id := range cands {
		labels[i] = ctx.Resolver.Name(id)
	}
	idx := ctx.ChooseOption("Name a card", labels)
	if idx < 0 || idx >= len(cands) {
		return
	}
	ctx.Resolver.SetNamedCard(ctx.Source, cands[idx])
}
