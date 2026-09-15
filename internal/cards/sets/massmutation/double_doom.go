//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DoubleDoom
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Return an enemy creature to its owner's hand. Your opponent discards a random card from their hand.
var DoubleDoom = set.New(
	"Double Doom",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "020"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
