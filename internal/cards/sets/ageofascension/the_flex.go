//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// TheFlex
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a ready friendly Brobnar creature. Exhaust it and gain A equal to half its power (rounding down the gain).
var TheFlex = card.New(
	"The Flex",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 31),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
