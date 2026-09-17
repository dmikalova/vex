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
	// Equal sets the number of +1 counters to a live count rather than a fixed
	// Amount — Mimic Gel gives itself counters equal to a chosen creature's power.
	// When set it overrides Amount and renders "equal to <count>".
	Equal Count
	// Walk places counters along an inward flank walk instead of on Target: Walk[i]
	// +1 power counters to the creature at step i from a chosen flank creature
	// (Growth Surge, Walk{3, 2, 1}). When set it overrides Target, Amount, Per, and
	// Equal.
	Walk []int
}

// validate requires an explicit target unless the counters are placed along a
// flank walk, which chooses its own creatures.
func (e AddPowerCounter) validate() error {
	if len(e.Walk) > 0 {
		return nil
	}
	if !e.Target.valid() {
		return errUnsetTarget("AddPowerCounter")
	}
	return nil
}

// Text renders the effect, e.g. "give Eater of the Dead a +1 power counter".
func (e AddPowerCounter) Text() string {
	if len(e.Walk) > 0 {
		return e.walkText()
	}
	if e.Equal != nil {
		return fmt.Sprintf("give %s +1 power counters equal to %s",
			e.Target.Text(), e.Equal.CountText())
	}
	return forEach(e.Per, fmt.Sprintf("give %s %s", e.Target.Text(), e.counters()))
}

// walkText renders the inward flank walk, e.g. "choose a flank creature. Give it
// three +1 power counters, its neighbor two +1 power counters, and the neighbor's
// other neighbor a +1 power counter."
func (e AddPowerCounter) walkText() string {
	parts := make([]string, len(e.Walk))
	for i, a := range e.Walk {
		parts[i] = flankWalkPhrase(i) + " " + counterTokens(a)
	}
	joined := parts[0]
	for i := 1; i < len(parts); i++ {
		sep := ", "
		if i == len(parts)-1 {
			sep = ", and "
		}
		joined += sep + parts[i]
	}
	return "choose a flank creature. Give " + joined
}

// counters renders the tokens placed, e.g. "a +1 power counter" or "two +1 power
// counters": a larger Amount is that many single counters, and the count is
// spelled out to match KeyForge's printed wording.
func (e AddPowerCounter) counters() string {
	return counterTokens(e.Amount)
}

// counterTokens renders n placed tokens, e.g. "a +1 power counter" or "three +1
// power counters": a larger count is that many single counters, spelled out to
// match KeyForge's printed wording.
func counterTokens(n int) string {
	unit := 1
	if n < 0 {
		unit = -1
	}
	if m := n * unit; m != 1 {
		return fmt.Sprintf("%s %+d power counters", spellCounters(m), unit)
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
	if len(e.Walk) > 0 {
		for _, st := range flankWalkSteps(ctx, e.Walk) {
			ctx.Resolver.AddPowerCounter(st.ID, st.Amount)
		}
		return
	}
	if e.Equal != nil {
		amount := e.Equal.Value(ctx)
		for _, id := range e.Target.Select(ctx) {
			ctx.Resolver.AddPowerCounter(id, amount)
		}
		return
	}
	amount := scaled(e.Amount, e.Per, ctx)
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.AddPowerCounter(id, amount)
	}
}
