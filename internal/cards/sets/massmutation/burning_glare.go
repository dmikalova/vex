//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BurningGlare
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Enhance Damage.
//	Play: Stun an enemy creature, or stun each enemy Mutant creature.
var BurningGlare = set.New(
	"Burning Glare",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "128"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
