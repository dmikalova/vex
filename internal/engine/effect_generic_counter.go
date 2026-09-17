package engine

import "fmt"

// PlaceCounter puts generic counters on each card its Target selects. A generic
// counter is a card-placed marker that does nothing on its own and matters only
// to a card that reads it — a doom counter (Wretched Doll destroys every creature
// carrying one) and its kin (ADR 0024).
type PlaceCounter struct {
	// Kind is which counter to place.
	Kind CounterKind
	// Target picks the cards the counters are placed on.
	Target Target
	// Amount is how many to place on each; zero means one.
	Amount int
	// Per multiplies Amount by a running count "for each" — Book of Malefaction
	// places a warrant counter for each Æmber stolen from its controller.
	Per Count
}

// amount returns how many counters to place, defaulting to one.
func (e PlaceCounter) amount() int {
	if e.Amount <= 0 {
		return 1
	}
	return e.Amount
}

// validate requires a kind and a target for the counters to land on.
func (e PlaceCounter) validate() error {
	if !e.Kind.valid() {
		return fmt.Errorf("PlaceCounter: counter kind not set")
	}
	if !e.Target.valid() {
		return errUnsetTarget("PlaceCounter")
	}
	return nil
}

// Text renders the effect, e.g. "put a doom counter on a creature" or "put 2 doom
// counters on a creature". A Per count leads the sentence, e.g. "for each Æmber
// stolen, put a warrant counter on {self}".
func (e PlaceCounter) Text() string {
	var body string
	if n := e.amount(); n > 1 {
		body = fmt.Sprintf("put %d %ss on %s", n, e.Kind.noun(), e.Target.Text())
	} else {
		body = fmt.Sprintf("put a %s on %s", e.Kind.noun(), e.Target.Text())
	}
	return forEach(e.Per, body)
}

// Resolve places the counters on each selected card, scaling by the Per count.
func (e PlaceCounter) Resolve(ctx *EffectContext) {
	n := scaled(e.amount(), e.Per, ctx)
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.PlaceCounter(id, e.Kind, n)
	}
}
