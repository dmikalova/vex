//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SauryAboutThat
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy a creature. Its controller gains 1A.
var SauryAboutThat = card.New(
	"Saury About That",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 228),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
