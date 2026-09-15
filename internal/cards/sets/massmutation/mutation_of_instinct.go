//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MutationOfInstinct
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Until the start of your next turn, a creature gains skirmish and the Mutant trait.
var MutationOfInstinct = set.New(
	"Mutation of Instinct",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MM, "415"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
