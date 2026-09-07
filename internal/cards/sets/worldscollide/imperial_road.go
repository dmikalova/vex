//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ImperialRoad
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Omni: Play a Saurian creature. That creature enters play stunned.
var ImperialRoad = card.New(
	"Imperial Road",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 223),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
