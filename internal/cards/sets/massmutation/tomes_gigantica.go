//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// TomesGigantica
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Search your deck and discard pile for two halves of a gigantic creature, reveal them, and put them in your hand. Purge Tomes Gigantica.
var TomesGigantica = set.New(
	"Tomes Gigantica",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "004"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
