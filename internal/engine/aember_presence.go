package engine

// AemberPresence is the Æmber filter a target narrows by: whether a card has
// Æmber on it. The zero value asks nothing, and the two real values are
// exclusive — a card either carries Æmber or does not, so a card can never ask
// for both.
type AemberPresence int

const (
	// aemberAny is the zero value: the filter asks nothing about Æmber.
	aemberAny AemberPresence = iota
	// AemberSome admits a card with Æmber on it ("a creature with Æmber on it").
	AemberSome
	// AemberNone admits a card with no Æmber on it (Draining Touch destroys a
	// creature with no Æmber on it).
	AemberNone
)

// allAemberPresences returns every real Æmber presence, so the census enumerates
// them rather than keeping a list a new value would fall out of.
func allAemberPresences() []AemberPresence {
	return []AemberPresence{AemberSome, AemberNone}
}

// filters reports whether the presence narrows its candidates at all.
func (a AemberPresence) filters() bool { return a != aemberAny }

// admits reports whether a card carrying this much Æmber passes the filter.
func (a AemberPresence) admits(aember int) bool {
	switch a {
	case AemberSome:
		return aember != 0
	case AemberNone:
		return aember == 0
	}
	return true
}

// clause appends the presence's own clause to a quantified phrase, e.g. "each
// creature" becomes "each creature with Æmber on it". A presence that asks
// nothing appends nothing.
func (a AemberPresence) clause(phrase string) string {
	switch a {
	case AemberSome:
		return phrase + " with Æmber on it"
	case AemberNone:
		return phrase + " with no Æmber on it"
	}
	return phrase
}
