//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ThePaleStar
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Omni: Destroy The Pale Star. For the remainder of the turn, each creature is considered to have 1 power and 0 armor.
var ThePaleStar = set.New(
	"The Pale Star",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "049"),
	card.WithTraits(card.Traits.Power),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
