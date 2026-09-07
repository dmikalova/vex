//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// PhalanxStrike
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a creature. Deal 1D to it for each friendly creature. You may exalt a friendly creature to repeat the preceding effect.
var PhalanxStrike = card.New(
	"Phalanx Strike",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 189),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
