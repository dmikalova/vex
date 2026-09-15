//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// GrowthSurge
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Give a flank creature three
//	+1 power counters. Give its neighbor two +1 power counters. Give the second creature's other neighbor
//	a +1 power counter.
var GrowthSurge = set.New(
	"Growth Surge",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "383"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
