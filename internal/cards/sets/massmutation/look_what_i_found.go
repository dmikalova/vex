//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LookWhatIFound
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Omega.
//	Play: Return 1 card of each type (action, artifact, creature, upgrade) from your discard pile to your hand.
var LookWhatIFound = set.New(
	"Look What I Found!",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "402"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
