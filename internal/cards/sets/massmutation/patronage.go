//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Patronage
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Move half the A from a creature to your pool (rounding up). Move the remaining A from that creature to your opponent's pool.
var Patronage = set.New(
	"Patronage",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "227"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
