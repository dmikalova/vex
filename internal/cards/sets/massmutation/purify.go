//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Purify
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Purge a Mutant creature. If you do, discard cards from the top of its controller's deck until you discard a non-Mutant creature or run out of cards. If you discard a non-Mutant creature this way, put it into play under its owner's control.
var Purify = set.New(
	"Purify",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "153"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
