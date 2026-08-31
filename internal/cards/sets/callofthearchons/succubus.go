package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Succubus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Demon
//
//	During their "draw cards" step, your opponent refills their hand to 1 less card.
var Succubus = card.New(
	"Succubus",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 99),
	card.WithPower(3),
	card.WithTraits("Demon"),
	card.WithDrawModifier(card.Opponent, -1),
)
