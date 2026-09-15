//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ReclaimedByNature
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Purge an artifact. Resolve its bonus icons as if you had played it.
var ReclaimedByNature = set.New(
	"Reclaimed by Nature",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "374"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
