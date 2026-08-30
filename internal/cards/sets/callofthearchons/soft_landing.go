//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SoftLanding
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//
//	Play: The next creature or artifact you play this turn enters play ready.
var SoftLanding = card.New(
	"Soft Landing",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 177),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
