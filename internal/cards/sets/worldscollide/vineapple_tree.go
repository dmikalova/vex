//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// VineappleTree
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	Keys cost +1A for each growth counter on Vineapple Tree.
//	After a player forges a key, remove each growth counter from Vineapple Tree.
//	Action: Put a growth counter on Vineapple Tree.
var VineappleTree = card.New(
	"Vineapple Tree",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "402"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
