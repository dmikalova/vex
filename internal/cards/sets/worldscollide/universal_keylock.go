//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// UniversalKeylock
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Keys cost +3A.
//	After a player forges a key, destroy Universal Keylock.
var UniversalKeylock = card.New(
	"Universal Keylock",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "178"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
