//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ShieldOfJustice
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For the remainder of the turn, each friendly creature cannot be dealt damage.
var ShieldOfJustice = card.New(
	"Shield of Justice",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 225),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
