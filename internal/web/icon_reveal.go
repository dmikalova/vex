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
func revealEffectGlyphs(e engine.Effect) ([]glyph, bool) {
	switch v := e.(type) {
	case engine.LookAtTopOfDeck:
		return append(
			[]glyph{{asset: "zone-deck"}, {asset: "glyph-look"}},
			topActsGlyphs(v.Then)...,
		), true
	case engine.RevealHand:
		// Every card in that hand, which is the stack a count already draws on a pile
		// (CardsInHand); RevealRandomFromHand below reveals one card out of the same
		// hand, so the stack is what tells the two apart.
		return []glyph{
			{asset: "zone-hand", decor: decorEach | playerDecor(v.Player)},
			{asset: "glyph-look", sigil: sigilReveal},
		}, true
	case engine.RevealRandomFromHand:
		return []glyph{{asset: "zone-hand"}, {asset: "glyph-look", sigil: sigilReveal}}, true
	case engine.RevealChosenFromHand:
		return []glyph{
			{asset: "zone-hand", decor: decorChosen},
			{asset: "glyph-look", sigil: sigilReveal},
		}, true
	case engine.RevealTopOfDeck:
		return append(
			[]glyph{{asset: "zone-deck"}, {asset: "glyph-look", sigil: sigilReveal}},
			topActsGlyphs(v.Then)...,
		), true
	case engine.PutRevealedCard:
		// The card the read put in context, drawn as every other card-in-context is,
		// which separates it from the face-down card under a card (ArchiveCardUnder).
		return []glyph{
			{asset: "card-back", decor: decorChosen},
			arrowTo(glyph{asset: deckDestZone(v.To)}),
		}, true
	case engine.Search:
		return []glyph{{asset: "zone-deck"}, {asset: "glyph-search"}}, true
	default:
		return nil, false
	}
}

// topActsGlyphs renders the ordered routing steps a deck read sends its cards
// through, one after another.
func topActsGlyphs(acts []engine.TopAct) []glyph {
	out := make([]glyph, 0, len(acts)*2)
	for _, act := range acts {
		out = append(out, topActGlyphs(act)...)
	}
	return out
}

// topActGlyphs renders one routing step as the cards it takes and the zone it
// sends them to. The TopAct family is gated (TestEveryTopActDraws,
// TestIconTotality), so every row here draws.
func topActGlyphs(act engine.TopAct) []glyph {
	switch v := act.(type) {
	case engine.ChooseAndMove:
		if a := deckDestZone(v.Dest); a != "" {
			// A single card prints no numeral, the way every other card-moving verb's
			// count does (engine.FixedCardCount).
			taken := glyph{asset: "card-back", decor: decorChosen}
			if v.Cards > 1 {
				taken.qty = v.Cards
			}
			return []glyph{taken, arrowTo(glyph{asset: a})}
		}
		return unmapped(v)
	case engine.PartitionByChosenHouse:
		// The split routes every read card two ways, which is the same shape a
		// Conditional's two branches take: the cards of the chosen house open the
		// clause, and the branch mark opens where the others go.
		matching := deckDestZone(v.Matching)
		rest := deckDestZone(v.Rest)
		if matching == "" || rest == "" {
			return unmapped(v)
		}
		return []glyph{
			{asset: "house-any", decor: decorChosen},
			arrowTo(glyph{asset: matching}),
			{asset: rest, punct: punctBranch},
		}
	case engine.MayDiscardLookedAt:
		// The one card the look read, which the controller may discard, so the may
		// mark rides the card rather than the pile it would go to.
		return []glyph{
			{asset: "card-back", sigil: sigilMay},
			arrowTo(glyph{asset: "zone-discard"}),
		}
	case engine.ReorderRest:
		// Back on top at an order the controller picks — a known position, which is
		// what tells it from the shuffle that puts them back at no known position.
		return []glyph{
			{asset: "card-back", decor: decorEach | decorChosen},
			arrowTo(glyph{asset: "zone-deck-top"}),
		}
	case engine.Shuffle:
		// The same deck wearing the same shuffle verb the Shuffle effect draws: it is
		// the one node that is both a step and an effect.
		return []glyph{{asset: "zone-deck", sigil: sigilShuffle}}
	}
	return unmapped(act)
}
