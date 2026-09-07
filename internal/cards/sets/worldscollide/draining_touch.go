//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DrainingTouch
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy a creature with no A on it.
var DrainingTouch = card.New(
	"Draining Touch",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 72),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
