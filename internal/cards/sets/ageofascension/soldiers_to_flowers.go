//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// SoldiersToFlowers
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Purge each Untamed creature from each player's discard pile. For each card purged this way, its owner gains 1A.
var SoldiersToFlowers = card.New(
	"Soldiers to Flowers",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 349),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
