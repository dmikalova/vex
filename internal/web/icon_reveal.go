package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the deck-and-hand reveal
// family (ADR 0047) — looking at or revealing hidden cards, and putting a
// revealed card into whichever zone its search resolves to.
//
// Looking and revealing share the eye noun, so the reveal sigil is what says the
// cards are shown to both players rather than read privately: without it
// LookAtTopOfDeck and RevealTopOfDeck draw the same strip (ADR 0047).

// revealEffectGlyphs claims every look/reveal/search effect.
func revealEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.LookAtTopOfDeck:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-look"}}, true, true
	case engine.RevealHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look", sigil: sigilReveal}}, true, true
	case engine.RevealRandomFromHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look", sigil: sigilReveal}}, true, true
	case engine.RevealChosenFromHand:
		return []glyph{
			{asset: "zone-hand", decor: decorChosen},
			{asset: "glyph-look", sigil: sigilReveal},
		}, true, true
	case engine.RevealTopOfDeck:
		g := []glyph{{asset: "zone-deck"}, {asset: "glyph-look", sigil: sigilReveal}}
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
