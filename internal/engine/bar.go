package engine

// A Bar is a turn-scoped restriction on a player: the value it imposes together
// with the card that imposed it, so a reminder can name the reason without a
// second list to keep in step with the bars themselves. The zero Bar imposes
// nothing, and clearing the restriction clears its source with it.
type Bar[T comparable] struct {
	Value  T
	Source LocalID
}

// A CreatureBar is a board-wide "creatures cannot fight/reap" restriction (Into
// the Night, Sow Salt): the barred Action (fighting or reaping) and the one house
// it spares (HouseNone spares none). It is comparable, so it rides in a Bar and
// in flat state; the zero value (an unset Action) bars nothing.
type CreatureBar struct {
	Action      UseKind
	ExceptHouse House
}

// A HouseWager is a bet on a player's active house on a future turn: if that
// player chooses House, the Predictor steals Amount (Snaglet). The zero value
// (Amount 0) arms no wager, and paying it out clears it.
type HouseWager struct {
	House  House
	Amount int
	// Predictor is the player who armed the wager and collects Amount if it hits.
	Predictor int
	Source    LocalID
}
