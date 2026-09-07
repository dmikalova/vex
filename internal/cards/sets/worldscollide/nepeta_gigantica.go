//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// NepetaGigantica
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Stun a creature with power 5 or higher, or stun a Giant creature.
var NepetaGigantica = card.New(
	"Nepeta Gigantica",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 394),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
