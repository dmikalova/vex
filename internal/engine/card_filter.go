package engine

// CardFilter is a predicate over a card's identity, admitting it by type, trait,
// and/or name. It is the identity axis a pile Selection (Chosen, Each) narrows by
// — "which cards does this effect act on" filtered by what a card is rather than
// where it sits or who owns it — and the element type of a Selection's Or
// disjunction.
//
// An empty CardFilter admits any card. A non-empty one admits a card that
// satisfies every predicate it sets (Type, Trait, Name — conjoined), or, failing
// that, any alternative listed in Or. This lets one value express both a
// conjunction ("Horseman creature" is Type Creature and Trait Horseman) and a
// disjunction ("upgrade or Robot card" is Type Upgrade with Or a Trait Robot
// clause).
type CardFilter struct {
	// Type requires the card to be of this type; the zero value requires no type.
	Type CardType
	// Trait requires the card to carry this trait; the zero value requires none.
	Trait Trait
	// Name requires the card to have this exact name; the zero value requires none.
	Name string
	// Gigantic requires the card to be a half of a gigantic creature (either half);
	// the zero value requires none. The tutors search for a gigantic's halves.
	Gigantic bool
	// Or lists alternative filters: a card also qualifies if it satisfies any of
	// them. Chief Engineer Walls admits an upgrade or a Robot card.
	Or []CardFilter
}

// empty reports that the filter sets no predicate, so it admits every card.
func (f CardFilter) empty() bool {
	return f.Type == TypeUnset && f.Trait == traitUnset && f.Name == "" &&
		!f.Gigantic && len(f.Or) == 0
}

// admits reports whether the card satisfies the filter.
func (f CardFilter) admits(r StateReader, id LocalID) bool {
	if f.empty() {
		return true
	}
	if f.satisfiesClause(r, id) {
		return true
	}
	for _, alt := range f.Or {
		if alt.admits(r, id) {
			return true
		}
	}
	return false
}

// satisfiesClause reports whether the card satisfies this filter's own conjoined
// predicates, ignoring Or. A clause that sets no predicate never matches on its
// own — it exists only to carry alternatives in Or.
func (f CardFilter) satisfiesClause(r StateReader, id LocalID) bool {
	if f.Type == TypeUnset && f.Trait == traitUnset && f.Name == "" && !f.Gigantic {
		return false
	}
	if f.Type != TypeUnset && r.TypeOf(id) != f.Type {
		return false
	}
	if f.Trait != traitUnset && !r.HasTrait(id, f.Trait) {
		return false
	}
	if f.Name != "" && r.Name(id) != f.Name {
		return false
	}
	if f.Gigantic && r.GiganticRoleOf(id) == GiganticNone {
		return false
	}
	return true
}

// noun renders the kind of card the filter admits — the card's own name when Name
// is set (e.g. "Ortannu's Binding"), otherwise the trait-and-type phrase (e.g.
// "creature" or "Horseman creature"), with each Or alternative joined by " or "
// ("upgrade or Robot card"). An empty filter renders the generic "card".
func (f CardFilter) noun() string {
	if f.Name != "" {
		return f.Name
	}
	base := "card"
	if f.Type != TypeUnset {
		base = typeWord(f.Type)
	}
	if f.Gigantic {
		base = "gigantic creature"
	}
	if f.Trait != traitUnset {
		base = f.Trait.String() + " " + base
	}
	for _, alt := range f.Or {
		base += " or " + alt.noun()
	}
	return base
}
