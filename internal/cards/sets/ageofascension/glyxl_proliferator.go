//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// GlyxlProliferator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Soldier
//
//	Reap: If Glyxl Proliferator is on a flank, archive a Mars card from your discard pile.
var GlyxlProliferator = card.New(
	"Glyxl Proliferator",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 163),
	card.WithPower(3),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
