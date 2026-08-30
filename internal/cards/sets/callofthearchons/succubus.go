//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Succubus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
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
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
