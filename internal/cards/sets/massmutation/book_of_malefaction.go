//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BookOfMalefaction
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item • Law
//
//	After your A is stolen, put a warrant counter on Book of Malefaction for each A stolen.
//	Omni: Remove a warrant counter from Book of Malefaction. If you do, purge a creature.
var BookOfMalefaction = set.New(
	"Book of Malefaction",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "159"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item, card.Traits.Law),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
