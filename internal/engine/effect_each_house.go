package engine

import "fmt"

// ForEachHouse resolves Do once per house, in canonical order (Brobnar through
// Untamed), binding that house in context so a nested "of that house" target
// (a HouseMatcher of kind MatchEachHouse) reads it. It backs "for each house,
// ..." cards: Gleeful Mayhem deals damage to a creature of each house, and a
// future card chooses a creature of each house before destroying the rest.
type ForEachHouse struct {
	Do Effect
}

// validate rejects a missing Do and descends into it.
func (e ForEachHouse) validate() error {
	if e.Do == nil {
		return fmt.Errorf("ForEachHouse: Do must be set")
	}
	return validateEffect(e.Do)
}

// Text renders the effect, e.g. "for each house, deal 5 damage to a creature of
// that house".
func (e ForEachHouse) Text() string {
	return "for each house, " + e.Do.Text()
}

// Resolve binds each house in turn and resolves Do against it, restoring the
// prior house binding afterward so a ForEachHouse nested in another house-scoped
// effect leaves it untouched.
func (e ForEachHouse) Resolve(ctx *EffectContext) {
	prev := ctx.ChosenHouse
	for h := Brobnar; h <= Untamed; h++ {
		ctx.ChosenHouse = h
		e.Do.Resolve(ctx)
	}
	ctx.ChosenHouse = prev
}
