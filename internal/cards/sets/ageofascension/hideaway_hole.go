//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// HideawayHole
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	Omni: Sacrifice Hideaway Hole. Creatures you control gain elusive until the start of your next turn.
var HideawayHole = card.New(
	"Hideaway Hole",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 287),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
