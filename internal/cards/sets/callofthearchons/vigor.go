//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Vigor
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Heal up to 3 damage from a creature. If you healed 3 damage, gain 1 Aember.
var Vigor = card.New(
	"Vigor",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 338),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
