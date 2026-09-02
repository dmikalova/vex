package engine

// A Bar is a turn-scoped restriction on a player: the value it imposes together
// with the card that imposed it, so a reminder can name the reason without a
// second list to keep in step with the bars themselves. The zero Bar imposes
// nothing, and clearing the restriction clears its source with it.
type Bar[T comparable] struct {
	Value  T
	Source LocalID
}
