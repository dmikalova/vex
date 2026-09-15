//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MonumentToOctavia
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Action: A friendly creature captures 1A. If Cornicen Octavia is in your discard pile, that creature captures 2A instead.
var MonumentToOctavia = set.New(
	"Monument to Octavia",
	card.House.Saurian,
	card.Type.Artifact,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "237"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
