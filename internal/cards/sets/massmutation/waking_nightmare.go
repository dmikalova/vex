//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// WakingNightmare
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Enhance Capture.
//	Play: Keys cost +1A for each Dis creature in play during your opponent's next turn.
var WakingNightmare = set.New(
	"Waking Nightmare",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "017"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
