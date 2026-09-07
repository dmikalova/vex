//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GeneralOrder24
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	At the start of each player's turn, they must choose a creature they control and destroy each creature of the chosen creature's house. If that player has no creatures in play, destroy General Order 24 instead.
var GeneralOrder24 = card.New(
	"General Order 24",
	card.House.Staralliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "333"),
	card.WithTraits(card.Traits.Law),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
