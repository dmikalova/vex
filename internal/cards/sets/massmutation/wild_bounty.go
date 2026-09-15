//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// WildBounty
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Enhance Aember Aember.
//	Play: The next time you play a card this turn, resolve each of its bonus icons an additional time.
var WildBounty = set.New(
	"Wild Bounty",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "392"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
