//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheGoldenSpiral
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Action: Exalt a friendly creature. Ready and use that creature.
var TheGoldenSpiral = card.New(
	"The Golden Spiral",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.WC, 194),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
