//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MonumentToFaust
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Action: Keys cost +1A during your opponent's next turn. If Faust the Great is in your discard pile, keys cost +2A during your opponent's next turn instead.
var MonumentToFaust = set.New(
	"Monument to Faust",
	card.House.Saurian,
	card.Type.Artifact,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "236"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
