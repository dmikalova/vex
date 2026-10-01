package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the board family
// (ADR 0047) — playing a card into play (from hand, deck, opponent's zones, or
// face-down under another card), permission grants to play or use out of house,
// taking control of a creature, and the battleline positions a creature can hold.
// Battleline position (Swap, RearrangeBattleline, MoveToFlank, ConsiderFlank,
// MoveWithinBattleline) stays inside this family rather than becoming a six-case
// file of its own.

// boardEffectGlyphs claims every play, control and battleline-position effect.
func boardEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.PlayFrom, engine.PlayTopOfDeck, engine.PutIntoPlay:
		return []glyph{{asset: "glyph-play"}}, true, true
	case engine.PlayFromOpponent:
		zone := "zone-deck"
		if v.From == engine.Archives {
			zone = "zone-archives"
		}
		return []glyph{
			{asset: zone, decor: decorEnemy},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true, true
	case engine.PlayItFromOpponentDiscard:
		return []glyph{
			{asset: "zone-discard", decor: decorEnemy},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true, true
	case engine.PlayRevealedCard:
		return []glyph{{asset: "glyph-play"}}, true, true
	case engine.PlayCardUnder:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	case engine.ArchiveCardUnder:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "zone-archives"})}, true, true
	case engine.PutUnderFromHand:
		return []glyph{{asset: "zone-hand"}, arrowTo(glyph{asset: "card-back"})}, true, true
	case engine.PutUnderIntoPlay:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	case engine.PutDiscardedIntoPlay:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	case engine.EachPlayerPutsHandCreaturesIntoPlay:
		return []glyph{
			{asset: "zone-hand", decor: decorEach},
			arrowTo(glyph{asset: "glyph-play"}),
		}, true, true
	case engine.MayPlayOrUse:
		// A standing permission the controller may decline, where PlayOrUse makes the
		// same play-or-use happen now; the may sigil is the only thing between them.
		return markVerb(mayPlayOrUseGlyphs(v), sigilMay), true, true
	case engine.PlayOrUse:
		// The any-house matcher prints no house word here ("play or use a card"), so
		// it draws no house glyph either; every narrowing matcher leads the verbs.
		var gs []glyph
		if v.House.Kind != engine.MatchAnyHouse {
			gs = houseMatcherGlyphs(v.House)
		}
		if v.Grant == 0 || v.Grant&engine.GrantPlay != 0 {
			gs = append(gs, glyph{asset: "glyph-play"})
		}
		if v.Grant == 0 || v.Grant&engine.GrantUse != 0 {
			gs = append(gs, glyph{asset: "glyph-action"})
		}
		return gs, true, true
	case engine.CannotPlay:
		return []glyph{{asset: "glyph-play"}, {asset: "glyph-ban"}}, true, true
	case engine.PlayersCannotPlay:
		if a := typeIconName(v.Type); a != "" {
			return []glyph{{asset: a}, {asset: "glyph-play"}, {asset: "glyph-ban"}}, true, true
		}
		return []glyph{{asset: "glyph-play"}, {asset: "glyph-ban"}}, true, true
	case engine.AttachSelfTo:
		return []glyph{
			{asset: "card-back", decor: decorThis},
			arrowTo(glyph{asset: "type-creature"}),
		}, true, true
	case engine.TakeControl:
		// The host-creature form (Collar of Subordination) has no Target; it takes
		// this creature, so render the "this creature" noun rather than a blank.
		subject := targetGlyph(v.Target)
		if v.Target == (engine.Target{}) {
			subject = glyph{
				asset: "type-creature",
				decor: decorThis,
			}
		}
		return []glyph{
			subject,
			arrowTo(glyph{
				asset: "type-creature",
				decor: decorFriendly,
			}),
		}, true, true
	case engine.Graft:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "card-back"})}, true, true
	case engine.TriggerGraftedPlayEffect:
		// Only the buried card's "Play:" ability fires; the card itself stays under
		// its host, which is what separates this from PutUnderIntoPlay (ADR 0047).
		return []glyph{
			{asset: "card-back"},
			arrowTo(glyph{asset: "glyph-play", sigil: sigilGrafted}),
		}, true, true
	case engine.Swap:
		return []glyph{
			{asset: "type-creature", decor: decorThis},
			{asset: "glyph-swap"},
			arrowTo(targetGlyph(v.With)),
		}, true, true
	case engine.SwapChosen:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			{asset: "glyph-swap"},
			{asset: "type-creature", decor: decorChosen},
		}, true, true
	case engine.RearrangeBattleline:
		return []glyph{
			{asset: "type-creature", decor: decorChosen},
			{asset: "glyph-swap"},
			{asset: "type-creature", decor: decorChosen},
		}, true, true
	case engine.MoveWithinBattleline:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "glyph-swap"})}, true, true
	case engine.MoveToFlank:
		return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: "glyph-flank"})}, true, true
	case engine.ConsiderFlank:
		return []glyph{{asset: "glyph-flank"}, arrowTo(targetGlyph(v.Target))}, true, true
	default:
		return nil, false, false
	}
}

// mayPlayOrUseGlyphs renders an out-of-house permission grant as the houses it
// frees followed by the verbs it frees — "Mars cards: play", "the chosen house:
// play, use". The houses lead because that is what the grant is about; before the
// house matcher had glyphs the strip drew the verbs alone, so a grant naming Mars
// and one naming the chosen house were the same picture.
//
// A trait-scoped grant (Mutagenic Serum's friendly Mutant creatures) names no
// house, so the trait mark leads instead of a house emblem; which trait it frees
// stays in the rules text, the way every other trait is drawn.
func mayPlayOrUseGlyphs(e engine.MayPlayOrUse) []glyph {
	if e.Trait != 0 {
		return append([]glyph{traitGlyph()}, houseGrantGlyphs(e.Grant)...)
	}
	return append(houseSelectorGlyphs(e.Houses), houseGrantGlyphs(e.Grant)...)
}

// houseSelectorGlyphs renders whose cards a grant frees. The controlled selector
// is the one selection no filter can make — every house the controller has a card
// in play for — so it draws the abstract house tinted friendly.
func houseSelectorGlyphs(s engine.HouseSelector) []glyph {
	if s.Controlled {
		return []glyph{{asset: "house-any", decor: decorFriendly}}
	}
	return houseMatcherGlyphs(s.Match)
}

// houseGrantGlyphs renders the verbs a grant frees, one glyph per freed verb,
// reusing the glyphs those actions draw on their own.
func houseGrantGlyphs(g engine.HouseGrant) []glyph {
	gs := make([]glyph, 0, 3)
	if g&engine.GrantPlay != 0 {
		gs = append(gs, glyph{asset: "glyph-play"})
	}
	if g&engine.GrantUse != 0 {
		gs = append(gs, glyph{asset: "glyph-action"})
	}
	if g&engine.GrantFight != 0 {
		gs = append(gs, glyph{asset: "glyph-fight"})
	}
	return gs
}
