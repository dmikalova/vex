package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the count family
// (ADR 0047) — the engine's Count nodes, the numbers an effect scales by ("for
// each friendly creature, gain 1 Æmber") and a condition compares against.
//
// A count is not an effect, so it has no place in the glyphFamilies chain: a
// count is drawn where the node holding it is drawn, and until this file existed
// a CountIs or a Per clause rendered as the effect it gated, saying nothing about
// what was being counted.
//
// The grammar is one rule — a count draws as the noun it totals — with three
// decorations carrying what kind of total it is:
//
//   - A count over a set of cards stacks its noun (decorEach), because a count is
//     "each of these" by construction, and tints it with the side it surveys.
//   - A count that reads the card holding it marks its noun decorThis, the way
//     every other self-referring glyph does; one that reads the card in context
//     marks it decorChosen.
//   - A threshold inside a count rides the noun's numeral badge with its
//     comparison mark — CardsInPlay{MinPower: 5} is the power noun wearing "≥5",
//     not a drawing of its own (ADR 0047).
//
// Fine filters that only narrow which cards qualify (ready, damaged, with Æmber
// on it, other than this one) stay in the rules text beside the strip, as they do
// for a Target: the strip summarises the noun.
//
// Fixed is the one plumbing row and the census declares it silent — a bare
// literal prints nothing before the noun it counts — so it folds and draws
// nothing, the way a plumbing effect folds into its children.
//
// Two counts asking one question of different subjects are one drawing with two
// subjects, not two drawings: the houses represented among every card in play and
// among a player's creatures share housesAmongGlyphs, and every "this way" tally
// a ProducedThisWay names draws what its dedicated count draws, tinted with whose
// share it reads.

// countGlyphs renders a Count as the noun it totals, reporting whether a real
// mapping covered it. A count with no mapping draws the abstract unknown glyph,
// which is what the icon gate fails on: the Count family is gated
// (TestEveryCountDraws, TestIconTotality), so every row here draws.
func countGlyphs(c engine.Count) ([]glyph, bool) {
	switch v := c.(type) {
	case engine.Fixed:
		return nil, true

	// The Æmber economy.
	case engine.AemberInPool:
		return []glyph{{asset: "aember", decor: playerDecor(v.Player)}}, true
	case engine.AemberOnThis:
		return []glyph{{asset: "aember", decor: decorThis}}, true
	case engine.AemberOnFriendlyCreatures:
		return []glyph{{asset: "aember"}, creaturesCounted(decorFriendly)}, true
	case engine.AemberStolenThisEvent:
		// The theft that fired the trigger took this Æmber from the controller's own
		// pool, so the noun is untinted; the steal verb is what it just suffered.
		return []glyph{{asset: "aember", sigil: sigilSteal}}, true

	// Keys. A forged key draws the forge, an unforged one the bare key, which is
	// the same pair the player bar already shows on the board.
	case engine.ForgedKeys:
		return []glyph{{asset: "forge", decor: playerDecor(v.Player)}}, true
	case engine.UnforgedKeys:
		return []glyph{{asset: "key", decor: playerDecor(v.Player)}}, true

	// The board: what is in play, and how the sides compare.
	case engine.CardsInPlay:
		return cardsInPlayCountGlyphs(v), true
	case engine.ExcessCreatures:
		return excessCreaturesGlyphs(v), true
	case engine.ArtifactsInPlay:
		return []glyph{{asset: "type-artifact", decor: decorEach}}, true
	case engine.UpgradesOn:
		return []glyph{
			{asset: "type-upgrade", decor: decorEach},
			targetGlyph(v.Target),
		}, true
	case engine.HousesInPlay:
		// Every card in play, both sides, minus one excepted house — the emblem
		// behind the refusal mark, the way every other barred thing is drawn.
		gs := housesAmongGlyphs(cardsCounted(0))
		if a := houseIconName(v.Except); a != "" {
			gs = append(gs, glyph{asset: a}, glyph{asset: "glyph-ban"})
		}
		return gs, true
	case engine.HousesAmong:
		return housesAmongGlyphs(glyph{
			asset: countedCardAsset(v.Type),
			decor: decorEach | playerDecor(v.Player),
		}), true

	// One card's own measurements.
	case engine.DamageOnThis:
		return []glyph{{asset: "damage", decor: decorThis}}, true
	case engine.CountersOnThis:
		if a := counterAsset(v.Kind); a != "" {
			return []glyph{{asset: a, decor: decorThis}}, true
		}
		return fallbackGlyphs(), false
	case engine.PowerCountersOnThis:
		return []glyph{{asset: "power-counter-plus", decor: decorThis}}, true
	case engine.PowerOfChosen:
		// The Of fraction ("half its power, rounded down") is a Loss, a family with
		// no glyphs yet, so it stays in the rules text beside the strip.
		return []glyph{{asset: "power", decor: decorChosen}}, true
	case engine.TraitsOfChosen:
		return []glyph{
			{asset: "glyph-trait", decor: decorEach},
			{asset: "type-creature", decor: decorChosen},
		}, true
	case engine.BonusIconsOf:
		over, covered := bonusIconSubjectGlyphs(v.Over)
		return append(
			[]glyph{{asset: bonusCountAsset(v.Kind), decor: decorEach}},
			over...,
		), covered

	// The battleline around a card.
	case engine.NeighborsOfThis:
		return neighborCountGlyphs(decorThis), true
	case engine.NeighborsMatching:
		// The neighbors of the card in context, not of this one, so the battleline
		// mark wears the chosen outline rather than the self marker.
		return neighborCountGlyphs(decorChosen, houseFilterGlyphs(v.House)...), true
	case engine.CombinedPowerOfNeighborsWithout:
		// A measure of the same set rather than a count of it, so the power noun
		// leads the neighbors it totals, as it does for the destroyed-power tally.
		return append(
			[]glyph{{asset: "power"}},
			neighborCountGlyphs(decorThis, traitGlyph(), glyph{asset: "glyph-ban"})...,
		), true

	// Zones: how many cards a pile holds.
	case engine.CardsInHand:
		return append(
			houseChoiceGlyphs(v.House),
			glyph{asset: "zone-hand", decor: decorEach | playerDecor(v.Player)},
		), true
	case engine.CardsInZone:
		if a := zoneAsset(v.Zone); a != "" {
			return []glyph{{asset: a, decor: decorEach | playerDecor(v.Player)}}, true
		}
		return fallbackGlyphs(), false
	case engine.CopiesInDiscard:
		// A copy of this card by name, so the nameplate wears the self marker: the
		// pile is searched for one printed name rather than for any card.
		return []glyph{
			{asset: "glyph-name", decor: decorEach | decorThis},
			{asset: "zone-discard"},
		}, true
	case engine.PurgedCards:
		// Every card in the purge pile, which is the pile itself — the same shape
		// CardsInZone draws for a pile, where CardsPurged counts one purge's tally.
		return []glyph{{asset: "zone-purge", decor: decorEach}}, true

	// The turn so far.
	case engine.CardsPlayed:
		return thisTurn(append(
			houseFilterGlyphs(v.House),
			glyph{asset: "glyph-play", decor: decorEach | playerDecor(v.Player)},
		)), true
	case engine.CreaturesUsed:
		return thisTurn([]glyph{
			creaturesCounted(playerDecor(v.Player)),
			{asset: "glyph-action"},
		}), true
	case engine.TurnCount:
		if gs := turnStatGlyphs(v.Of); len(gs) > 0 {
			return thisTurn(tintSubject(gs, v.Player)), true
		}
		return fallbackGlyphs(), false

	// Tallies an earlier effect in the same resolution left behind.
	case engine.CardsDestroyed:
		return destroyedThisWayGlyphs(cardsCounted(0)), true
	case engine.CreaturesDestroyed:
		return destroyedThisWayGlyphs(creaturesCounted(0)), true
	case engine.PowerDestroyedThisWay:
		// The power the destruction took off the board rather than its head count,
		// so the measure leads the same noun the head count draws.
		return append(
			[]glyph{{asset: "power"}},
			destroyedThisWayGlyphs(creaturesCounted(0))...,
		), true
	case engine.CardsShuffledIntoDeck:
		return shuffledHomeGlyphs(0), true
	case engine.ProducedThisWay:
		return producedTallyGlyphs(v)
	case engine.CardsPurged:
		return []glyph{
			{asset: countedCardAsset(v.Type), decor: decorEach},
			{asset: "zone-purge"},
		}, true
	case engine.CardsRevealed:
		revealed := cardsCounted(0)
		revealed.sigil = sigilReveal
		return []glyph{revealed}, true
	case engine.CreaturesHealed:
		return []glyph{creaturesCounted(0), {asset: "glyph-heal"}}, true
	case engine.DamageHealed:
		return []glyph{{asset: "damage"}, {asset: "glyph-heal"}}, true
	case engine.DamagePrevented:
		// The damage the armor absorbed, so the damage noun leads the armor that
		// spent itself stopping it.
		return []glyph{{asset: "damage"}, {asset: "shield"}}, true
	}
	return fallbackGlyphs(), false
}

// cardsInPlayCountGlyphs draws the one count that carries filters of its own: the
// cards a player has in play, narrowed by house, trait, name, and a power floor.
// The narrowing marks lead the noun, the way a house matcher leads the verbs it
// frees (mayPlayOrUseGlyphs), and the power floor trails it as a threshold badge
// rather than as a second noun.
//
// The Amount and None fields belong to this node's Condition role rather than its
// Count role — "if there are 2 or more friendly creatures in play" — so the
// Condition family draws them; a count ignores them.
func cardsInPlayCountGlyphs(e engine.CardsInPlay) []glyph {
	gs := append(houseFilterGlyphs(e.House), glyph{
		asset: countedCardAsset(e.Type),
		decor: decorEach | playerDecor(e.Player),
	})
	if e.Trait != 0 {
		gs = append(gs, traitGlyph())
	}
	if e.Name != "" {
		gs = append(gs, nameGlyph())
	}
	if e.MinPower > 0 {
		gs = append(gs, glyph{
			asset:      "power",
			qty:        e.MinPower,
			comparison: engine.AtLeast,
		})
	}
	return gs
}

// excessCreaturesGlyphs draws how many more creatures one side has than the
// other: the creature noun tinted with the side holding the excess, wearing the
// relative comparison's ">" on its badge. The comparison is the engine's own, so
// which side is ahead is said once — by the tint — rather than twice (ADR 0047).
// It is the same mark the forged-key lead condition wears, so the two borrow one
// helper (aheadOfOther, icon_condition.go).
//
// NotCountingSelf shifts the total by one without changing what is compared, so
// like the other fine filters it stays in the rules text beside the strip.
func excessCreaturesGlyphs(e engine.ExcessCreatures) []glyph {
	noun := creaturesCounted(playerDecor(e.Player))
	noun.comparison = aheadOfOther(e.Player)
	gs := []glyph{noun}
	if e.Trait != 0 {
		gs = append(gs, traitGlyph())
	}
	return gs
}

// housesAmongGlyphs draws a count of the houses represented among a set of
// cards: the abstract house, stacked because it is the houses that are counted,
// then the cards surveyed for them. The two house counts differ only in that
// subject — every card in play for HousesInPlay, a side and a type for
// HousesAmong — so they are one drawing with two subjects rather than two.
func housesAmongGlyphs(among glyph) []glyph {
	return []glyph{{asset: "house-any", decor: decorEach}, among}
}

// bonusCountAsset is the icon a bonus-icon count totals: the kind's own icon, or
// the abstract bonus icon for a count that totals every kind (Mindfire counts
// each bonus icon on the card it discarded, whatever it is).
func bonusCountAsset(kind engine.BonusIcon) string {
	if a := bonusIconStem(kind); a != "" {
		return a
	}
	return "glyph-bonus"
}

// bonusIconSubjectGlyphs draws the cards a bonus-icon count reads, reporting
// whether a real mapping covered it: the single card an effect just put in
// context, or the whole set its last purge removed. The BonusIconSubject family is
// gated (TestEveryBonusIconSubjectDraws, TestIconTotality), so every row here
// draws.
//
// Which noun a TheCardInContext names the card by — the discarded card, the
// revealed one — is the wording of the clause rather than a different card, so the
// strip draws the one card in context and leaves the noun to the rules text.
func bonusIconSubjectGlyphs(s engine.BonusIconSubject) ([]glyph, bool) {
	switch s.(type) {
	case engine.TheCardInContext:
		return []glyph{{asset: "card-back", decor: decorChosen}}, true
	case engine.ThePurgedCards:
		// The purge this effect just made, not the pile it added to, so the mark takes
		// no stack: the stacked purge pile is what the PurgedCards count draws.
		return []glyph{{asset: "zone-purge"}}, true
	}
	return fallbackGlyphs(), false
}

// destroyedThisWayGlyphs draws a "destroyed this way" tally: the noun that was
// destroyed, followed by the tombstone the leave-play window and the
// "Destroyed:" trigger already draw. The tallies count the same cards and differ
// only in the noun they name them by, so they differ only in that glyph.
func destroyedThisWayGlyphs(noun glyph) []glyph {
	return []glyph{noun, {asset: "glyph-destroyed"}}
}

// shuffledHomeGlyphs draws a "shuffled into the deck this way" tally: the cards
// it moved, then the deck wearing the shuffle verb that says they went back at no
// known position, exactly as the shuffle effect itself draws it.
func shuffledHomeGlyphs(side decor) []glyph {
	return []glyph{cardsCounted(side), {asset: "zone-deck", sigil: sigilShuffle}}
}

// producedTallyGlyphs draws a tally an earlier effect in this resolution left
// behind, reporting whether the tally has a drawing. Each one draws the picture
// its dedicated count draws, tinted with whose share is being read, because the
// two say the same thing about the same event — the tally is the per-player view
// of it.
func producedTallyGlyphs(c engine.ProducedThisWay) ([]glyph, bool) {
	side := playerDecor(c.Player)
	switch c.Tally {
	case engine.TallyCreaturesDestroyed:
		return destroyedThisWayGlyphs(creaturesCounted(side)), true
	case engine.TallyCreaturesShuffledIntoDeck:
		return shuffledHomeGlyphs(side), true
	case engine.TallyAemberLost:
		return []glyph{{asset: "aember", decor: side, sigil: sigilLose}}, true
	case engine.TallyCardsReturned:
		// The whole tally rather than one player's share, so it takes no tint.
		return []glyph{cardsCounted(0), {asset: "glyph-return"}}, true
	case engine.TallyCardsPurged:
		return []glyph{cardsCounted(side), {asset: "zone-purge"}}, true
	}
	return fallbackGlyphs(), false
}

// turnStatGlyphs draws one of the engine's turn tallies as the noun it counts and
// the verb that put it there — a creature and the reap mark for "creature that has
// reaped this turn". It returns no glyphs for a tally with no drawing, which the
// caller renders as the unknown glyph rather than as a bare turn marker.
//
// Which turn a tally reads — the current one or the player's previous one — is not
// drawn, for the reason durationGlyph does not draw which turn a window lifts at:
// the marker says the tally belongs to a turn, and saying which would need a second
// picture per tally. TestTurnStatGlyphsCoverEveryTally holds the eleven to drawing
// something; TurnStat is not a catalogued enum, so that list is local and a twelfth
// tally has to be added to both.
func turnStatGlyphs(s engine.TurnStat) []glyph {
	switch s {
	case engine.KeysForgedThisTurn, engine.KeysForgedLastTurn:
		return []glyph{{asset: "forge", decor: decorEach}}
	case engine.CreaturesPlayedLastTurn:
		return []glyph{creaturesCounted(0), {asset: "glyph-play"}}
	case engine.EnemyCreaturesFightKilled:
		return []glyph{
			creaturesCounted(decorEnemy),
			{asset: "glyph-fight"},
			{asset: "glyph-destroyed"},
		}
	case engine.EnemyCreaturesDestroyed:
		return destroyedThisWayGlyphs(creaturesCounted(decorEnemy))
	case engine.FriendlyCreaturesDestroyed:
		return destroyedThisWayGlyphs(creaturesCounted(decorFriendly))
	case engine.CreaturesReapedThisTurn:
		return []glyph{creaturesCounted(0), {asset: "glyph-reap"}}
	case engine.CreaturesFoughtThisTurn:
		return []glyph{creaturesCounted(0), {asset: "glyph-fight"}}
	case engine.CreaturesUsedThisTurn:
		return []glyph{creaturesCounted(0), {asset: "glyph-action"}}
	case engine.AemberStolenFromThisTurn, engine.AemberStolenFromLastTurn:
		return []glyph{{asset: "aember", sigil: sigilSteal}}
	}
	return nil
}

// neighborCountGlyphs draws the creatures beside a card: the battleline mark,
// decorated with whose position it reads from — this card, or the card in
// context — then any marks narrowing which neighbors count, then the stacked
// creature noun.
func neighborCountGlyphs(of decor, narrowing ...glyph) []glyph {
	gs := append([]glyph{{asset: "glyph-flank", decor: of}}, narrowing...)
	return append(gs, creaturesCounted(0))
}

// creaturesCounted is the creature noun as a count names it: stacked, because a
// count is "each of these" by construction, and tinted with the side it surveys.
func creaturesCounted(side decor) glyph {
	return glyph{asset: "type-creature", decor: decorEach | side}
}

// cardsCounted is the same for a count that names no type: the abstract card,
// stacked and tinted with the side it surveys.
func cardsCounted(side decor) glyph {
	return glyph{asset: "card-back", decor: decorEach | side}
}

// tintSubject tints a run's head with the side a count surveys, leaving a noun
// that already names a side alone: "each enemy creature destroyed this turn" is
// enemy from the tally's own wording, and a second tint would argue with it.
func tintSubject(gs []glyph, p engine.Player) []glyph {
	if len(gs) > 0 && gs[0].decor&(decorEnemy|decorFriendly) == 0 {
		gs[0].decor |= playerDecor(p)
	}
	return gs
}

// countedCardAsset is the noun a count names its cards by: the card type's own
// icon, or the abstract card for a count that takes any type.
func countedCardAsset(t engine.CardType) string {
	if a := typeIconName(t); a != "" {
		return a
	}
	return "card-back"
}

// thisTurn closes a run with the turn marker, which is how a tally of what a
// player has done so far this turn says "this turn" — the counterpart of
// inWindow, which marks the window an effect lifts at.
func thisTurn(gs []glyph) []glyph {
	return append(gs, glyph{asset: "phase-turn"})
}

// houseFilterGlyphs renders a house matcher that actually narrows a count, and
// nothing for the any-house matcher every unfiltered count carries, so a plain
// "for each card" does not draw a house it never names.
func houseFilterGlyphs(m engine.HouseMatcher) []glyph {
	if m.Kind == engine.MatchAnyHouse {
		return nil
	}
	return houseMatcherGlyphs(m)
}
