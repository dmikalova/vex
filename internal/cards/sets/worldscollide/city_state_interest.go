//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CityStateInterest
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Each friendly creature captures 1A.
var CityStateInterest = card.New(
	"City-State Interest",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 200),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
