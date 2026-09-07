//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// IntoTheNight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Until the start of your next turn, non-Shadows creatures cannot be used to fight.
var IntoTheNight = card.New(
	"Into the Night",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "256"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
