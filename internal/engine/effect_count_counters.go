package engine

// CountersOnThis counts the generic counters of one kind sitting on the source
// card — the running "+1 Æmber for each growth counter on Vineapple Tree" tally a
// key-cost change scales by. It reads its own card, so a leading clause names the
// source rather than saying "it".
type CountersOnThis struct {
	// Kind is which counter to count.
	Kind CounterKind
}

// Value returns how many counters of the kind sit on the source card.
func (c CountersOnThis) Value(ctx *EffectContext) int {
	return ctx.Resolver.CountersOn(ctx.Source, c.Kind)
}

// CountText renders the singular noun the "for each" clause repeats, e.g.
// "growth counter on {self}".
func (c CountersOnThis) CountText() string {
	return c.Kind.noun() + " on " + SelfName
}
