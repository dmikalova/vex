//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TrustNoOne
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1A. If there are no friendly creatures in play, instead steal 1A for each house represented among enemy creatures (to a maximum of 3).
var TrustNoOne = card.New(
	"Trust No One",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 248),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
