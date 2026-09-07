//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// HunterOrHunted
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Ward a creature, or move a ward from a creature to another creature.
var HunterOrHunted = card.New(
	"Hunter or Hunted?",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "269"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
