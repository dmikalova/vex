//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// PhaseShift
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: You may play one non-Logos card this turn.
var PhaseShift = card.New(
	"Phase Shift",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 117),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
