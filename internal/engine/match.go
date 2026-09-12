package engine

import "strings"

// Match is a predicate over a card, selecting it by type, trait, and/or name. It
// is the shared vocabulary for "which cards does this effect act on" when the
// choice is filtered by what a card is rather than where it sits — a card
// recurred from the discard pile, for instance.
//
// An empty Match admits any card. A non-empty Match admits a card that satisfies
// every predicate it sets (Type, Trait, Name — conjoined), or, failing that, any
// alternative listed in Or. This lets one field express both a conjunction
// ("Horseman creature" is Type Creature and Trait Horseman) and a disjunction
// ("upgrade or Robot card" is Type Upgrade with Or a Trait Robot clause).
type Match struct {
	// Type requires the card to be of this type; the zero value requires no type.
	Type CardType
	// Trait requires the card to carry this trait; the zero value requires none.
	Trait Trait
	// Name requires the card to have this exact name; the zero value requires none.
	Name string
	// Or lists alternative matches: a card also qualifies if it satisfies any of
	// them. Chief Engineer Walls admits an upgrade or a Robot card.
	Or []Match
}

// empty reports that the match sets no predicate, so it admits every card.
func (m Match) empty() bool {
	return m.Type == TypeUnset && m.Trait == traitUnset && m.Name == "" &&
		len(m.Or) == 0
}

// admits reports whether the card satisfies the match.
func (m Match) admits(r StateReader, id LocalID) bool {
	if m.empty() {
		return true
	}
	if m.satisfiesClause(r, id) {
		return true
	}
	for _, alt := range m.Or {
		if alt.admits(r, id) {
			return true
		}
	}
	return false
}

// satisfiesClause reports whether the card satisfies this match's own conjoined
// predicates, ignoring Or. A clause that sets no predicate never matches on its
// own — it exists only to carry alternatives in Or.
func (m Match) satisfiesClause(r StateReader, id LocalID) bool {
	if m.Type == TypeUnset && m.Trait == traitUnset && m.Name == "" {
		return false
	}
	if m.Type != TypeUnset && r.TypeOf(id) != m.Type {
		return false
	}
	if m.Trait != traitUnset && !r.HasTrait(id, m.Trait) {
		return false
	}
	if m.Name != "" && r.Name(id) != m.Name {
		return false
	}
	return true
}

// noun renders the kind of card the match selects — the card's own name when Name
// is set (e.g. "Ortannu's Binding"), otherwise the trait-and-type phrase (e.g.
// "creature" or "Horseman creature"), with each Or alternative joined by " or "
// ("upgrade or Robot card"). An empty match renders the generic "card".
func (m Match) noun() string {
	if m.Name != "" {
		return m.Name
	}
	base := "card"
	if m.Type != TypeUnset {
		base = strings.ToLower(m.Type.String())
	}
	if m.Trait != traitUnset {
		base = m.Trait.String() + " " + base
	}
	for _, alt := range m.Or {
		base += " or " + alt.noun()
	}
	return base
}
