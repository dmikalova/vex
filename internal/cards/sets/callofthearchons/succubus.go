package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Succubus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Demon
//
//	Your opponent's hand size is 1 less.
var Succubus = set.New(
	"Succubus",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, "99"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon),
	card.WithDrawModifier(card.Opponent, -1),
)
