//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// NoSafetyInNumbers
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 3D to each creature that belongs to a house that has 3 or more creatures in play.
var NoSafetyInNumbers = card.New(
	"No Safety in Numbers",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 257),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
