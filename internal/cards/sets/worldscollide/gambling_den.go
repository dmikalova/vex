//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GamblingDen
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	At the start of each player's turn, that player may name a house. If they do, reveal the top card of their deck. If it is of the named house, they gain 2A. Otherwise, they lose 2A.
var GamblingDen = card.New(
	"Gambling Den",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 268),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
