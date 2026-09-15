//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SeekerOfTruth
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Armor:  1
//	Traits: Human
//
//	Fight: You may fight with a friendly non-Sanctum creature.
var SeekerOfTruth = set.New(
	"Seeker of Truth",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "138"),
	card.WithPower(3),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
