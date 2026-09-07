//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// InformationExchange
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1A. If your opponent stole A from you on their previous turn, steal 2A instead.
var InformationExchange = card.New(
	"Information Exchange",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "136"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
