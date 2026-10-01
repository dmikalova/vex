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
//     glyph, decorated with the card it reads (inContext for the card an effect or
//     trigger just put in context, decorThis for the card the ability is printed
//     on). ItIsStunned is the stun glyph, ItIsOfTrait the trait mark, ItIsNamed the
//     nameplate, ItIs a house matcher, OnFlank the battleline mark.
//   - A condition that tests what has happened this turn draws the verb glyph the
//     action itself draws, closed by the turn marker — which is the turn tally the
//     Count family already draws, wearing a "≥1" for "you have" and an "=1" for
//     "this is the first". "First" and "at least once" are the same picture read
//     against a different figure, not two drawings.
//
// The tide is the one noun this family adds to the vocabulary, and it adds one:
// high and low are two readings of the same mark, carried by the increase/decrease
// sigil pair, because a tide that is high for you is the same tide that is low for
// your opponent (TestTideReadsOneMarkTwoWays).
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
	case engine.AlwaysMet:
		// The census declares it silent: an always-met condition adds no clause, so
		// it draws neither a noun nor the gate mark that would open one.
		return nil, true
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
	case engine.ControlsNamed:
		// "If you control a Velum" is the cards-in-play condition narrowed by a name,
		// so it draws that condition's picture rather than a second one of its own.
		return cardsInPlayConditionGlyphs(engine.CardsInPlay{
			Player: engine.Controller,
			Name:   v.Name,
		}), true
	case engine.NamedCardPurged:
		return []glyph{nameGlyph(), {asset: "zone-purge"}}, true
	case engine.ActiveHouseMatchesNoCardsInPlay:
		// The active house, then every card in play, behind the refusal mark: no card
		// on either board shares it.
		return append(
			houseMatcherGlyphs(engine.HouseMatcher{Kind: engine.MatchActiveHouse}),
			cardsCounted(0),
			glyph{asset: "glyph-ban"},
		), true

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
	case engine.NamedCardInDiscard:
		return []glyph{nameGlyph(), {asset: "zone-discard"}}, true
	case engine.Haunted:
		// Haunted is a word for a pile size, so it draws the pile the word counts
		// rather than a picture of its own: your discard pile wearing a "≥10".
		return countedThreshold(
			engine.CardsInZone{Player: engine.Controller, Zone: engine.Discard},
			engine.AtLeast,
			hauntedDiscardSize,
		)

	// The card in context, and the cards in play beside it.
	case engine.HasAember:
		return []glyph{{asset: "aember", decor: subjectDecor(v.Subject)}}, true
	case engine.ItIs:
		// Which houses the card may belong to, then the card itself. Other ("another
		// creature") only excludes the source card, so like every other fine filter it
		// stays in the rules text beside the strip.
		return append(
			houseFilterGlyphs(v.House),
			inContext(glyph{asset: countedCardAsset(v.Type)}),
		), true
	case engine.ItIsFriendly:
		return []glyph{inContext(glyph{asset: "type-creature", decor: decorFriendly})}, true
	case engine.ItIsEnemy:
		return []glyph{inContext(glyph{asset: "type-creature", decor: decorEnemy})}, true
	case engine.ItIsOfTrait:
		return []glyph{inContext(glyph{asset: "type-creature"}), traitGlyph()}, true
	case engine.ItIsNamed:
		return []glyph{inContext(nameGlyph())}, true
	case engine.ItIsStunned:
		return []glyph{inContext(glyph{asset: "stun"})}, true
	case engine.ItHasBonusIcon:
		return []glyph{inContext(glyph{asset: "glyph-bonus"})}, true
	case engine.ItIsOffIdentity:
		// Your identity is the three houses of your deck, which is the abstract house
		// stacked and tinted friendly; the card in context belongs to none of them.
		return []glyph{
			inContext(glyph{asset: "card-back"}),
			{asset: "house-any", decor: decorEach | decorFriendly},
			{asset: "glyph-ban"},
		}, true
	case engine.ItIsNotOfNamedHouse:
		// The named house draws as the chosen house, which is what the effect that
		// names one draws (OpponentNamesHouse); the refusal mark is this condition's
		// own, since the card in context must not be of it.
		return append(
			houseMatcherGlyphs(engine.HouseMatcher{Kind: engine.MatchChosenHouse}),
			inContext(glyph{asset: "card-back"}),
			glyph{asset: "glyph-ban"},
		), true
	case engine.ItAttachedToThisOrNeighbor:
		// The upgrade in context, on the stretch of battleline this card occupies —
		// the same battleline mark a neighbor count draws, read as "here or beside".
		return []glyph{
			inContext(glyph{asset: "type-upgrade"}),
			{asset: "glyph-flank", decor: decorThis},
		}, true
	case engine.ItIsAmong:
		// The card in context, then the set it must be one of, drawn as the Target
		// itself draws. Which refinement narrows the set ("the most powerful") is a
		// Refinement, a family with no glyphs yet, so it stays in the rules text.
		return []glyph{inContext(glyph{asset: "type-creature"}), targetGlyph(v.Target)}, true
	case engine.CardsInPlay:
		return cardsInPlayConditionGlyphs(v), true

	// The source card: where it stands and how it stands.
	case engine.OnFlank:
		// Which flank — either end, the left, the right — is not drawn, for the reason
		// durationGlyph does not draw which turn a window lifts at: the mark says the
		// card is on a flank, and saying which end would need a second picture.
		return []glyph{{asset: "glyph-flank", decor: flankSubjectDecor(v.OfIt)}}, true
	case engine.SourceInCenterOfBattleline:
		// Equal creatures to the left and to the right, which is what a center is. It
		// is deliberately not the flank mark behind a refusal mark: "not on a flank"
		// is met by three of the five positions in a five-creature line, and only one
		// of them is the center.
		return []glyph{
			creaturesCounted(0),
			{asset: "glyph-flank", decor: decorThis},
			creaturesCounted(0),
		}, true
	case engine.SourceHasNoNeighbor:
		return append(
			neighborCountGlyphs(decorThis, houseFilterGlyphs(v.House)...),
			glyph{asset: "glyph-ban"},
		), true
	case engine.SourceReady:
		// Ready is drawn as the exhausted mark tinted friendly, which is how the ready
		// effect draws it; the self marker says it is this card that is ready.
		return []glyph{{asset: "exhausted", decor: decorFriendly | decorThis}}, true
	case engine.SourceIsFighting:
		return []glyph{{asset: "glyph-fight", decor: decorThis}}, true
	case engine.CountersOnThisAtLeast:
		return countedThreshold(
			engine.CountersOnThis{Kind: v.Kind},
			engine.AtLeast,
			v.N,
		)

	// The tide, high and low being two readings of one mark.
	case engine.TideIsHigh:
		return tideGlyphs(sigilGain), true
	case engine.TideIsLow:
		return tideGlyphs(sigilLose), true

	// The turn so far.
	case engine.ItIsYourTurn:
		return []glyph{{asset: "phase-turn", decor: decorFriendly}}, true
	case engine.ChoseHouse:
		return choseHouseGlyphs(v)
	case engine.CardsPlayed:
		return countedThreshold(v, engine.AtLeast, atLeastOne(v.Amount))
	case engine.FirstCreaturePlayedThisTurn:
		// The card in context is the one creature played this turn, so the turn tally
		// wears an "=1" on it rather than on the stacked noun an open tally counts.
		return threshold(thisTurn([]glyph{
			inContext(glyph{asset: "type-creature"}),
			{asset: "glyph-play"},
		}), engine.Exactly, 1), true
	case engine.NoCreaturesPlayedThisTurn:
		return append(thisTurn([]glyph{
			creaturesCounted(0),
			{asset: "glyph-play"},
		}), glyph{asset: "glyph-ban"}), true
	case engine.CreatureDestroyedThisTurn:
		return turnTally(creatureDestroyedStat(v.Player), engine.AtLeast)
	case engine.UsedCreatureToReap:
		return turnTally(engine.CreaturesReapedThisTurn, engine.AtLeast)
	case engine.UsedCreatureToFight:
		return turnTally(engine.CreaturesFoughtThisTurn, engine.AtLeast)
	case engine.FirstReapOfTurn:
		// The same tally UsedCreatureToReap reads, against a different figure: one
		// creature has reaped, so this reap is the first.
		return turnTally(engine.CreaturesReapedThisTurn, engine.Exactly)
	case engine.UsedNoCreatures:
		gs, _ := countGlyphs(engine.CreaturesUsed{Player: engine.Controller})
		return append(gs, glyph{asset: "glyph-ban"}), true
	case engine.SourceFirstUseThisTurn:
		// The uses of this card rather than of any creature, so the self marker
		// replaces the stacked noun the board-wide tally counts.
		return threshold(thisTurn([]glyph{
			{asset: "type-creature", decor: decorThis},
			{asset: "glyph-action"},
		}), engine.Exactly, 1), true
	case engine.ArchivedCreaturesShareHouse:
		return archivedShareHouseGlyphs(), true
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

// inContext marks a glyph as being about the card an effect or trigger just put
// in context — the "it" the whole It… family asks its question of. It is the mark
// that turns a bare noun or property glyph into a test on one card, so the family
// reads as one grammar rather than as a dozen hand-decorated cases.
func inContext(g glyph) glyph {
	g.decor |= decorChosen
	return g
}

// flankSubjectDecor marks whose position a flank condition reads: the card in
// context, or the card the ability is printed on.
func flankSubjectDecor(ofIt bool) decor {
	if ofIt {
		return decorChosen
	}
	return decorThis
}

// hauntedDiscardSize is the discard pile a player is haunted at, mirroring the
// engine's own rule (Haunted.Met) so the badge cannot print a figure the rule
// never reads.
const hauntedDiscardSize = 10

// tideGlyphs draws the tide, with the sigil saying which way it stands: the
// increase mark for high, the decrease mark for low. High and low are one state
// seen from the two sides of the table — the tide high for you is low for your
// opponent — so they are one mark with two readings rather than two drawings
// (TestTideReadsOneMarkTwoWays).
func tideGlyphs(s sigil) []glyph {
	return []glyph{{asset: "glyph-tide", sigil: s}}
}

// turnTally draws a condition that reads one of the engine's turn tallies: the
// tally's own picture from the Count family, wearing the figure the condition
// compares it against — "≥1" for "you have reaped this turn", "=1" for "this is
// the first reap this turn".
func turnTally(of engine.TurnStat, is engine.Comparison) ([]glyph, bool) {
	return countedThreshold(engine.TurnCount{Player: engine.Controller, Of: of}, is, 1)
}

// creatureDestroyedStat is the turn tally a creature-destroyed condition reads,
// mirroring the engine's own side-to-stat pairing: a friendly creature for the
// controller, an enemy one for the opponent. The side is drawn by the tally's own
// tint rather than by this condition, so the two cannot disagree.
func creatureDestroyedStat(p engine.Player) engine.TurnStat {
	if p == engine.Opponent {
		return engine.EnemyCreaturesDestroyed
	}
	return engine.FriendlyCreaturesDestroyed
}

// choseHouseGlyphs draws "you choose Dis as your active house": the active-house
// mark the house vocabulary already uses, then the emblem of the house the ability
// watches for. A condition naming no house is a half-built value no card can
// carry, and it draws the unknown glyph rather than a bare turn marker.
func choseHouseGlyphs(c engine.ChoseHouse) ([]glyph, bool) {
	a := houseIconName(c.House)
	if a == "" {
		return fallbackGlyphs(), false
	}
	return []glyph{{asset: "phase-turn"}, {asset: a}}, true
}

// archivedShareHouseGlyphs draws "if those creatures share a house": the houses
// counted among the creatures a preceding archive set aside, wearing an "=1",
// then the pile they were set aside into. Sharing a house is one house among
// them, so it is the house count the Count family already draws rather than a
// picture of its own.
func archivedShareHouseGlyphs() []glyph {
	gs := housesAmongGlyphs(glyph{asset: "type-creature", decor: decorEach | decorChosen})
	return append(threshold(gs, engine.Exactly, 1), glyph{asset: "zone-archives"})
}
