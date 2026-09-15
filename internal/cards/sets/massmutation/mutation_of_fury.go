//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MutationOfFury
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Until the start of your next turn, a creature gains assault 3 and the Mutant trait.
var MutationOfFury = set.New(
	"Mutation of Fury",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MM, "414"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
