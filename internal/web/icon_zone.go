package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the zone-movement family
// (ADR 0047) — drawing, archiving, discarding, purging, shuffling and putting a
// card into another zone. The "under" cousins of a few of these cases (PlayCardUnder,
// ArchiveCardUnder, PutUnderFromHand, PutUnderIntoPlay) live in the board family
// instead (icon_board.go): they read as playing/burying a card face-down, the
// board's own idiom, rather than a plain zone move.
//
// The Quantity a movement verb carries is drawn here too, beside the verb holding
// it, the way a Count is: it says how many cards the verb takes, where the verb's
// Selection says which. On a verb it is a badge rather than a glyph of its own,
// because "purge 2 cards" is one picture — the numeral the verb already prints,
// plus the may mark when the controller may stop short of it.

// zoneEffectGlyphs claims every card-zone-movement effect.
func zoneEffectGlyphs(e engine.Effect) ([]glyph, bool) {
	switch v := e.(type) {
	case engine.Draw:
		return []glyph{{asset: "zone-hand", qty: v.Amount}}, true
	case engine.ArchiveCard:
		if v.Zone == engine.Purged {
			return []glyph{
				{asset: "zone-purge"},
				arrowTo(glyph{asset: "zone-archives"}),
			}, true
		}
		return []glyph{quantityOn(glyph{asset: "zone-archives"}, v.Quantity)}, true
	case engine.ArchiveFromPlay:
		return []glyph{{asset: "zone-archives"}}, true
	case engine.ArchiveSource:
		return []glyph{{asset: "zone-archives", decor: decorThis}}, true
	case engine.ArchiveGrantingUpgrade:
		return []glyph{
			{asset: "card-back", decor: decorThis},
			arrowTo(glyph{asset: "zone-archives"}),
		}, true
	case engine.ArchiveDiscardedThisWay:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-archives"})}, true
	case engine.DiscardArchives:
		return []glyph{{asset: "zone-archives"}, arrowTo(glyph{asset: "zone-discard"})}, true
	case engine.DiscardCard:
		return []glyph{quantityOn(glyph{asset: "zone-discard"}, v.Quantity)}, true
	case engine.DiscardHand:
		h := glyph{asset: "zone-hand"}
		if v.Player == engine.EachPlayer {
			h.decor = decorEach
		}
		return []glyph{h, arrowTo(glyph{asset: "zone-discard"})}, true
	case engine.DiscardTop:
		return []glyph{
			{asset: "zone-discard", qty: v.Amount, decor: playerDecor(v.Player)},
		}, true
	case engine.DiscardUntil:
		return []glyph{{asset: "zone-deck"}, {asset: "zone-discard"}}, true
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
		}, true
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
		return []glyph{
			quantityOn(src, v.Quantity),
			arrowTo(glyph{asset: "zone-purge"}),
		}, true
	case engine.PurgeArchivedCardThen:
		gs := []glyph{{asset: "zone-purge"}}
		return append(gs, effectGlyphs(v.Then)...), true
	case engine.Shuffle:
		return []glyph{{asset: "zone-deck", sigil: sigilShuffle}}, true
	case engine.ShuffleIntoDeck:
		// A multi-zone shuffle has no single source glyph, so it shows what it takes.
		// The shuffle sigil on the deck is what separates it from a plain PutCard,
		// which lands a card on top at a known position (ADR 0047).
		src := glyph{asset: "zone-discard"}
		if len(v.From) > 1 {
			src = glyph{asset: "type-creature"}
		}
		return []glyph{
			quantityOn(src, v.Quantity),
			arrowTo(glyph{asset: "zone-deck", sigil: sigilShuffle}),
		}, true
	case engine.SwapDeckAndDiscard:
		return []glyph{
			{asset: "zone-deck"},
			{asset: "glyph-swap"},
			{asset: "zone-discard"},
		}, true
	case engine.RefillHand:
		h := glyph{asset: "zone-hand"}
		if v.Player == engine.EachPlayer {
			h.decor = decorEach
		}
		return []glyph{{asset: "zone-deck"}, arrowTo(h)}, true
	case engine.PutItIntoHand:
		return []glyph{{asset: "glyph-return"}}, true
	case engine.PutFromPlay:
		if a := destinationGlyph(v.Destination); a != "" {
			return append(targetGlyphs(v.Target), arrowTo(glyph{asset: a})), true
		}
		return unmapped(v.Destination), true
	case engine.PutChosen:
		if a := destinationGlyph(v.Destination); a != "" {
			// PutChosen is the bounded-choice counterpart to PutFromPlay: the
			// controller picks how many of the pool move, so the subject is outlined
			// as chosen and carries the count.
			src := targetGlyphs(v.Target)
			src[0].decor |= decorChosen
			src[0] = quantityOn(src[0], v.Quantity)
			return append(src, arrowTo(glyph{asset: a})), true
		}
		return unmapped(v.Destination), true
	case engine.PutCard:
		if a := destinationGlyph(v.Destination); a != "" {
			return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: a})}, true
		}
		return unmapped(v.Destination), true
	case engine.PutFromHand:
		return []glyph{{asset: "zone-hand"}, arrowTo(glyph{asset: "glyph-play"})}, true
	case engine.PutDiscardedIntoHand:
		return []glyph{{asset: "zone-discard"}, arrowTo(glyph{asset: "zone-hand"})}, true
	case engine.PutNextTacticIntoHand:
		return []glyph{{asset: "type-tactic"}, arrowTo(glyph{asset: "zone-hand"})}, true
	default:
		return nil, false
	}
}

// quantityGlyphs renders a Quantity on its own as the cards it takes. The
// Quantity family is gated (TestEveryQuantityDraws, TestIconTotality), so every
// row here draws.
//
// On a verb the same quantity is a badge on the verb's own noun (quantityOn)
// rather than a card of its own: a quantity says how many, so it has nothing to
// draw until there is something to count. Standing alone — in the style gallery,
// and to the gate — it counts the abstract card.
func quantityGlyphs(q engine.Quantity) []glyph {
	switch q.(type) {
	case engine.Takes, engine.UpTo, engine.AnyNumber:
		return []glyph{quantityOn(glyph{asset: "card-back"}, q)}
	}
	return unmapped(q)
}

// quantityOn badges a verb's noun with how many cards it takes: the numeral when
// the number is a constant above one, the stack when there is no ceiling at all,
// and the may mark when the controller may stop short of the number.
//
// A board-scaled number prints no numeral, for the reason a scaled damage amount
// prints none — "for each friendly Shard" is a clause the strip cannot say, so it
// stays in the rules text beside it.
func quantityOn(g glyph, q engine.Quantity) glyph {
	g.qty = engine.FixedCardCount(q)
	switch q.(type) {
	case engine.UpTo:
		g.sigil = sigilMay
	case engine.AnyNumber:
		g.decor |= decorEach
		g.sigil = sigilMay
	}
	return g
}
