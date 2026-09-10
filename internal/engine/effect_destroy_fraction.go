package engine

import "fmt"

// A Fraction is a portion of a count, rounded up — one third, one half. It is the
// creature-count analog of the Æmber-pool Half Loss: a reusable share so an effect
// that acts on "one third" or "one half" of a battleline needs no card-specific
// node. The zero value is invalid; use OneThird or OneHalf.
type Fraction struct {
	denominator int
	ordinal     string
}

var (
	// OneThird is a third of a count, rounded up.
	OneThird = Fraction{denominator: 3, ordinal: "third"}
	// OneHalf is half of a count, rounded up.
	OneHalf = Fraction{denominator: 2, ordinal: "half"}
)

// of returns the fraction of n, rounded up.
func (f Fraction) of(n int) int { return (n + f.denominator - 1) / f.denominator }

// word renders the fraction as it reads in "one third of ...".
func (f Fraction) word() string { return f.ordinal }

// DestroyFractionOfEachBattleline destroys a fraction of all enemy creatures and
// the same fraction of all friendly creatures, rounding each count up. It reads
// both counts from the pre-destruction board, lets the controller choose which
// creatures on each side are destroyed, then destroys the whole chosen set
// together so their Destroyed abilities see one another still in play. Tertiate is
// Portion OneThird.
type DestroyFractionOfEachBattleline struct {
	Portion Fraction
}

// validate requires a portion.
func (e DestroyFractionOfEachBattleline) validate() error {
	if e.Portion.denominator == 0 {
		return fmt.Errorf("DestroyFractionOfEachBattleline: Portion must be set")
	}
	return nil
}

// Text renders the effect.
func (e DestroyFractionOfEachBattleline) Text() string {
	return fmt.Sprintf(
		"destroy one %s of all enemy creatures and one %s of all friendly "+
			"creatures (rounding up each time)",
		e.Portion.word(), e.Portion.word(),
	)
}

// Resolve destroys the chosen fraction of each side, chosen by the controller, all
// at once.
func (e DestroyFractionOfEachBattleline) Resolve(ctx *EffectContext) {
	doomed := e.choose(ctx, ctx.Opponent())
	doomed = append(doomed, e.choose(ctx, ctx.Controller)...)
	Destroy{}.destroy(ctx, doomed)
}

// choose asks the controller to pick the Portion fraction of player p's creatures,
// one at a time, returning the chosen ids.
func (e DestroyFractionOfEachBattleline) choose(ctx *EffectContext, p int) []LocalID {
	line := ctx.Resolver.Battleline(p)
	want := e.Portion.of(len(line))
	picked := map[LocalID]bool{}
	chosen := make([]LocalID, 0, want)
	for len(chosen) < want {
		cands := make([]LocalID, 0, len(line))
		for _, id := range line {
			if !picked[id] {
				cands = append(cands, id)
			}
		}
		pick, ok := ctx.ChooseCard("Choose a creature to destroy", cands)
		if !ok {
			break
		}
		picked[pick] = true
		chosen = append(chosen, pick)
	}
	return chosen
}
