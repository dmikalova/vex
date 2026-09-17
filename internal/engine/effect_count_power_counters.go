package engine

// PowerCountersOnThis counts the net +1/-1 power counters on the source card — the
// tally Chonkers doubles by adding that many again. It reads its own card, so its
// text names the source rather than saying "it".
type PowerCountersOnThis struct{}

// Value returns the net power counters on the source card.
func (PowerCountersOnThis) Value(ctx *EffectContext) int {
	return ctx.Resolver.PowerCountersOn(ctx.Source)
}

// CountText renders the "equal to ..." noun, e.g. "the number of +1 power counters
// on {self}".
func (PowerCountersOnThis) CountText() string {
	return "the number of +1 power counters on " + SelfName
}
