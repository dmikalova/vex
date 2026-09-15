//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// TemptingOffer
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Enhance Capture.
//	Play: Return an enemy creature to its owner's hand. If you do, your opponent gains 1A.
var TemptingOffer = set.New(
	"Tempting Offer",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "259"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
