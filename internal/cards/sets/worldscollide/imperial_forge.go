//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ImperialForge
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Forge a key at +8A current cost, reduced by 1A for each A on friendly creatures.
var ImperialForge = card.New(
	"Imperial Forge",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 222),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
