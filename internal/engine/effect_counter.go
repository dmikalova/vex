package engine

import "fmt"

// Power counters are permanent +1/-1 power tokens placed on a creature: they
// raise (or lower) its power for as long as it stays in play, and are shed when
// it leaves. A creature can hold any number of them. AddPowerCounter places
// counters on each creature its Target selects.
type AddPowerCounter struct {
	// Target picks the creatures the counters are placed on.
	Target Target
	// Amount is the total power the counters add, placed as that many +1 (or -1)
	// counters — Amount: 2 is "two +1 power counters", not one +2 counter.
	Amount int
	// Per scales the counters by a board quantity, choosing the target only once
	// (Martian Hounds counts the damaged creatures).
	Per Count
}

// validate requires an explicit target.
func (e AddPowerCounter) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("AddPowerCounter")
	}
	return nil
}

// Text renders the effect, e.g. "give Eater of the Dead a +1 power counter".
func (e AddPowerCounter) Text() string {
	return forEach(e.Per, fmt.Sprintf("give %s %s", e.Target.Text(), e.counters()))
}

// counters renders the tokens placed, e.g. "a +1 power counter" or "two +1 power
// counters": a larger Amount is that many single counters, and the count is
// spelled out to match KeyForge's printed wording.
func (e AddPowerCounter) counters() string {
	unit := 1
	if e.Amount < 0 {
		unit = -1
	}
	if n := e.Amount * unit; n != 1 {
		return fmt.Sprintf("%s %+d power counters", spellCounters(n), unit)
	}
	return fmt.Sprintf("a %+d power counter", unit)
}

// counterWords spells out the small counts KeyForge prints as words.
var counterWords = map[int]string{
	2: "two", 3: "three", 4: "four", 5: "five",
	6: "six", 7: "seven", 8: "eight", 9: "nine", 10: "ten",
}

// spellCounters renders a token count as its English word, falling back to the
// digits for counts larger than KeyForge ever prints.
func spellCounters(n int) string {
	if w, ok := counterWords[n]; ok {
		return w
	}
	return fmt.Sprintf("%d", n)
}

// Resolve places the counters on each selected creature, scaled by Per.
func (e AddPowerCounter) Resolve(ctx *EffectContext) {
	amount := scaled(e.Amount, e.Per, ctx)
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddPowerCounter(id, amount)
	}
}
