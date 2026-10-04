package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the Refinement family —
// the set-relative rules a Target carries, which compare the candidates to each
// other rather than testing each on its own.
//
// A refinement is not an effect, so it has no place in the glyphFamilies chain:
// it narrows a noun, so it draws after the noun it narrows (targetGlyphs) rather
// than as a clause of its own, the way a Count draws where the verb holding it
// draws.
//
// The grammar is two rules, and neither needs art of its own:
//
//   - A rule about power is the power noun carrying the comparison it makes, with
//     the decoration saying how much of the set it keeps — the each stack for a
//     tier, which keeps every creature at the extreme, and the chosen outline for a
//     singular selector, which keeps one and lets the controller break the tie. So
//     "the most powerful creature" and "each creature with the highest power" are
//     one noun under two decorations rather than two drawings.
//   - A rule about a battleline draws the creatures it picks on each side, tinted
//     friendly and enemy, with its number on the badge. A refinement that *spares*
//     what it names puts the refusal mark after it, the way every other barred
//     thing is drawn (AllExcept, KeepPerSide).
//
// A combinator has no picture of its own: it draws the refinements it composes,
// which is what the census calls it — plumbing that folds into its children.

// refinementGlyphs renders a Refinement as the narrowing it makes. The Refinement
// family is gated (TestEveryRefinementDraws, TestIconTotality), so every row here
// draws.
func refinementGlyphs(r engine.Refinement) []glyph {
	switch v := r.(type) {
	// The power rules: which end of the power scale, and how much of it.
	case engine.MostPowerfulN:
		// The numeral is drawn only for the top-N form, as the clause prints it only
		// there: "the most powerful creature", but "the 3 most powerful creatures".
		top := v.N
		if top == 1 {
			top = 0
		}
		return []glyph{powerRank(decorChosen, engine.AtLeast, top)}
	case engine.LeastPowerful:
		return []glyph{powerRank(decorChosen, engine.AtMost, 0)}
	case engine.HighestPower:
		return []glyph{powerRank(decorEach, engine.AtLeast, 0)}
	case engine.LowestPower:
		return []glyph{powerRank(decorEach, engine.AtMost, 0)}
	case engine.PowerLessThan:
		// The threshold is a running count rather than a number, so it draws as the
		// count's own glyphs after the bounded power noun. "Lower than" and "at most"
		// share the one ≤ mark: the strip says which way the bound points, and the
		// rules text beside it carries where exactly it falls.
		return append([]glyph{powerRank(0, engine.AtMost, 0)}, countGlyphs(v.Limit)...)
	case engine.SamePowerAs:
		return samePowerGlyphs(v)

	// The battleline rules, which act on each side in turn.
	case engine.KeepPerSide:
		// The creatures kept on each side are the ones the refinement spares, so they
		// draw behind the refusal mark: what the target reaches is everything else.
		return append(perSideGlyphs(v.N), glyph{asset: "glyph-ban"})
	case engine.PortionPerSide:
		// A share of each side rather than a count of it, so the numeral is the
		// denominator the side is cut into — half draws a 2, a third a 3 — which is
		// the same bargain lossGlyphs makes for a share of a pool.
		if d := fractionDenominator(v.Share); d > 0 {
			return perSideGlyphs(d)
		}
		return unmapped(v)

	// The combinators, which draw what they compose.
	case engine.AllExcept:
		return append(refinementGlyphs(v.Inner), glyph{asset: "glyph-ban"})
	case engine.AnyOf:
		gs := make([]glyph, 0, len(v.Refinements))
		for _, member := range v.Refinements {
			gs = append(gs, refinementGlyphs(member)...)
		}
		return gs
	case engine.Filter:
		// A Filter is a per-card narrowing that happens to sit inside a union, so it
		// draws the axes it sets and no noun: the noun is the one the target holding
		// the union already drew (filterAxisGlyphs).
		return filterAxisGlyphs(v)
	}
	return unmapped(r)
}

// powerRank is the power noun a refinement draws: the comparison it makes on the
// badge, the decoration saying how much of the set it keeps, and a numeral only
// where the refinement prints one. It is the power threshold cardsInPlayCountGlyphs
// already draws, borrowed so a rank and a filter on power read alike.
func powerRank(d decor, cmp engine.Comparison, qty int) glyph {
	return glyph{asset: "power", qty: qty, decor: d, comparison: cmp}
}

// samePowerGlyphs draws a power bracket measured against chosen creatures: the
// power noun under the equality mark, stacked because the whole bracket is kept,
// then the creature each pick takes. The picks carry their own side tint, so
// Quintrino Flux's friendly-and-enemy pair reads as two creatures rather than one.
func samePowerGlyphs(s engine.SamePowerAs) []glyph {
	gs := []glyph{powerRank(decorEach, engine.Exactly, 0)}
	for _, side := range s.Sides {
		gs = append(gs, glyph{
			asset: "type-creature",
			decor: decorChosen | sideDecor(side),
		})
	}
	return gs
}

// sideDecor tints a glyph by the battleline a pick draws from. A pick over both
// battlelines names no side, so it takes no tint.
func sideDecor(s engine.CreatureSide) decor {
	switch s {
	case engine.FriendlySide:
		return decorFriendly
	case engine.EnemySide:
		return decorEnemy
	}
	return 0
}

// perSideGlyphs draws n creatures chosen on each battleline — the shape both
// per-battleline refinements share, differing only in what the number means and in
// whether the creatures are the ones kept or the ones taken. The sides draw in one
// order whichever refinement draws them, so the pair reads as "on each battleline"
// rather than as a sequence.
func perSideGlyphs(n int) []glyph {
	return []glyph{
		{asset: "type-creature", qty: n, decor: decorChosen | decorEnemy},
		{asset: "type-creature", qty: n, decor: decorChosen | decorFriendly},
	}
}
