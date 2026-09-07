//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheBigOne
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	After a creature is played, put a fuse counter on The Big One.
//	If there are 10 or more fuse counters on The Big One, destroy each creature and artifact.
var TheBigOne = card.New(
	"The Big One",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "50"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Weapon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
