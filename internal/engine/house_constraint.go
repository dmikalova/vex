package engine

// This file defines the delayed constraint table that governs a player's next
// active-house choice (ADR 0035). A player's forward house constraints are one
// table of typed entries — musts, cannots, and armed wagers — not three paired
// state slots. Entries accumulate (a must and a cannot stack; cannot overrides
// must), StartTurn promotes a player's armed entries onto their own next turn,
// and ChooseHouse resolves the whole table at once (see allowedHouses in
// game_read.go and the arming methods in game_turn.go).

// maxHouseConstraints bounds how many active-house constraints can bind one
// player at once — musts, cannots, and armed wagers together. It is generous: no
// realistic board arms this many forward house constraints on a single player.
const maxHouseConstraints = 8

// houseConstraintKind is what a HouseConstraint says about a player's next house
// choice.
type houseConstraintKind uint8

const (
	// houseConstraintUnset is the zero value: an empty table slot.
	houseConstraintUnset houseConstraintKind = iota
	// constraintMustHouse requires the player choose House — a fixed house the
	// arming card already named (Control the Weak).
	constraintMustHouse
	// constraintMustCreature requires the player choose the house of Creature, read
	// live at choice time (Snag — the fought creature's current house, which may
	// change between the fight and the choice).
	constraintMustCreature
	// constraintCannotHouse bars the player from choosing House (Tezmal, Snag's
	// Mirror).
	constraintCannotHouse
	// constraintWager is an armed bet: if the player chooses House, Predictor steals
	// Amount (Snaglet). It does not constrain the choice — it is a reaction to it,
	// riding the same next-turn promotion so the bet outlives the arming card.
	constraintWager
)

// A HouseConstraint is one forward constraint on a player's next active-house
// choice: what it says (Kind), the house or creature it reads, and the card that
// imposed it. It is comparable so it rides in the flat, value-copyable GameState
// (ADR 0005); the zero value is an empty table slot.
type HouseConstraint struct {
	// Kind is what the entry says — a must, a cannot, or an armed wager.
	Kind houseConstraintKind
	// House is the fixed house a must or cannot names, or a wager's predicted house.
	House House
	// Creature is the card a constraintMustCreature reads its houses from at choice
	// time (Snag stores the fought creature, not a frozen house).
	Creature LocalID
	// Amount is how much a winning wager steals.
	Amount int
	// Predictor is the player a winning wager pays.
	Predictor int
	// Source is the card that imposed the constraint, so a reminder can name it.
	Source LocalID
}
