//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mother
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Robot • Scientist
//
//	During your "draw cards" step, refill your hand to 1 additional card.
var Mother = card.New(
	"Mother",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 145),
	card.WithPower(5),
	card.WithTraits("Robot", "Scientist"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
