//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MatterMaker
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	You may play upgrades as if they were in the active house.
var MatterMaker = set.New(
	"Matter Maker",
	card.House.Staralliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "349"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
