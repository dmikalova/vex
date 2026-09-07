//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// AutoLegionary
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Robot • Ally
//
//	Action: Put Auto-Legionary on a flank of your battleline. While in the battleline, it is considered a creature with 5 power and belongs to all houses.
var AutoLegionary = card.New(
	"Auto-Legionary",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "214"),
	card.WithTraits(card.Traits.Robot, card.Traits.Ally),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
