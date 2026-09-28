package engine

// Exclusion is the one card a target leaves out of the set it would otherwise
// select — what a card means when it says "other". The zero value excludes
// nothing, and the three real values differ only in which card they name, since
// every one of them prints the same word.
type Exclusion int

const (
	// excludeNone is the zero value: the target excludes no card.
	excludeNone Exclusion = iota
	// ExcludeSource excludes the card the ability is printed on ("each other
	// friendly card").
	ExcludeSource
	// ExcludeIt excludes the card in context (ctx.It) and nothing when no effect
	// has put one there, which is what "another" means for a selection made after
	// a first pick (pinned by TestChosenAnotherExcludesTheCardInContext).
	ExcludeIt
	// ExcludeFocus excludes the card in context when an effect has put one there
	// and the source card otherwise — the fallback that keeps "another creature"
	// from landing on the very card it is defined against (Guardian Demon deals to
	// another creature than the one it healed).
	ExcludeFocus
)

// allExclusions returns every real exclusion, so the census enumerates them
// rather than keeping a list a new value would fall out of.
func allExclusions() []Exclusion {
	return []Exclusion{ExcludeSource, ExcludeIt, ExcludeFocus}
}

// filters reports whether the exclusion drops any card at all.
func (e Exclusion) filters() bool { return e != excludeNone }

// excluded returns the card the exclusion drops in this context, and whether it
// drops one: ExcludeIt drops nothing when no effect has put a card in context.
func (e Exclusion) excluded(ctx *EffectContext) (LocalID, bool) {
	switch e {
	case ExcludeSource:
		return ctx.Source, true
	case ExcludeIt:
		return ctx.It, ctx.HasIt
	case ExcludeFocus:
		return focus(ctx), true
	}
	return 0, false
}

// admits reports whether a card survives the exclusion.
func (e Exclusion) admits(ctx *EffectContext, id LocalID) bool {
	excluded, drops := e.excluded(ctx)
	return !drops || id != excluded
}

// quantifyOne renders the whole phrase for a target that names a single card,
// given the noun with its side adjective already attached ("friendly creature").
// An exclusion reads "another", because "another friendly creature" is exactly "a
// friendly creature" that is not the card the exclusion names; without one the
// noun takes its indefinite article.
func (e Exclusion) quantifyOne(noun string) string {
	if e.filters() {
		return "another " + noun
	}
	return indefinite(noun)
}

// qualifyNoun prefixes the "other" qualifier to the noun the target names, e.g.
// "card" becomes "other card" for "each other friendly card". Every exclusion
// prints the same word: printed text names an exclusion by contrast, never by
// which card it leaves out.
func (e Exclusion) qualifyNoun(noun string) string {
	if e.filters() {
		return "other " + noun
	}
	return noun
}
