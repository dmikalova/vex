//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Pincerator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	At the end of each turn, deal 1D to each flank creature.
var Pincerator = set.New(
	"Pincerator",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "289"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
