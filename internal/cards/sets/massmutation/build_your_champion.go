//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BuildYourChampion
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Search your deck and discard pile for two halves of a gigantic creature, reveal them, and archive them.
var BuildYourChampion = set.New(
	"Build Your Champion",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "002"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
