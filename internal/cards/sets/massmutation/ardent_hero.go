//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ArdentHero
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Taunt.
//	Ardent Hero cannot be dealt damage by Mutant creatures or creatures with power 5 or higher.
var ArdentHero = set.New(
	"Ardent Hero",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "126"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
