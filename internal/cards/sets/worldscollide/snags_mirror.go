//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SnagsMirror
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	After a player chooses an active house, their opponent cannot choose the same house as their active house on their next turn.
var SnagsMirror = card.New(
	"Snag's Mirror",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "117"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
