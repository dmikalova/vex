//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// SackOfCoins
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Deal 1D to a creature for each
//	A in your pool.
var SackOfCoins = card.New(
	"Sack of Coins",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 312),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
