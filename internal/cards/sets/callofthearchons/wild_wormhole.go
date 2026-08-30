//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// WildWormhole
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Play the top card of your deck.
var WildWormhole = card.New(
	"Wild Wormhole",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 125),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
