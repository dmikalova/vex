//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// InstrumentOfSilence
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains skirmish and, "Fight: Gain 1A."
var InstrumentOfSilence = card.New(
	"Instrument of Silence",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 375),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
