//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MothershipSupport
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each friendly ready Mars creature, deal 2 Damage to a creature. (You may choose a different creature each time.)
var MothershipSupport = card.New(
	"Mothership Support",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 171),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
