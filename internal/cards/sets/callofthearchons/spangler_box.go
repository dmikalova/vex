//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SpanglerBox
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Purge a creature in play. If you do, your opponent gains control of Spangler Box. If Spangler Box leaves play, return to play all cards purged by Spangler Box.
var SpanglerBox = card.New(
	"Spangler Box",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 132),
	card.WithTraits("Item"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
