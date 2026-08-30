//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// CustomVirus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	Omni: Sacrifice Custom Virus. Purge a creature from your hand. Destroy each creature that shares a trait with the purged creature.
var CustomVirus = card.New(
	"Custom Virus",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 183),
	card.WithAemberBonus(1),
	card.WithTraits("Weapon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
