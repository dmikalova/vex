package engine

// Position is the battleline place a target narrows by. The zero value asks
// nothing, and the five real values are exclusive: a card says "a flank
// creature" or "a creature to the right of <self>", never both. A card that ever
// wants the union of two places writes a target per place rather than two
// positions in one.
type Position int

const (
	// positionAny is the zero value: the filter asks nothing about position.
	positionAny Position = iota
	// PositionOnFlank admits a creature on a flank of its battleline — its
	// leftmost or rightmost creature. A flank is a battleline position, so the
	// filter only constrains creatures: on a target that also reaches artifacts
	// ("an artifact or flank creature", Snudge) an artifact passes it untouched.
	PositionOnFlank
	// PositionNotOnFlank admits a creature that is neither its battleline's
	// leftmost nor its rightmost.
	PositionNotOnFlank
	// PositionCenter admits the creature in the center of its controller's
	// battleline; an even-sized line has no center (Beware the Ides).
	PositionCenter
	// PositionRightOfSource admits a creature standing to the source card's right
	// in its battleline (Panpaca, Anga).
	PositionRightOfSource
	// PositionLeftOfSource admits a creature standing to the source card's left in
	// its battleline (Panpaca, Jaga).
	PositionLeftOfSource
)

// allPositions returns every real position, so the census enumerates them rather
// than keeping a list a new value would fall out of.
func allPositions() []Position {
	return []Position{
		PositionOnFlank,
		PositionNotOnFlank,
		PositionCenter,
		PositionRightOfSource,
		PositionLeftOfSource,
	}
}

// filters reports whether the position narrows its candidates at all.
func (p Position) filters() bool { return p != positionAny }

// admits reports whether a card standing where it does passes the filter.
func (p Position) admits(ctx *EffectContext, id LocalID) bool {
	switch p {
	case PositionOnFlank:
		return !ctx.Resolver.IsCreature(id) || onFlank(ctx, id)
	case PositionNotOnFlank:
		return !onFlank(ctx, id)
	case PositionCenter:
		return ctx.Resolver.InCenterOfBattleline(id)
	case PositionRightOfSource:
		return toSideOfSource(ctx, ctx.Source, id, +1)
	case PositionLeftOfSource:
		return toSideOfSource(ctx, ctx.Source, id, -1)
	}
	return true
}

// adjective is the word the position prefixes to the noun it narrows, e.g.
// "flank" in "each flank creature". It is empty for a position that renders as a
// trailing clause instead.
func (p Position) adjective() string {
	if p == PositionOnFlank {
		return "flank"
	}
	return ""
}

// clause appends the position's own clause to a quantified phrase, e.g. "each
// creature" becomes "each creature that is not on a flank". A position that
// renders as an adjective, or asks nothing, appends nothing.
func (p Position) clause(phrase string) string {
	switch p {
	case PositionNotOnFlank:
		return phrase + " that is not on a flank"
	case PositionCenter:
		return phrase + " in the center of its controller's battleline"
	case PositionRightOfSource:
		return phrase + " to the right of " + SelfName
	case PositionLeftOfSource:
		return phrase + " to the left of " + SelfName
	}
	return phrase
}
