//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TwinBoltEmission
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 Damage to a creature and deal 2 Damage to a different creature.
var TwinBoltEmission = card.New(
	"Twin Bolt Emission",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 124),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
