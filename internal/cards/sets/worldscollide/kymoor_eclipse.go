//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// KymoorEclipse
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Shuffle each flank creature into its owner's deck.
var KymoorEclipse = card.New(
	"Kymoor Eclipse",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 243),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
