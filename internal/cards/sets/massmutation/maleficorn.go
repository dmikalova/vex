//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Maleficorn
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Mutant
//
//	Enhance Damage Damage Damage Damage.
//	After an enemy creature is dealt damage by a D bonus icon, deal 1D to that creature.
var Maleficorn = set.New(
	"Maleficorn",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "040"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
