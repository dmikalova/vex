//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// RemoteAccess
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Use an opponent's artifact as if it were yours.
var RemoteAccess = card.New(
	"Remote Access",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 120),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
