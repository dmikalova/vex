package engine

import (
	"fmt"
	"strings"
)

// DamageSourceMatcher selects the creatures a card's CannotBeDealtDamageBy passive
// refuses damage from. A creature matches if it carries Trait, or its power is at
// least MinPower — Ardent Hero refuses Mutant creatures or creatures with power 5
// or higher. The zero value matches nothing, so the passive is absent.
type DamageSourceMatcher struct {
	// Trait matches a source carrying this trait; traitUnset (the zero value)
	// matches on trait not at all.
	Trait Trait
	// MinPower matches a source whose power is at least this; 0 matches on power not
	// at all. Power is read live at damage time, so a buff into or out of the
	// threshold counts.
	MinPower int
}

// Active reports whether the matcher names any source, so a card carries the
// passive only when at least one criterion is set.
func (m DamageSourceMatcher) Active() bool {
	return m.Trait != traitUnset || m.MinPower > 0
}

// clause renders the sources the matcher names, e.g. "Mutant creatures or
// creatures with power 5 or higher".
func (m DamageSourceMatcher) clause() string {
	parts := make([]string, 0, 2)
	if m.Trait != traitUnset {
		parts = append(parts, m.Trait.String()+" creatures")
	}
	if m.MinPower > 0 {
		parts = append(parts, fmt.Sprintf("creatures with power %d or higher", m.MinPower))
	}
	return strings.Join(parts, " or ")
}

// refusesDamageFrom reports whether id's CannotBeDealtDamageBy passive blocks
// damage dealt to it by source. Damage with no credited source (source 0) is never
// blocked: only a named creature can match.
func (g *Game) refusesDamageFrom(id, source LocalID) bool {
	m := g.cat.def(id).CannotBeDealtDamageBy
	if !m.Active() || source == 0 {
		return false
	}
	if m.Trait != traitUnset && g.HasTrait(source, m.Trait) {
		return true
	}
	return m.MinPower > 0 && g.Power(source) >= m.MinPower
}
