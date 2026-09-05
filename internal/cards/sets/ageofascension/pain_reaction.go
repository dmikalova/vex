//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// PainReaction
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 2D to an enemy creature. If this damage destroys that creature, deal 2D to each of that creature's neighbors.
var PainReaction = card.New(
	"Pain Reaction",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 78),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
