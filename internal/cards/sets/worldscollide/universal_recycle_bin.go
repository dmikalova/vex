//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// UniversalRecycleBin
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Action: Archive a purged card you own.
var UniversalRecycleBin = card.New(
	"Universal Recycle Bin",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 179),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
