//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// RedHotArmor
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Each enemy creature with armor loses all of its armor until the end of the turn and is dealt 1 Damage for each point of armor it lost this way.
var RedHotArmor = card.New(
	"Red-Hot Armor",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 70),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
