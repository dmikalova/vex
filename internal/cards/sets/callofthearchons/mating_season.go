//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MatingSeason
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Shuffle each Mars creature into its owner's deck. Each player gains 1 Aember for each creature shuffled into their deck this way.
var MatingSeason = card.New(
	"Mating Season",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 170),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
