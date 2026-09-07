//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CityGates
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: A friendly creature captures 1A. If that creature is a Dinosaur, it captures 2A instead.
var CityGates = card.New(
	"City Gates",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "216"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
