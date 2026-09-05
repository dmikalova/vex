//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// RiggedLottery
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Each player discards the top 5 cards of their deck. For each Shadows card discarded, its owner gains 1A.
var RiggedLottery = card.New(
	"Rigged Lottery",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 309),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
