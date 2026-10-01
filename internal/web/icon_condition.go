package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the condition family
// (ADR 0047) — the engine's Condition nodes, the "if …" clause a card gates an
// effect behind ("While your opponent has 6 or more Æmber, destroy each
// creature").
//
// A condition is not an effect, so like a count it has no place in the
// glyphFamilies chain: a condition is drawn where the node holding it is drawn.
// Until this file existed none of them drew anything at all, so a card reading
// "While your opponent has 6 or more Æmber, destroy each creature" drew exactly
// the strip a bare "destroy each creature" drew.
//
// The grammar is one rule — a condition draws as the gate mark, then the noun it
// tests — and the noun is composed from vocabulary that is already here:
//
//   - A condition that compares a quantity against a figure draws the Count's own
//     picture, wearing the comparison and the figure on its badge
//     (countedThreshold). It restates no count: PoolAember is the Æmber-in-pool
//     count with an enemy tint and a "≥6", ForgedKey is the keys-forged-this-turn
//     tally with a "≥1".
//   - A condition that tests a property of one card draws that property's own
//     glyph, decorated with the card it reads.
//   - A condition that tests what has happened this turn draws the verb glyph the
//     action itself draws, closed by the turn marker.
//
// Several conditions ask a question a Target filter also asks — a creature with
// Æmber on it, a stunned one, one on a flank. They draw the same picture on
// purpose: "each friendly creature with Æmber on it" and "if it has Æmber on it"
// are one idea, and the strip should not spell them two ways.
//
// And, Or and Not carry no test of their own and fold into their children, the
// way a plumbing effect folds into the effects it wraps. Or divides its children
// with the alternative mark, because they are alternatives rather than a run.
// Not must still negate, so it closes its child with the refusal mark already in
// the vocabulary: folding into a child unmarked would draw the opposite of what
// the card says.

// conditionGlyphs renders a Condition as the clause it opens: the gate mark, then
// the test itself. It reports whether a real mapping covered the condition; one
// with no mapping draws the abstract unknown glyph, which is what the icon gate
// fails on.
//
// The gate mark lives here rather than on the effect the condition gates, so a
// condition drawn anywhere — a Conditional, an OrAmount, a constant ability's
// "while …" — opens its clause the same way.
func conditionGlyphs(c engine.Condition) ([]glyph, bool) {
	gs, covered := conditionTestGlyphs(c)
	return punctuate(gs, punctGate), covered
}

// conditionTestGlyphs renders the test a condition makes, without the gate mark
// that opens the clause. It is the recursion point: a combinator folds into its
// children through this rather than through conditionGlyphs, so an And of three
// conditions opens one clause rather than three.
func conditionTestGlyphs(c engine.Condition) ([]glyph, bool) {
	switch v := c.(type) {
	// Combinators, which carry no test of their own.
	case engine.And:
		return foldedConditions(v.Conditions, punctNone)
	case engine.Or:
		return foldedConditions(v.Conditions, punctDivider)
	case engine.Not:
		gs, covered := conditionTestGlyphs(v.Cond)
		return append(gs, glyph{asset: "glyph-ban"}), covered
	case engine.CountIs:
		return countedThreshold(v.Count, v.Is, v.Amount)

	// The board: what each player controls, holds, and has forged.
	case engine.PoolAember:
		return countedThreshold(engine.AemberInPool{Player: v.Player}, v.Is, v.Amount)
	case engine.ControlsMoreCreatures:
		// The condition is a threshold on the excess count, and the count already
		// carries the "more than the other side" comparison, so it draws itself.
		return countGlyphs(engine.ExcessCreatures{Player: engine.Controller, Trait: v.Trait})
	case engine.Overwhelmed:
		return countGlyphs(engine.ExcessCreatures{Player: engine.Opponent})
	case engine.HousesRepresented:
		return countedThreshold(v.Among, v.Is, v.Amount)
	case engine.CounterInPlay:
		return counterInPlayGlyphs(v)
	case engine.ForgedKey:
		return countedThreshold(
			engine.TurnCount{Player: v.Player, Of: forgedKeyStat(v.Previous)},
			engine.AtLeast,
			1,
		)
	case engine.HasMoreForgedKeys:
		gs, covered := countGlyphs(engine.ForgedKeys(v))
		return threshold(gs, aheadOfOther(v.Player), 0), covered
	case engine.KeyColorForged:
		return keyColorForgedGlyphs(v)

	// The discard pile.
	case engine.CardsInDiscardAtLeast:
		gs := append(
			houseFilterGlyphs(v.House),
			glyph{asset: countedCardAsset(v.Type), decor: decorEach},
			glyph{asset: "zone-discard"},
		)
		return threshold(gs, engine.AtLeast, v.Amount), true
	case engine.DiscardedThisWay:
		// The cards this resolution just put in the pile rather than the cards
		// sitting in it, so the pile takes the flow arrow the discard effects draw.
		return append(
			houseFilterGlyphs(v.House),
			glyph{asset: countedCardAsset(v.Type), decor: decorEach},
			arrowTo(glyph{asset: "zone-discard"}),
		), true
	case engine.CardsDiscarded:
		// The cards a player sent to their discard pile this turn. There is no Count
		// for this tally to borrow, so it draws the shape every turn tally draws: the
		// cards, what became of them, and the turn marker.
		gs := thisTurn(append(
			houseFilterGlyphs(v.House),
			cardsCounted(playerDecor(v.Player)),
			arrowTo(glyph{asset: "zone-discard"}),
		))
		return threshold(gs, engine.AtLeast, v.Amount), true

	// The card in context, and the cards in play beside it.
	case engine.HasAember:
		return []glyph{{asset: "aember", decor: subjectDecor(v.Subject)}}, true
	case engine.CardsInPlay:
		return cardsInPlayConditionGlyphs(v), true
	case engine.CountersOnThisAtLeast:
		return countedThreshold(
			engine.CountersOnThis{Kind: v.Kind},
			engine.AtLeast,
			v.N,
		)

	// The turn so far.
	case engine.CardsPlayed:
		return countedThreshold(v, engine.AtLeast, atLeastOne(v.Amount))
	case engine.AemberStolenFromYou:
		return countedThreshold(
			engine.TurnCount{Player: engine.Controller, Of: engine.AemberStolenFromLastTurn},
			engine.AtLeast,
			1,
		)
	case engine.MovedAnyAember:
		// A tally of what a preceding MoveAember relocated, so it draws that move —
		// Æmber onto the card in context — with the "at least one" the condition asks
		// for on its badge.
		return threshold([]glyph{
			{asset: "aember"},
			arrowTo(glyph{asset: "card-back", decor: decorChosen}),
		}, engine.AtLeast, 1), true
	}
	return fallbackGlyphs(), false
}

// foldedConditions renders a combinator's children one after another, reporting
// covered only when every child maps to a real glyph. between opens every child
// after the first: an Or's children are alternatives and take the divider mark,
// an And's are a run and take none.
func foldedConditions(conds []engine.Condition, between punct) ([]glyph, bool) {
	out := make([]glyph, 0, len(conds)*2)
	covered := true
	for i, cond := range conds {
		gs, ok := conditionTestGlyphs(cond)
		if i > 0 {
			gs = punctuate(gs, between)
		}
		out = append(out, gs...)
		covered = covered && ok
	}
	return out, covered
}

// countedThreshold draws a condition that compares a count against a figure: the
// count's own picture, wearing the comparison and the figure on its head. It is
// how the condition family reuses the count grammar rather than restating it —
// "if your opponent has 6 or more Æmber" is the Æmber-in-pool count with an enemy
// tint and a "≥6", not a second drawing of a pool.
func countedThreshold(c engine.Count, is engine.Comparison, amount int) ([]glyph, bool) {
	gs, covered := countGlyphs(c)
	return threshold(gs, is, amount), covered
}

// threshold stamps a comparison and its figure on the head of an already-composed
// run, which is how a condition says what it compares a count against without
// threading the figure through the helper that drew the count. An empty run is
// left alone.
//
// It is the counterpart of markVerb and punctuate: a run marked as a whole rather
// than glyph by glyph. Unlike those two the badge does not ride the head, because
// a count's run often opens with a mark that narrows it rather than with the noun
// it totals — "3 or more Untamed creatures in your discard pile" draws the house
// emblem first, and a "≥3" on the emblem would read as three houses. It rides the
// stacked noun instead (countedNoun), which is the thing being totalled.
func threshold(gs []glyph, is engine.Comparison, amount int) []glyph {
	i := countedNoun(gs)
	if i < 0 {
		return gs
	}
	gs[i].comparison = is
	gs[i].qty = comparedFigure(is, amount)
	return gs
}

// countedNoun is the index of the glyph a threshold badges: the first stacked
// one, because a count stacks the noun it totals (creaturesCounted,
// cardsCounted), or the head for a count of something singular — a pool, the
// counters on this card — which stacks nothing. It is -1 for an empty run.
func countedNoun(gs []glyph) int {
	for i := range gs {
		if gs[i].decor&decorEach != 0 {
			return i
		}
	}
	if len(gs) == 0 {
		return -1
	}
	return 0
}

// comparedFigure is the number a comparison prints on its badge, or 0 for a
// comparison that ignores Amount — the two parity marks compare against no
// figure at all, and the two relative marks compare against the other player's
// total rather than a number. The engine's Met ignores Amount for all four, so
// drawing it would print a figure the rule never reads.
func comparedFigure(is engine.Comparison, amount int) int {
	switch is {
	case engine.Even, engine.Odd, engine.MoreThanYou, engine.MoreThanOpponent:
		return 0
	}
	return amount
}

// aheadOfOther is the comparison for "this side has more of it than the other
// side". Which side is ahead is said once — by the glyph's tint — so the two
// relative comparisons differ only in whose total is the subject (ADR 0047).
func aheadOfOther(p engine.Player) engine.Comparison {
	if p == engine.Opponent {
		return engine.MoreThanYou
	}
	return engine.MoreThanOpponent
}

// atLeastOne is the threshold a condition whose Amount defaults to one compares
// against, mirroring the engine's own default so the strip cannot print a "≥0"
// the rule never means.
func atLeastOne(amount int) int {
	if amount < 1 {
		return 1
	}
	return amount
}

// forgedKeyStat is the turn tally a forged-key condition reads: this turn's keys,
// or the player's previous turn's. Which turn is not drawn, for the reason
// turnStatGlyphs does not draw it — the turn marker says the tally belongs to a
// turn, and saying which would need a second picture.
func forgedKeyStat(previous bool) engine.TurnStat {
	if previous {
		return engine.KeysForgedLastTurn
	}
	return engine.KeysForgedThisTurn
}

// keyColorForgedGlyphs draws "your red key is forged": the forge the key count
// already draws, tinted with whose keys are read, then the key of the colour the
// condition names. A condition naming no colour is a half-built value no card can
// carry, and it draws the unknown glyph rather than a bare forge.
func keyColorForgedGlyphs(c engine.KeyColorForged) ([]glyph, bool) {
	a := keyColorIconName(c.Color)
	if a == "" {
		return fallbackGlyphs(), false
	}
	return []glyph{
		{asset: "forge", decor: playerDecor(c.Player)},
		{asset: a},
	}, true
}

// counterInPlayGlyphs draws "there is a card in play with a doom counter on it":
// the cards in play, holding the counter the condition looks for, with the "at
// least one" it asks for on the badge. A kind with no icon draws the unknown
// glyph, as the matching count does.
func counterInPlayGlyphs(c engine.CounterInPlay) ([]glyph, bool) {
	a := counterAsset(c.Kind)
	if a == "" {
		return fallbackGlyphs(), false
	}
	gs := []glyph{cardsCounted(0), {asset: a}}
	return threshold(gs, engine.AtLeast, 1), true
}

// cardsInPlayConditionGlyphs draws the Condition role of the one node that serves
// both roles from one description. The picture is the Count role's, because the
// cards surveyed are the same cards; the Condition role adds only what it alone
// carries — the refusal mark for "if there are no …", and the threshold for "if
// there are N or more …". A condition that names neither draws exactly what the
// count draws, which is the point: one Go value, one picture, read in whichever
// slot it sits.
func cardsInPlayConditionGlyphs(e engine.CardsInPlay) []glyph {
	gs := cardsInPlayCountGlyphs(e)
	if e.None {
		return append(gs, glyph{asset: "glyph-ban"})
	}
	if e.Amount > 0 {
		return threshold(gs, engine.AtLeast, e.Amount)
	}
	return gs
}

// subjectDecor marks the card a condition reads: the card the ability is printed
// on, or the card an effect or trigger just put in context. It is the condition
// family's counterpart of playerDecor, which marks whose side a glyph is about.
func subjectDecor(s engine.Subject) decor {
	if s == engine.This {
		return decorThis
	}
	return decorChosen
}
