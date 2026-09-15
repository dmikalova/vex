//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MonumentToPrimus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Action: Move 1A from a friendly creature to another friendly creature. If Consul Primus is in your discard pile, move 1A from a creature to another creature instead.
var MonumentToPrimus = set.New(
	"Monument to Primus",
	card.House.Saurian,
	card.Type.Artifact,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "238"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
