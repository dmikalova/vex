package engine

// DamagePresence is the damage filter a target narrows by: whether a creature
// has damage on it. The zero value asks nothing, and the two real values are
// exclusive — a creature either carries damage or does not, so a card can never
// ask for both.
type DamagePresence int

const (
	// damageAny is the zero value: the filter asks nothing about damage.
	damageAny DamagePresence = iota
	// DamageSome admits a creature that has damage on it ("a damaged creature").
	DamageSome
	// DamageNone admits a creature that has no damage on it ("an undamaged
	// creature").
	DamageNone
)

// allDamagePresences returns every real damage presence, so the census
// enumerates them rather than keeping a list a new value would fall out of.
func allDamagePresences() []DamagePresence {
	return []DamagePresence{DamageSome, DamageNone}
}

// filters reports whether the presence narrows its candidates at all.
func (d DamagePresence) filters() bool { return d != damageAny }

// admits reports whether a creature carrying this much damage passes the filter.
func (d DamagePresence) admits(damage int) bool {
	switch d {
	case DamageSome:
		return damage != 0
	case DamageNone:
		return damage == 0
	}
	return true
}

// adjective is the word the presence prefixes to the noun it narrows, e.g.
// "damaged" in "each damaged creature". It is empty when the filter asks
// nothing.
func (d DamagePresence) adjective() string {
	switch d {
	case DamageSome:
		return "damaged"
	case DamageNone:
		return "undamaged"
	}
	return ""
}
