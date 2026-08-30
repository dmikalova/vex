//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// PileOfSkulls
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Each time an enemy creature is destroyed during your turn, a friendly creature captures 1 Aember.
var PileOfSkulls = card.New(
	"Pile of Skulls",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 25),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
