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
// the Night, Sow Salt): the barred Action (fighting or reaping) and the houses it
// reaches. Houses admits the barred creatures — an unset matcher (MatchAnyHouse)
// bars every house, a MatchExceptHouse spares one. It is comparable, so it rides
// in a Bar and in flat state; the zero value (an unset Action) bars nothing.
type CreatureBar struct {
	Action UseKind
	Houses HouseMatcher
}

// A perHouseKeySurcharge is a counted key-cost raise: Per extra Æmber for each
// creature of House in play, recomputed at each forge rather than frozen. It is
// comparable, so it rides in a Bar and in flat state; the zero value (HouseNone)
// surcharges nothing.
type perHouseKeySurcharge struct {
	House House
	Per   int
}
