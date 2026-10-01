package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the deck-and-hand reveal
// family (ADR 0047) — looking at or revealing hidden cards, and putting a
// revealed card into whichever zone its search resolves to.
//
// Looking and revealing share the eye noun, so the reveal sigil is what says the
// cards are shown to both players rather than read privately: without it
// LookAtTopOfDeck and RevealTopOfDeck draw the same strip (ADR 0047).
//
// The TopAct steps the two reads route their cards through are drawn here too,
// beside the verb that holds them, the way a Count is drawn where its holder is.
// A step is an ordinary zone move wearing a step's clothes, so each one draws the
// zone its cards end in; what tells the steps apart is where the cards go and
// which cards go there, never a drawing of the step itself.

// revealEffectGlyphs claims every look/reveal/search effect.
func revealEffectGlyphs(e engine.Effect) ([]glyph, bool, bool) {
	switch v := e.(type) {
	case engine.LookAtTopOfDeck:
		gs, covered := topActsGlyphs(v.Then)
		return append([]glyph{{asset: "zone-deck"}, {asset: "glyph-look"}}, gs...), covered, true
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
		gs, covered := topActsGlyphs(v.Then)
		return append(
			[]glyph{{asset: "zone-deck"}, {asset: "glyph-look", sigil: sigilReveal}},
			gs...,
		), covered, true
	case engine.PutRevealedCard:
		return []glyph{{asset: "card-back"}, arrowTo(glyph{asset: deckDestZone(v.To)})}, true, true
	case engine.Search:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-search"}}, true, true
	default:
		return nil, false, false
	}
}

// topActsGlyphs renders the ordered routing steps a deck read sends its cards
// through, one after another, reporting covered only when every step maps.
func topActsGlyphs(acts []engine.TopAct) ([]glyph, bool) {
	out := make([]glyph, 0, len(acts)*2)
	covered := true
	for _, act := range acts {
		gs, c := topActGlyphs(act)
		out = append(out, gs...)
		covered = covered && c
	}
	return out, covered
}

// topActGlyphs renders one routing step as the cards it takes and the zone it
// sends them to, reporting whether a real mapping covered it. The TopAct family
// is gated (TestEveryTopActDraws, TestIconTotality), so every row here draws.
func topActGlyphs(act engine.TopAct) ([]glyph, bool) {
	switch v := act.(type) {
	case engine.ChooseAndMove:
		if a := deckDestZone(v.Dest); a != "" {
			// A single card prints no numeral, the way every other card-moving verb's
			// count does (engine.FixedCardCount).
			taken := glyph{asset: "card-back", decor: decorChosen}
			if v.Cards > 1 {
				taken.qty = v.Cards
			}
			return []glyph{taken, arrowTo(glyph{asset: a})}, true
		}
		return fallbackGlyphs(), false
	case engine.PartitionByChosenHouse:
		// The split routes every read card two ways, which is the same shape a
		// Conditional's two branches take: the cards of the chosen house open the
		// clause, and the branch mark opens where the others go.
		matching := deckDestZone(v.Matching)
		rest := deckDestZone(v.Rest)
		if matching == "" || rest == "" {
			return fallbackGlyphs(), false
		}
		return []glyph{
			{asset: "house-any", decor: decorChosen},
			arrowTo(glyph{asset: matching}),
			{asset: rest, punct: punctBranch},
		}, true
	case engine.MayDiscardLookedAt:
		// The one card the look read, which the controller may discard, so the may
		// mark rides the card rather than the pile it would go to.
		return []glyph{
			{asset: "card-back", sigil: sigilMay},
			arrowTo(glyph{asset: "zone-discard"}),
		}, true
	case engine.ReorderRest:
		// Back on top at an order the controller picks — a known position, which is
		// what tells it from the shuffle that puts them back at no known position.
		return []glyph{
			{asset: "card-back", decor: decorEach | decorChosen},
			arrowTo(glyph{asset: "zone-deck-top"}),
		}, true
	case engine.Shuffle:
		// The same deck wearing the same shuffle verb the Shuffle effect draws: it is
		// the one node that is both a step and an effect.
		return []glyph{{asset: "zone-deck", sigil: sigilShuffle}}, true
	}
	return fallbackGlyphs(), false
}
