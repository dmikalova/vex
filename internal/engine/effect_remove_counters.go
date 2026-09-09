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

// Text renders the effect, e.g. "remove each growth counter from {self}" or
// "remove 6 glory counters from {self}".
func (e RemoveCounters) Text() string {
	if e.Amount > 0 {
		return fmt.Sprintf(
			"remove %d %ss from %s", e.Amount, e.Kind.noun(), e.Target.Text())
	}
	return fmt.Sprintf("remove each %s from %s", e.Kind.noun(), e.Target.Text())
}

// Resolve drops the kind's counters from each selected card.
func (e RemoveCounters) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if e.Amount > 0 {
			ctx.Resolver.RemoveCountersN(id, e.Kind, e.Amount)
			continue
		}
		ctx.Resolver.RemoveCounters(id, e.Kind)
	}
}
