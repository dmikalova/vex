package engine

import "fmt"

// A PowerBound is the power filter a target narrows by: the comparison it makes
// and the number it measures against. The zero value bounds nothing, so a target
// that names no power admits every creature.
//
// Threshold and parity share one axis because they answer the same question —
// which powers does this card reach. No card asks two power questions at once,
// and one field cannot contradict itself the way a bound beside a parity flag
// could ("with power 3 or lower with odd power").
type PowerBound struct {
	// Kind is the comparison the bound makes; the zero value makes none.
	Kind PowerBoundKind
	// Amount is the power the comparison measures against. A parity kind reads
	// only whether the power is odd or even, so it ignores this.
	Amount int
}

// PowerBoundKind names the comparison a PowerBound makes against a power.
type PowerBoundKind int

const (
	// boundUnset is the zero value: the bound compares nothing and admits every
	// power.
	boundUnset PowerBoundKind = iota
	// BoundAtMost admits a power at or below Amount ("with power 3 or lower").
	BoundAtMost
	// BoundAtLeast admits a power at or above Amount ("with power 3 or higher").
	BoundAtLeast
	// BoundExactly admits exactly Amount power ("with power 1").
	BoundExactly
	// BoundOdd admits an odd power, whatever Amount says (Onyx Knight).
	BoundOdd
	// BoundEven admits an even power, whatever Amount says (Opal Knight).
	BoundEven
	// BoundLessThanSource admits a power below the source card's own, read live when
	// the effect resolves (Dreadbone Decimus destroys a creature with lower power
	// than itself). It takes no Amount, which is what keeps it a bound rather than a
	// Refinement: a threshold that is a literal number, or no number at all, is a
	// comparable filter; one whose threshold is a live Count stays a Refinement
	// (PowerLessThan).
	BoundLessThanSource
)

// allPowerBoundKinds returns every real bound kind, so the census enumerates them
// rather than keeping a list a new kind would fall out of.
func allPowerBoundKinds() []PowerBoundKind {
	return []PowerBoundKind{
		BoundAtMost, BoundAtLeast, BoundExactly, BoundOdd, BoundEven, BoundLessThanSource,
	}
}

// filters reports whether the bound narrows its candidates at all.
func (b PowerBound) filters() bool { return b.Kind != boundUnset }

// admits reports whether a creature of this power passes the bound. The zero
// value admits every power, pinned by TestZeroPowerBoundAdmitsEveryPower. The
// context is read only by the bounds measured against the board rather than
// against a printed number.
func (b PowerBound) admits(ctx *EffectContext, power int) bool {
	switch b.Kind {
	case BoundLessThanSource:
		return power < ctx.Resolver.Power(ctx.Source)
	case BoundAtMost:
		return power <= b.Amount
	case BoundAtLeast:
		return power >= b.Amount
	case BoundExactly:
		return power == b.Amount
	case BoundOdd:
		return power%2 != 0
	case BoundEven:
		return power%2 == 0
	}
	return true
}

// clause appends the bound's own clause to a quantified phrase, e.g. "each
// creature" becomes "each creature with power 3 or lower". A bound that compares
// nothing appends nothing.
func (b PowerBound) clause(phrase string) string {
	switch b.Kind {
	case BoundAtMost:
		return phrase + fmt.Sprintf(" with power %d or lower", b.Amount)
	case BoundAtLeast:
		return phrase + fmt.Sprintf(" with power %d or higher", b.Amount)
	case BoundExactly:
		return phrase + fmt.Sprintf(" with power %d", b.Amount)
	case BoundOdd:
		return phrase + " with odd power"
	case BoundEven:
		return phrase + " with even power"
	case BoundLessThanSource:
		return phrase + " with lower power than " + SelfName
	}
	return phrase
}
