package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the zone-movement family
// (ADR 0047) — drawing, archiving, discarding, purging, shuffling and putting a
// card into another zone. The "under" cousins of a few of these cases (PlayCardUnder,
// ArchiveCardUnder, PutUnderFromHand, PutUnderIntoPlay) live in the board family
// instead (icon_board.go): they read as playing/burying a card face-down, the
// board's own idiom, rather than a plain zone move.

// zoneEffectGlyphs claims every card-zone-movement effect.
func zoneEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.Draw:
		return []glyph{{asset: "zone-hand", qty: v.Amount}}, true, true
	case engine.ArchiveCard:
		if v.Zone == engine.Purged {
			return []glyph{
				{asset: "zone-purge"},
				arrowTo(glyph{asset: "zone-archives"}),
			}, true, true
		}
		return []glyph{{asset: "zone-archives", qty: engine.FixedCardCount(v.Quantity)}}, true, true
	case engine.ArchiveFromPlay:
		return []glyph{{asset: "zone-archives"}}, true, true
	case engine.ArchiveSource:
		return []glyph{{asset: "zone-archives", decor: decorThis}}, true, true
	case engine.ArchiveGrantingUpgrade:
		return []glyph{
			{asset: "card-back", decor: decorThis},
			arrowTo(glyph{asset: "zone-archives"}),
		}, true, true
	case engine.ArchiveDiscardedThisWay:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-archives"})}, true, true
	case engine.DiscardArchives:
		return []glyph{{asset: "zone-archives"}, arrowTo(glyph{asset: "zone-discard"})}, true, true
	case engine.DiscardCard:
		return []glyph{{asset: "zone-discard"}}, true, true
	case engine.DiscardHand:
		h := glyph{asset: "zone-hand"}
		if v.Player == engine.EachPlayer {
			h.decor = decorEach
		}
		return []glyph{h, arrowTo(glyph{asset: "zone-discard"})}, true, true
	case engine.DiscardTop:
		return []glyph{
			{asset: "zone-discard", qty: v.Amount, decor: playerDecor(v.Player)},
		}, true, true
	case engine.DiscardUntil:
		return []glyph{{asset: "zone-deck"}, {asset: "zone-discard"}}, true, true
	case engine.DiscardFromOpponent:
		zone := "zone-deck"
		if len(v.Sources) > 0 && v.Sources[0] == engine.Archives {
			zone = "zone-archives"
		}
		return []glyph{
			{asset: zone, decor: decorEnemy},
			arrowTo(glyph{
				asset: "zone-discard",
				decor: decorEnemy,
			}),
		}, true, true
	case engine.PurgeCard:
		src := glyph{asset: "zone-hand"}
		switch v.Zones[0] {
		case engine.Discard:
			src.asset = "zone-discard"
		case engine.Archives:
			src.asset = "zone-archives"
		}
		if _, each := v.Selection.(engine.Each); each {
			src.decor = decorEach
		}
		return []glyph{src, arrowTo(glyph{asset: "zone-purge"})}, true, true
	case engine.PurgeArchivedCardThen:
		gs := []glyph{{asset: "zone-purge"}}
		more, _ := effectGlyphs(v.Then)
		return append(gs, more...), true, true
	case engine.Shuffle:
		return []glyph{{asset: "zone-deck"}}, true, true
	case engine.ShuffleIntoDeck:
		// A multi-zone shuffle has no single source glyph, so it shows what it takes.
		src := glyph{asset: "zone-discard"}
		if len(v.From) > 1 {
			src = glyph{asset: "type-creature"}
		}
		return []glyph{src, arrowTo(glyph{asset: "zone-deck"})}, true, true
	case engine.SwapDeckAndDiscard:
		return []glyph{
			{asset: "zone-deck"},
			{asset: "glyph-swap"},
			{asset: "zone-discard"},
		}, true, true
	case engine.RefillHand:
		h := glyph{asset: "zone-hand"}
		if v.Player == engine.EachPlayer {
			h.decor = decorEach
		}
		return []glyph{{asset: "zone-deck"}, arrowTo(h)}, true, true
	case engine.PutItIntoHand:
		return []glyph{{asset: "glyph-return"}}, true, true
	case engine.PutFromPlay:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: a})}, true, true
		}
		return fallbackGlyphs(e), false, true
	case engine.PutChosen:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{targetGlyph(v.Target), arrowTo(glyph{asset: a})}, true, true
		}
		return fallbackGlyphs(e), false, true
	case engine.PutCard:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: a})}, true, true
		}
		return fallbackGlyphs(e), false, true
	case engine.PutFromHand:
		return []glyph{{asset: "zone-hand"}, arrowTo(glyph{asset: "glyph-play"})}, true, true
	case engine.PutDiscardedIntoHand:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-hand"})}, true, true
	case engine.PutNextTacticIntoHand:
		return []glyph{{asset: "type-tactic"}, arrowTo(glyph{asset: "zone-hand"})}, true, true
	default:
		return nil, false, false
	}
}
