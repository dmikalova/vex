//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MonumentToLudo
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Action: Move 1A from a creature to the common supply. If Praefectus Ludo is in your discard pile, move 2A from that creature to the common supply instead.
var MonumentToLudo = set.New(
	"Monument to Ludo",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.MM, "194"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
