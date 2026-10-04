package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the Filter family — the
// per-card axes a Target narrows by, from the house a card belongs to down to
// the counters sitting on it (engine/filter.go).
//
// A filter is drawn one axis at a time, never as a summary of the kind of target
// holding it. Three axes reshape the noun itself — the card type, a gigantic
// creature, a card named outright — because that is what they do in the printed
// text, where "each creature" becomes "each artifact" or "each Ancient Bear".
// Every other axis draws a run of its own after the noun, in the order the text
// prints its adjectives and clauses, so the strip narrows in the order the rules
// text beside it does.
//
// One axis per glyph is what keeps the drawing composable: whichever noun a
// target ends up handing in — the Kind's today, a quantifier/side/type
// composition once the board-scanning kinds are decomposed — nothing here
// changes.
//
// MatchAny is the one axis with no mark of its own. It is the filter's "or", so
// it punctuates instead: every axis run after the first opens with the divider a
// ChooseOne's options already use, and "each Mars or Robot creature" reads as one
// noun with two alternatives rather than as two filters.

// filterGlyphs renders a Filter over the noun it narrows: that noun carrying the
// axes that reshape it, then one run per remaining axis. A filter that narrows
// nothing draws the noun alone, which is what most targets do.
func filterGlyphs(f engine.Filter, noun glyph) []glyph {
	return filterRuns(filterNoun(f, noun), f)
}

// filterRowGlyphs draws one row of the Filter census: the axis over the "each
// creature" the census renders that row's text on (catalog_filter.go), so the
// strip a row draws is the strip of the phrase beside it. The Filter family is
// gated (TestEveryFilterAxisDraws, TestIconTotality), so every axis here draws.
func filterRowGlyphs(f engine.Filter) []glyph {
	return targetGlyphs(engine.Target{Kind: engine.TargetEachCreature, Filter: f})
}

// filterAxisGlyphs renders a Filter with no noun of its own, for a filter
// serving as a Refinement inside a union: the noun is the one its target already
// drew, so only the narrowing is left to draw (Regrettable Meteor unions a trait
// against a power floor). A member that names a type or a card draws that noun,
// since the target's own is wider than the member.
func filterAxisGlyphs(f engine.Filter) []glyph {
	return filterRuns(filterNoun(f, glyph{}), f)
}

// filterRuns assembles a drawn filter: the noun first, unless there is none to
// draw, then the axis runs in printed order, joined by the "or" divider when the
// filter disjoins them.
func filterRuns(noun glyph, f engine.Filter) []glyph {
	var gs []glyph
	if noun.asset != "" {
		gs = append(gs, noun)
	}
	for i, run := range filterAxisRuns(f) {
		if f.MatchAny && i > 0 {
			run = punctuate(run, punctDivider)
		}
		gs = append(gs, run...)
	}
	return gs
}

// filterNoun applies the axes that reshape the noun rather than trailing it, in
// the order the text substitutes them: a card type replaces the noun outright, a
// gigantic creature is a noun of its own, and a card named outright drops the
// noun for the nameplate. A filter that sets none of them leaves the noun as it
// found it.
func filterNoun(f engine.Filter, noun glyph) glyph {
	if a := typeIconName(f.Type); a != "" {
		noun.asset = a
	}
	if f.Gigantic {
		noun.asset = "glyph-gigantic"
	}
	if f.Name != "" {
		noun.asset = nameGlyph().asset
	}
	return noun
}

// filterAxisRuns returns one run of glyphs per axis the filter sets, in printed
// order. They are grouped the way the Filter struct is: what a card is, what has
// happened to it, and where it stands. One run per axis rather than one flat
// list is what lets MatchAny punctuate between axes.
func filterAxisRuns(f engine.Filter) [][]glyph {
	runs := filterIdentityRuns(f)
	runs = append(runs, filterStateRuns(f)...)
	return append(runs, filterPositionRuns(f)...)
}

// filterIdentityRuns draws the axes decided by what a card is: the house it
// belongs to, the traits it carries or lacks, and the two rules that measure a
// house by the creatures in it.
func filterIdentityRuns(f engine.Filter) [][]glyph {
	var runs [][]glyph
	if f.House != (engine.HouseMatcher{}) {
		runs = append(runs, houseMatcherGlyphs(f.House))
	}
	if f.Trait != 0 {
		runs = append(runs, []glyph{traitGlyph()})
	}
	if f.ExceptTrait != 0 {
		runs = append(runs, []glyph{traitGlyph(), {asset: "glyph-ban"}})
	}
	if f.SharesTrait {
		// A trait held in common with the card in context, which is drawn as the card
		// wearing the chosen outline wherever a condition reads "it".
		runs = append(runs, []glyph{traitGlyph(), inContext(glyph{asset: "card-back"})})
	}
	if f.WithoutSharedTrait {
		// The loner: a trait shared with its own battleline, behind the refusal mark.
		runs = append(runs, []glyph{
			traitGlyph(),
			{asset: "type-creature", decor: decorEach | decorFriendly},
			{asset: "glyph-ban"},
		})
	}
	if f.HouseWithMostCreatures {
		runs = append(runs, housePopulationGlyphs(0))
	}
	if f.HouseWithAtLeast > 0 {
		runs = append(runs, housePopulationGlyphs(f.HouseWithAtLeast))
	}
	return runs
}

// housePopulationGlyphs draws a house measured by how many creatures it has in
// play: the abstract house, then its creatures stacked under the at-least mark.
// n is 0 for the superlative — "the house with the most creatures in play" names
// no number, the way a refinement's top rank prints none (refinementGlyphs).
func housePopulationGlyphs(n int) []glyph {
	return []glyph{
		{asset: "house-any"},
		{asset: "type-creature", qty: n, decor: decorEach, comparison: engine.AtLeast},
	}
}

// filterStateRuns draws the axes decided by what has happened to a card: its
// power, the damage and Æmber and counters on it, its armor and upgrades, the
// keyword it has, and whether it is stunned or ready.
func filterStateRuns(f engine.Filter) [][]glyph {
	var runs [][]glyph
	if f.Power != (engine.PowerBound{}) {
		runs = append(runs, powerBoundGlyphs(f.Power))
	}
	if f.Damage != 0 {
		runs = append(runs, damagePresenceGlyphs(f.Damage))
	}
	if f.Aember != 0 {
		runs = append(runs, aemberPresenceGlyphs(f.Aember))
	}
	if f.Keyword != 0 {
		runs = append(runs, keywordFilterGlyphs(f.Keyword))
	}
	if f.Counter != engine.CounterNone {
		runs = append(runs, counterFilterGlyphs(f.Counter))
	}
	if f.Stunned {
		runs = append(runs, []glyph{{asset: "stun"}})
	}
	if f.Ready {
		// Ready is drawn the way the ready verb is: the exhausted mark tinted the way
		// a creature that is upright reads, so one picture means one state.
		runs = append(runs, []glyph{{asset: "exhausted", decor: decorFriendly}})
	}
	if f.Armor {
		runs = append(runs, []glyph{{asset: "shield"}})
	}
	if f.Upgrade {
		runs = append(runs, []glyph{{asset: "type-upgrade"}})
	}
	if f.NoBonusIcons {
		runs = append(runs, []glyph{{asset: "glyph-bonus"}, {asset: "glyph-ban"}})
	}
	return runs
}

// filterPositionRuns draws the axes decided by where a card stands: its place in
// the battleline, its neighbors, and the one card the filter leaves out.
func filterPositionRuns(f engine.Filter) [][]glyph {
	var runs [][]glyph
	if f.Position != 0 {
		runs = append(runs, positionGlyphs(f.Position))
	}
	if f.Neighboring {
		runs = append(runs, []glyph{{asset: "glyph-neighbors"}})
	}
	if f.SharesHouseWithNeighbors > 0 {
		runs = append(runs, []glyph{
			{asset: "house-any"},
			{
				asset:      "glyph-neighbors",
				qty:        f.SharesHouseWithNeighbors,
				comparison: engine.AtLeast,
			},
		})
	}
	if f.Except != 0 {
		runs = append(runs, exclusionGlyphs(f.Except))
	}
	return runs
}

// powerBoundGlyphs draws a power filter as the power noun under the comparison
// it makes, with the number it measures against on the badge. It is the glyph a
// refinement's power rank draws (powerRank), borrowed so a bound and a rank read
// alike. A parity bound has no number to print, and a bound measured against the
// source card's own power wears the self marker instead of one.
func powerBoundGlyphs(b engine.PowerBound) []glyph {
	switch b.Kind {
	case engine.BoundAtMost:
		return []glyph{powerRank(0, engine.AtMost, b.Amount)}
	case engine.BoundAtLeast:
		return []glyph{powerRank(0, engine.AtLeast, b.Amount)}
	case engine.BoundExactly:
		return []glyph{powerRank(0, engine.Exactly, b.Amount)}
	case engine.BoundOdd:
		return []glyph{powerRank(0, engine.Odd, 0)}
	case engine.BoundEven:
		return []glyph{powerRank(0, engine.Even, 0)}
	case engine.BoundLessThanSource:
		return []glyph{powerRank(decorThis, engine.AtMost, 0)}
	}
	return unmapped(b)
}

// damagePresenceGlyphs draws whether a creature carries damage: the damage mark,
// or the damage mark behind the refusal mark for one carrying none.
func damagePresenceGlyphs(d engine.DamagePresence) []glyph {
	switch d {
	case engine.DamageSome:
		return []glyph{{asset: "damage"}}
	case engine.DamageNone:
		return []glyph{{asset: "damage"}, {asset: "glyph-ban"}}
	}
	return unmapped(d)
}

// aemberPresenceGlyphs draws whether a card carries Æmber, the same way the
// damage axis draws damage.
func aemberPresenceGlyphs(a engine.AemberPresence) []glyph {
	switch a {
	case engine.AemberSome:
		return []glyph{{asset: "aember"}}
	case engine.AemberNone:
		return []glyph{{asset: "aember"}, {asset: "glyph-ban"}}
	}
	return unmapped(a)
}

// keywordFilterGlyphs draws the keyword a filter narrows by, as the keyword's own
// mark: "each elusive creature" is the creature and the elusive badge.
func keywordFilterGlyphs(k engine.Keyword) []glyph {
	if a := keywordIcon(k); a != "" {
		return []glyph{{asset: a}}
	}
	return unmapped(k)
}

// counterFilterGlyphs draws the generic counter a filter narrows by, as that
// counter's own icon.
func counterFilterGlyphs(k engine.CounterKind) []glyph {
	if a := counterAsset(k); a != "" {
		return []glyph{{asset: a}}
	}
	return unmapped(k)
}

// positionGlyphs draws the place in a battleline a filter narrows to. A flank is
// the mark the flank rules already draw; the places measured from the source card
// — its neighbors' slots to its left and right — draw the battleline with the
// slot filled and wear the self marker, since the line is read from that card.
func positionGlyphs(p engine.Position) []glyph {
	switch p {
	case engine.PositionOnFlank:
		return []glyph{{asset: "glyph-flank"}}
	case engine.PositionNotOnFlank:
		return []glyph{{asset: "glyph-flank"}, {asset: "glyph-ban"}}
	case engine.PositionCenter:
		return []glyph{{asset: "glyph-center"}}
	case engine.PositionRightOfSource:
		return []glyph{{asset: "glyph-right-of", decor: decorThis}}
	case engine.PositionLeftOfSource:
		return []glyph{{asset: "glyph-left-of", decor: decorThis}}
	}
	return unmapped(p)
}

// exclusionGlyphs draws the one card a filter leaves out, behind the refusal
// mark. The printed text says "other" whichever card that is, but the strip can
// say which: the card the ability is printed on wears the self marker, the card
// in context the chosen outline, and the fallback that means either wears both.
func exclusionGlyphs(e engine.Exclusion) []glyph {
	switch e {
	case engine.ExcludeSource:
		return []glyph{{asset: "card-back", decor: decorThis}, {asset: "glyph-ban"}}
	case engine.ExcludeIt:
		return []glyph{inContext(glyph{asset: "card-back"}), {asset: "glyph-ban"}}
	case engine.ExcludeFocus:
		return []glyph{
			inContext(glyph{asset: "card-back", decor: decorThis}),
			{asset: "glyph-ban"},
		}
	}
	return unmapped(e)
}
