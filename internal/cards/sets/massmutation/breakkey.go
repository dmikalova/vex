//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Breakkey
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If your opponent has more forged keys than you, unforge an opponent's key. If you unforge an opponent's key this way, your opponent gains 6A.
var Breakkey = set.New(
	"Break-key",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "019"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
