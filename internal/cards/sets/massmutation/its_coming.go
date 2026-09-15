//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ItsComing
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Search your deck and discard pile for either half of a gigantic creature, reveal it, and put it into your hand. Shuffle your deck.
var ItsComing = set.New(
	"It's Coming...",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MM, "117"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
