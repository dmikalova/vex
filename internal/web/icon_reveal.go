package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the deck-and-hand reveal
// family (ADR 0047) — looking at or revealing hidden cards, and putting a
// revealed card into whichever zone its search resolves to.

// revealEffectGlyphs claims every look/reveal/search effect.
func revealEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.LookAtTopOfDeck:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-look"}}, true, true
	case engine.RevealHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true, true
	case engine.RevealRandomFromHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true, true
	case engine.RevealChosenFromHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look"}}, true, true
	case engine.RevealTopOfDeck:
		g := []glyph{{asset: "zone-deck"}, {asset: "glyph-look"}}
		for _, act := range v.Then {
			if m, ok := act.(engine.ChooseAndMove); ok && m.Dest == engine.IntoPurge {
				g = append(g, arrowTo(glyph{asset: "zone-purge"}))
			}
		}
		return g, true, true
	case engine.PutRevealedCard:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: deckDestZone(v.To)})}, true, true
	case engine.Search:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-search"}}, true, true
	default:
		return nil, false, false
	}
}
