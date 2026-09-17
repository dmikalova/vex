package engine

import "fmt"

// RemoveCounters takes every generic counter of one kind off each card its Target
// selects, leaving the card's other kinds untouched — Vineapple Tree removes each
// growth counter from itself after a key is forged.
type RemoveCounters struct {
	// Kind is which counter to remove.
	Kind CounterKind
	// Target picks the cards the counters are removed from.
	Target Target
	// Amount removes just that many counters of the kind from each selected card;
	// zero removes every counter of the kind (The Colosseum removes six).
	Amount int
}

// validate requires a kind and a target for the counters to leave.
func (e RemoveCounters) validate() error {
	if !e.Kind.valid() {
		return fmt.Errorf("RemoveCounters: counter kind not set")
	}
	if !e.Target.valid() {
		return errUnsetTarget("RemoveCounters")
	}
	if e.Amount < 0 {
		return fmt.Errorf("RemoveCounters: Amount cannot be negative")
	}
	return nil
}

// Text renders the effect, e.g. "remove each growth counter from {self}", "remove
// a warrant counter from {self}", or "remove 6 glory counters from {self}".
func (e RemoveCounters) Text() string {
	switch {
	case e.Amount == 1:
		return fmt.Sprintf("remove a %s from %s", e.Kind.noun(), e.Target.Text())
	case e.Amount > 1:
		return fmt.Sprintf(
			"remove %d %ss from %s", e.Amount, e.Kind.noun(), e.Target.Text())
	default:
		return fmt.Sprintf("remove each %s from %s", e.Kind.noun(), e.Target.Text())
	}
}

// Resolve drops the kind's counters from each selected card.
func (e RemoveCounters) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate drops the kind's counters from each selected card and reports
// whether it took any off, so it can be the first half of a Then — Book of
// Malefaction removes a warrant counter and, only if it had one, purges a creature.
func (e RemoveCounters) resolveGate(ctx *EffectContext) bool {
	removed := false
	for _, id := range e.Target.Select(ctx) {
		if ctx.Resolver.CountersOn(id, e.Kind) == 0 {
			continue
		}
		removed = true
		if e.Amount > 0 {
			ctx.Resolver.RemoveCountersN(id, e.Kind, e.Amount)
			continue
		}
		ctx.Resolver.RemoveCounters(id, e.Kind)
	}
	return removed
}
