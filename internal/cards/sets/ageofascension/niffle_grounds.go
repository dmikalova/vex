//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// NiffleGrounds
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
//	Action: Choose a creature. For the remainder of the turn, that creature loses taunt and elusive.
var NiffleGrounds = card.New(
	"Niffle Grounds",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 346),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
