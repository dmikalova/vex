//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// SongOfSpring
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Shuffle any number of friendly Untamed creatures from your hand, discard pile, or battleline back into your deck.
var SongOfSpring = card.New(
	"Song of Spring",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 332),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
