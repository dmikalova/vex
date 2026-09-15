//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MonumentToShrix
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	You may spend A on Monument to Shrix as if it were in your pool.
//	Action: Move 1A from your pool to Monument to Shrix. If Citizen Shrix is in your discard pile, move 1A from any player's pool to Monument to Shrix instead.
var MonumentToShrix = set.New(
	"Monument to Shrix",
	card.House.Saurian,
	card.Type.Artifact,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "239"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
