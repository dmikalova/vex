//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// FollowTheLeader
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For the remainder of the turn, each friendly creature may fight.
var FollowTheLeader = card.New(
	"Follow the Leader",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 8),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
