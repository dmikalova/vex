//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// FangtoothCavern
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	At the end of your turn, destroy the least powerful creature.
var FangtoothCavern = card.New(
	"Fangtooth Cavern",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 370),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
