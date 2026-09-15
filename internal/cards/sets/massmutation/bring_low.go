//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BringLow
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Enhance Capture.
//	Play: Capture all but 5 of your opponent's A, distributed among any number of friendly creatures.
var BringLow = set.New(
	"Bring Low",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "147"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
