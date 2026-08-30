//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// PsychicNetwork
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Steal 1 Aember for each friendly ready Mars creature.
var PsychicNetwork = card.New(
	"Psychic Network",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 174),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
