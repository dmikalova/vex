//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LegionsMarch
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: For the remainder of the turn, after you use a Dinosaur creature, deal 1D to each non-Dinosaur creature.
var LegionsMarch = set.New(
	"Legion's March",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "224"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
