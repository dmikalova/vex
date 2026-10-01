package web

import "github.com/dmikalova/vex/internal/engine"

// This file is part of the Iconography pass (ADR 0022): the Selection family
// (ADR 0031) — how a zone verb picks the cards it acts on.
//
// A selection is not an effect, so it has no place in the glyphFamilies chain:
// it is a sentence part the verb holding it draws, the way a Count is. The
// grammar is one rule — a selection draws the card it picks — with the way it was
// picked carried by the decorations the strip already has:
//
//   - Chosen is the chosen outline, Each the stack, Self the self marker. The
//     three are one noun under three decorations rather than three drawings, which
//     is the economy a Target already makes (targetGlyph).
//   - Named is the nameplate, because the only thing identifying the card is its
//     printed name and the strip speaks no words.
//   - Random, Top and Bottom name a position in a pile rather than a card: a
//     random card is at no known position, so it wears the shuffle verb, and the
//     two positional picks draw the end of the pile they take from.
//
// The identity filters a Chosen or an Each narrows by — house, type, trait, name —
// are drawn, because they decide which cards the verb can reach at all. Another
// only excludes the card a previous pick already took, so like every other fine
// filter it stays in the rules text beside the strip.

// selectionGlyphs renders a Selection as the card it picks, reporting whether a
// real mapping covered it. The Selection family is gated (TestEverySelectionDraws,
// TestIconTotality), so every row here draws.
func selectionGlyphs(s engine.Selection) ([]glyph, bool) {
	switch v := s.(type) {
	case engine.Chosen:
		return pickedCardGlyphs(v.House, cardPick{
			cardType: v.Type,
			trait:    v.Trait,
			name:     v.Name,
			how:      decorChosen,
			verb:     optionalVerb(v.Optional),
		}), true
	case engine.Each:
		return pickedCardGlyphs(v.House, cardPick{
			cardType: v.Type,
			trait:    v.Trait,
			name:     v.Name,
			how:      decorEach,
		}), true
	case engine.Random:
		// A card at no known position in its pile, which is what the shuffle verb
		// already says of a card put back into a deck.
		return []glyph{{asset: "card-back", sigil: sigilShuffle}}, true
	case engine.Named:
		return []glyph{nameGlyph()}, true
	case engine.Self:
		return []glyph{{asset: "card-back", decor: decorThis}}, true
	case engine.Top:
		return []glyph{{asset: "zone-deck-top"}}, true
	case engine.Bottom:
		return []glyph{{asset: "zone-deck-bottom"}}, true
	}
	return fallbackGlyphs(), false
}

// cardPick is how a filtering selection reaches its card: the type it names, the
// trait and name it narrows by, the decoration saying how the card was picked,
// and the verb it wears. It is a record rather than five parameters because the
// two filtering selections differ only in the fields they fill.
type cardPick struct {
	cardType engine.CardType
	trait    engine.Trait
	name     string
	how      decor
	verb     sigil
}

// pickedCardGlyphs draws the card a filtering selection reaches: the house it must
// belong to, then the card itself under the decoration and verb that say how it
// was picked, then the trait or name the pick narrows by. It is the Selection
// counterpart of cardsInPlayCountGlyphs and narrows in the same order, so a pick
// and a count of the same cards read alike.
func pickedCardGlyphs(house engine.HouseMatcher, p cardPick) []glyph {
	gs := append(houseFilterGlyphs(house), glyph{
		asset: countedCardAsset(p.cardType),
		decor: p.how,
		sigil: p.verb,
	})
	if p.trait != 0 {
		gs = append(gs, traitGlyph())
	}
	if p.name != "" {
		gs = append(gs, nameGlyph())
	}
	return gs
}

// optionalVerb is the verb a declinable pick wears: the may mark when the
// controller can pass, and no mark at all when the pick is forced. The verb rides
// the card rather than the head of the run, which a house filter may occupy — "you
// may discard a Dis card" is the card that is optional, not the house.
func optionalVerb(optional bool) sigil {
	if optional {
		return sigilMay
	}
	return sigilNone
}
