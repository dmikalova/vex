//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// AmphoraCaptura
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Enhance Aember Aember Damage Damage Draw Draw.
//	When resolving a bonus icon, you may choose to resolve it as a PT bonus icon instead.
var AmphoraCaptura = set.New(
	"Amphora Captura",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "215"),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
