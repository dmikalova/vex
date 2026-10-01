package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the shared noun and
// decor vocabulary — a Target's noun glyph, the zone glyphs a card moves
// between, the result-gate arrow, and the enemy/friendly/each/chosen tint an
// economy glyph carries.

// traitGlyph is the mark an effect draws for a trait it names — the trait tag,
// with *which* trait left to the rules text beside the strip. engine.Trait has 70
// members (Agent through Wolf), so one abstract mark is the vocabulary rather than
// seventy drawings: the strip says this is about a kind of creature, and the text
// says which kind.
//
// It is the trait counterpart of nameGlyph, and the two are drawn apart
// (TestTraitAndNameMarksAreDistinct).
func traitGlyph() glyph {
	return glyph{asset: "glyph-trait"}
}

// nameGlyph is the mark an effect draws for a card it names by its printed name —
// the nameplate, with *which* card left to the rules text. The value space is
// every card in the pool, which is a second card database rather than a
// vocabulary, so it draws one mark the way traitGlyph does.
//
// It is deliberately not the unknown glyph: an unknown says the pass failed to map
// a mechanic, where this says the mechanic is mapped and its one parameter is
// unprintable in icons.
func nameGlyph() glyph {
	return glyph{asset: "glyph-name"}
}

// targetGlyph renders a Target as its noun glyph plus the decorations that carry
// its enemy/friendly/each/chosen shape. Fine filters (power, house, trait) stay
// in the rules text; the strip summarises the noun.
//
// Every kind draws a noun — TestTargetGlyphCoversEveryTargetKind holds that — so
// the fall-through is reached only by a Target with no kind, a half-built value no
// card can carry, and it draws the abstract card rather than the target's printed
// text. The strip speaks no words (ADR 0022), so the text chip the fall-through
// used to draw was the one place it broke that rule.
func targetGlyph(t engine.Target) glyph {
	switch t.Kind {
	case engine.TargetThisCreature:
		return glyph{
			asset: "type-creature",
			decor: decorThis,
		}
	case engine.TargetTriggeringCreature, engine.TargetTheOtherCreature,
		engine.TargetTheChosenCreature, engine.TargetCreatureFought,
		engine.TargetTheFoughtCreature, engine.TargetTheSameCreature,
		engine.TargetAttachedHost:
		return glyph{
			asset: "type-creature",
			decor: decorChosen,
		}
	case engine.TargetEachCreature:
		return glyph{
			asset: "type-creature",
			decor: decorEach,
		}
	case engine.TargetEachNeighbor, engine.TargetFormerNeighbors:
		// The strip carries "each creature"; which creatures are neighbors stays in
		// the rules text, the way the trigger pass folds AfterNeighborFights into the
		// plain fight glyph.
		return glyph{
			asset: "type-creature",
			decor: decorEach,
		}
	case engine.TargetEachUpgradeOnThis:
		return glyph{
			asset: "type-upgrade",
			decor: decorEach | decorThis,
		}
	case engine.TargetGrantingCard:
		return glyph{
			asset: "card-back",
			decor: decorChosen,
		}
	case engine.TargetEachFriendlyCreature:
		return glyph{
			asset: "type-creature",
			decor: decorEach | decorFriendly,
		}
	case engine.TargetEachEnemyCreature:
		return glyph{
			asset: "type-creature",
			decor: decorEach | decorEnemy,
		}
	case engine.TargetChosenCreature:
		return glyph{
			asset: "type-creature",
			decor: decorChosen,
		}
	case engine.TargetChosenFriendlyCreature:
		return glyph{
			asset: "type-creature",
			decor: decorChosen | decorFriendly,
		}
	case engine.TargetChosenEnemyCreature:
		return glyph{
			asset: "type-creature",
			decor: decorChosen | decorEnemy,
		}
	case engine.TargetEachArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorEach,
		}
	case engine.TargetEachFriendlyArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorEach | decorFriendly,
		}
	case engine.TargetEachEnemyArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorEach | decorEnemy,
		}
	case engine.TargetChosenArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorChosen,
		}
	case engine.TargetChosenFriendlyArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorChosen | decorFriendly,
		}
	case engine.TargetChosenEnemyArtifact:
		return glyph{
			asset: "type-artifact",
			decor: decorChosen | decorEnemy,
		}
	case engine.TargetChosenUpgrade:
		return glyph{
			asset: "type-upgrade",
			decor: decorChosen,
		}
	case engine.TargetEachCardInPlay, engine.TargetEachFriendlyCardInPlay,
		engine.TargetChosenCreatureOrArtifact,
		engine.TargetChosenFriendlyCreatureOrArtifact,
		engine.TargetChosenEnemyCreatureOrArtifact:
		return glyph{
			asset: "card-back",
			decor: cardInPlayDecor(t.Kind),
		}
	}
	return glyph{asset: "card-back"}
}

// houseMatcherGlyphs renders which houses a matcher admits (ADR 0038). Only
// MatchNamedHouse has a picture of its own — one of the nine house emblems — so
// every other kind is drawn as a mark on the one abstract house glyph rather than
// as a drawing of its own: house-any is a house the matcher does not name, and
// the kind says which unnamed house it is. The chosen house wears the chosen
// outline, the loop's house the each stack, the active house is preceded by the
// turn marker ("this turn's house"), and the contextual house by the card it
// reads ("that card's house"). MatchExceptHouse draws the emblem it excludes
// behind the refusal mark, the way every other barred thing is drawn.
//
// TestHouseMatcherGlyphsCoverEveryKind walks engine.HouseMatchKinds() and fails on
// a kind that draws nothing, so an eighth kind cannot ship invisible.
func houseMatcherGlyphs(m engine.HouseMatcher) []glyph {
	// A named matcher with no house is rejected by the engine, so the emblem is
	// only missing when a half-built matcher reaches the gallery; it draws the
	// abstract house rather than an empty slot.
	named := houseIconName(m.House)
	if named == "" {
		named = "house-any"
	}
	switch m.Kind {
	case engine.MatchNamedHouse:
		return []glyph{{asset: named}}
	case engine.MatchExceptHouse:
		return []glyph{{asset: named}, {asset: "glyph-ban"}}
	case engine.MatchChosenHouse:
		return []glyph{{asset: "house-any", decor: decorChosen}}
	case engine.MatchEachHouse:
		return []glyph{{asset: "house-any", decor: decorEach}}
	case engine.MatchActiveHouse:
		return []glyph{{asset: "phase-turn"}, {asset: "house-any"}}
	case engine.MatchContextualHouse:
		return []glyph{{asset: "card-back", decor: decorChosen}, {asset: "house-any"}}
	}
	return []glyph{{asset: "house-any"}}
}

// houseChoiceGlyphs renders a referenced house — the house chosen this turn, the
// active house, or the house of the card in context — by mapping it onto the
// matcher kind that says the same thing, so the two house vocabularies draw one
// picture rather than two. AnyHouse names no house at all and so draws nothing.
func houseChoiceGlyphs(h engine.HouseChoice) []glyph {
	switch h {
	case engine.TheChosenHouse:
		return houseMatcherGlyphs(engine.HouseMatcher{Kind: engine.MatchChosenHouse})
	case engine.TheActiveHouse:
		return houseMatcherGlyphs(engine.HouseMatcher{Kind: engine.MatchActiveHouse})
	case engine.TheContextualHouse:
		return houseMatcherGlyphs(engine.HouseMatcher{Kind: engine.MatchContextualHouse})
	}
	return nil
}

// zoneAsset is the glyph a zone draws when a card reads it in place — the pile a
// count measures — where destinationGlyph draws the zone a card is moved to. It
// is "" for the two zones no count can name (InPlay, which is the board rather
// than a pile, and the unset zero value), so a caller guards the empty stem
// rather than drawing an assetless glyph.
func zoneAsset(z engine.Zone) string {
	switch z {
	case engine.Deck:
		return "zone-deck"
	case engine.Hand:
		return "zone-hand"
	case engine.Archives:
		return "zone-archives"
	case engine.Discard:
		return "zone-discard"
	case engine.Purged:
		return "zone-purge"
	}
	return ""
}

// counterAsset maps a generic counter kind to its icon-strip asset, or "" for a
// kind with no icon yet — TestCounterIconNamesHaveAssets fails on the "" so every
// counter ships its own unique SVG.
func counterAsset(kind engine.CounterKind) string {
	switch kind {
	case engine.CounterDoom:
		return "generic-counter-doom"
	case engine.CounterFuse:
		return "generic-counter-fuse"
	case engine.CounterGrowth:
		return "generic-counter-growth"
	case engine.CounterGlory:
		return "generic-counter-glory"
	case engine.CounterDisruption:
		return "generic-counter-disruption"
	case engine.CounterScheme:
		return "generic-counter-scheme"
	case engine.CounterWarrant:
		return "generic-counter-warrant"
	default:
		return ""
	}
}

// destinationGlyph is the zone glyph a card is moved to, or "" for a destination
// with no drawn zone yet. The three deck spots share the one deck glyph.
func destinationGlyph(d engine.Destination) string {
	switch d {
	case engine.ToHand:
		return "zone-hand"
	case engine.ToTopOfDeck, engine.ToBottomOfDeck, engine.ToDeckShuffled:
		return "zone-deck"
	case engine.ToArchives, engine.ToArchives.Yours():
		return "zone-archives"
	}
	return ""
}

// deckDestZone is the zone glyph a revealed deck card is moved to, or "" for a
// routing with no drawn zone yet — TestDeckDestZoneCoversEveryDeckDest fails on
// the "", which an empty asset would otherwise drop silently at render.
//
// IntoBottomOfDeck is a position in the deck the cards came from rather than
// another zone, so it draws the deck under the line its cards go beneath rather
// than a second deck: the leftover of a look stays on top, and the routed cards
// go under it.
func deckDestZone(d engine.DeckDest) string {
	switch d {
	case engine.IntoHand:
		return "zone-hand"
	case engine.IntoArchives:
		return "zone-archives"
	case engine.IntoDiscard:
		return "zone-discard"
	case engine.IntoPurge:
		return "zone-purge"
	case engine.IntoBottomOfDeck:
		return "zone-deck-bottom"
	}
	return ""
}

// arrowTo marks a glyph as the result of a result-gate arrow, so a target reads
// as "→ <target>".
func arrowTo(g glyph) glyph {
	g.punct = punctArrow
	return g
}

// cardInPlayDecor picks the enemy/friendly/each shading for the "card in play"
// target kinds that share one noun glyph.
func cardInPlayDecor(k engine.TargetKind) decor {
	switch k {
	case engine.TargetEachCardInPlay:
		return decorEach
	case engine.TargetEachFriendlyCardInPlay:
		return decorEach | decorFriendly
	case engine.TargetChosenFriendlyCreatureOrArtifact:
		return decorChosen | decorFriendly
	case engine.TargetChosenEnemyCreatureOrArtifact:
		return decorChosen | decorEnemy
	default:
		return decorChosen
	}
}

// playerDecor tints an economy glyph by who it acts on: an opponent-facing effect
// reads as enemy, an each-player effect stacks.
func playerDecor(p engine.Player) decor {
	switch p {
	case engine.Opponent:
		return decorEnemy
	case engine.EachPlayer:
		return decorEach
	default:
		return 0
	}
}

// verbGlyphs renders the verbs a chosen-creature effect applies, in order. It
// stays here rather than in a family file because more than one family
// (creature-state, composition) reads it.
func verbGlyphs(verbs []engine.CreatureVerb) []glyph {
	gs := make([]glyph, 0, len(verbs))
	for _, verb := range verbs {
		more, _ := creatureVerbGlyphs(verb)
		gs = append(gs, more...)
	}
	return gs
}

// creatureVerbGlyphs renders one CreatureVerb, reporting whether a real mapping
// covered it. Each verb reuses the glyph the action draws on its own, so "ready a
// creature" and the ready step of a turn draw one picture.
//
// The CreatureVerb family is gated (TestEveryCreatureVerbDraws, TestIconTotality),
// so every row here draws. The default arm is what makes an eighth verb fail that
// gate: without it an unlisted verb contributed no glyph at all, and the strip
// quietly said the creature was chosen and nothing happened to it.
func creatureVerbGlyphs(verb engine.CreatureVerb) ([]glyph, bool) {
	switch v := verb.(type) {
	case engine.ReadyVerb:
		return []glyph{{asset: "exhausted", decor: decorFriendly}}, true
	case engine.ReapVerb:
		return []glyph{{asset: "glyph-reap"}}, true
	case engine.FightVerb:
		return []glyph{{asset: "glyph-fight"}}, true
	case engine.UseVerb:
		return []glyph{{asset: "glyph-action"}}, true
	case engine.StunVerb:
		return []glyph{{asset: "stun"}}, true
	case engine.ExhaustVerb:
		return []glyph{{asset: "exhausted"}}, true
	case engine.GainKeywordVerb:
		if a := keywordIcon(v.Keyword); a != "" {
			return []glyph{{asset: a}}, true
		}
		return fallbackGlyphs(), false
	}
	return fallbackGlyphs(), false
}
