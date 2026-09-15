//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// EssenceScale
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Action: Destroy a friendly creature. If you do, ready and use a friendly creature that shares a house with the destroyed creature.
var EssenceScale = set.New(
	"Essence Scale",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "021"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
