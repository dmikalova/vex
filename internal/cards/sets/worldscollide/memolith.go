//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Memolith
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Action: Graft an action card from your hand onto Memolith (place it faceup under this card), or trigger the play effect of an action card grafted onto Memolith.
var Memolith = card.New(
	"Memolith",
	card.House.Brobnar,
	card.Type.Artifact,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A04"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
