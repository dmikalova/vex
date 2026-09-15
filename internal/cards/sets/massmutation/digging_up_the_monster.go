//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DiggingUpTheMonster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Search your deck and discard pile for two halves of a gigantic creature and reveal them. Put them on top of your deck in any order.
var DiggingUpTheMonster = set.New(
	"Digging Up the Monster",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "003"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
