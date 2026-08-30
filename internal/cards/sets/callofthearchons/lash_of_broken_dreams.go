//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// LashOfBrokenDreams
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Keys cost +3 Aember during your opponent's next turn.
var LashOfBrokenDreams = card.New(
	"Lash of Broken Dreams",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, 75),
	card.WithTraits("Weapon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
