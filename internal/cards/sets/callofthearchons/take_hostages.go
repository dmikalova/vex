//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TakeHostages
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For the remainder of the turn, each time a friendly creature fights, it captures 1 Aember.
var TakeHostages = card.New(
	"Take Hostages",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 226),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
