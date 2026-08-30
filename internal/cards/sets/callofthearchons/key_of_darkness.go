//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// KeyOfDarkness
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Forge a key at +6 Aember current cost. If your opponent has no Aember, forge a key at +2 Aember current cost instead.
var KeyOfDarkness = card.New(
	"Key of Darkness",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 273),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
