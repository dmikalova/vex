//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// PocketUniverse
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	You may spend Aember on Pocket Universe when forging keys.
//	Action: Move 1 Aember from your pool to Pocket Universe.
var PocketUniverse = card.New(
	"Pocket Universe",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 131),
	card.WithTraits("Item"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
