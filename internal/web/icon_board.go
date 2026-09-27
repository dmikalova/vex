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
		return mayPlayOrUseGlyphs(v), true, true
	case engine.PlayOrUse:
		var gs []glyph
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
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
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

// mayPlayOrUseGlyphs renders an out-of-house permission grant, narrowing to the
// verbs and houses its axes select: a fight grant to a fight glyph, an
// artifacts-any-house grant to an artifact-and-action pair, a named-house grant to
// play/action, and an exclusion or controlled grant to play (plus action when it
// also frees use).
func mayPlayOrUseGlyphs(e engine.MayPlayOrUse) []glyph {
	if e.Houses.Controlled {
		if e.Grant&engine.GrantUse != 0 {
			return []glyph{{asset: "glyph-play"}, {asset: "glyph-action"}}
		}
		return []glyph{{asset: "glyph-play"}}
	}
	switch e.Houses.Match.Kind {
	case engine.MatchExceptHouse:
		if e.Grant&engine.GrantUse != 0 {
			return []glyph{{asset: "glyph-play"}, {asset: "glyph-action"}}
		}
		return []glyph{{asset: "glyph-play"}}
	case engine.MatchAnyHouse:
		if e.Grant&engine.GrantFight != 0 {
			return []glyph{{asset: "glyph-fight", decor: decorFriendly | decorEach}}
		}
		return []glyph{
			{asset: "type-artifact", decor: decorFriendly},
			{asset: "glyph-action"},
		}
	default: // MatchNamedHouse, MatchChosenHouse
		if e.Grant == engine.GrantFight {
			d := decorEach
			if e.Houses.Match.House != engine.HouseNone {
				d |= decorFriendly
			}
			return []glyph{{asset: "glyph-fight", decor: d}}
		}
		if e.Grant&engine.GrantPlay != 0 {
			return []glyph{{asset: "glyph-play"}, {asset: "glyph-action"}}
		}
		return []glyph{{asset: "glyph-action", decor: decorFriendly}}
	}
}
