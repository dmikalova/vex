package engine

// NameCard has the controller name any card, recording it on the source permanent
// (Etan's Jar). While the source stays in play a name-based play gate reads the
// stored name and bars every card of that name from being played, by either
// player; the bar lifts when the source leaves play (resetCore clears the name).
//
// The names offered are the whole implemented card database, injected by the match
// (ADR 0003). Offering only the names present in the match would show a player
// their opponent's deck list.
type NameCard struct{}

// Text renders the effect as the naming and its lasting consequence, two
// sentences the way every other lead-in choice reads.
func (NameCard) Text() string {
	return leadInSentence("name a card", "cards with that name cannot be played "+
		durationClause(UntilThisLeavesPlay, SelfName))
}

func (NameCard) endsSentence() bool { return true }

// Resolve asks the controller to name a card and records it on the source. A name
// no card in the match carries is a legal but idle choice — nothing of that name
// can be played — so it records nothing rather than a card that is not there.
func (NameCard) Resolve(ctx *EffectContext) {
	names := ctx.Resolver.NameableNames()
	if len(names) == 0 {
		return
	}
	idx := ctx.ChooseOption("Name a card", names)
	if idx < 0 || idx >= len(names) {
		return
	}
	if id, ok := ctx.Resolver.CardNamed(names[idx]); ok {
		ctx.Resolver.SetNamedCard(ctx.Source, id)
	}
}
