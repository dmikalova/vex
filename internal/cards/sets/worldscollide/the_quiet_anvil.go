//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheQuietAnvil
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
//	Keys cost -2A.
//	After a player forges a key, destroy The Quiet Anvil.
var TheQuietAnvil = card.New(
	"The Quiet Anvil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 282),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
