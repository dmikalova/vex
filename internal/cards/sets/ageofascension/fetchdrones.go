//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Fetchdrones
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Discard the top 2 cards of your deck. For each Logos card discarded this way, a friendly creature captures 2A.
var Fetchdrones = card.New(
	"Fetchdrones",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 144),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
