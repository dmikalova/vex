package engine

// Rounding is how a Fraction rounds a share that does not divide evenly. It is
// always stated explicitly on the Fraction — no effect assumes a direction.
type Rounding uint8

const (
	roundingUnset Rounding = iota
	// RoundDown rounds a fractional share down, toward zero.
	RoundDown
	// RoundUp rounds a fractional share up.
	RoundUp
)

// rounded renders the direction as it reads after an amount, "rounded down" /
// "rounded up".
func (r Rounding) rounded() string {
	if r == RoundUp {
		return "rounded up"
	}
	return "rounded down"
}

// rounding renders the direction as a gerund, "rounding down" / "rounding up".
func (r Rounding) rounding() string {
	if r == RoundUp {
		return "rounding up"
	}
	return "rounding down"
}

// Fraction is a share of a quantity — a half or a third — with its rounding stated
// explicitly. One Fraction is the whole game's "portion of N": a share of an Æmber
// pool (LoseAember.By, since a Fraction is a Loss), of a creature's power
// (PowerOfChosen.Of), or of a battleline count (Portion, the both-sides
// refinement Tertiate carries). KeyForge only ever needs a half or a third. The
// zero value is invalid; use one of the four package vars.
type Fraction struct {
	denominator int
	ordinal     string
	round       Rounding
}

var (
	// HalfRoundedDown is half a quantity, rounded down (lose half your Æmber).
	HalfRoundedDown = Fraction{
		denominator: 2,
		ordinal:     "half",
		round:       RoundDown,
	}
	// HalfRoundedUp is half a quantity, rounded up.
	HalfRoundedUp = Fraction{
		denominator: 2,
		ordinal:     "half",
		round:       RoundUp,
	}
	// ThirdRoundedDown is a third of a quantity, rounded down.
	ThirdRoundedDown = Fraction{
		denominator: 3,
		ordinal:     "third",
		round:       RoundDown,
	}
	// ThirdRoundedUp is a third of a quantity, rounded up (Tertiate).
	ThirdRoundedUp = Fraction{
		denominator: 3,
		ordinal:     "third",
		round:       RoundUp,
	}
)

// valid reports that the fraction was built from one of the package vars.
func (f Fraction) valid() bool { return f.denominator != 0 && f.round != roundingUnset }

// of returns the fraction of n, rounded as the fraction states.
func (f Fraction) of(n int) int {
	if f.round == RoundUp {
		return (n + f.denominator - 1) / f.denominator
	}
	return n / f.denominator
}

// word renders the fraction as it reads in "one third of ...".
func (f Fraction) word() string { return f.ordinal }

// roundingPhrase renders the fraction's rounding as a gerund, "rounding up".
func (f Fraction) roundingPhrase() string { return f.round.rounding() }

// lose, object, and qualifier make a Fraction a Loss, so By: HalfRoundedDown
// removes half an Æmber pool.
func (f Fraction) lose(pool int) int { return f.of(pool) }

// object renders the amount as the object of "loses …", e.g. "half of their Æmber,
// rounded down".
func (f Fraction) object(possessive string) string {
	return f.ordinal + " of " + possessive + " Æmber, " + f.round.rounded()
}

// qualifier is empty: a fractional loss narrows no players.
func (f Fraction) qualifier() string { return "" }

// countPhrase makes a Fraction a portionPhraser, so PowerOfChosen{Of:
// HalfRoundedDown} reads "half its power, rounded down".
func (f Fraction) countPhrase(noun string) string {
	return f.ordinal + " " + noun + ", " + f.round.rounded()
}
