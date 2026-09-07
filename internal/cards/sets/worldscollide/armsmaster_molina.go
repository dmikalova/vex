//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ArmsmasterMolina
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Hazardous 3. (Before this creature is attacked, deal 3D to the attacking enemy.)
//	Each of Armsmaster Molina's neighbors gains hazardous 3.
var ArmsmasterMolina = card.New(
	"Armsmaster Molina",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 292),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
