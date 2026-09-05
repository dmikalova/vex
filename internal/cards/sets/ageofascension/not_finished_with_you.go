//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// NotFinishedWithYou
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Shuffle any number of creatures from your discard pile into your deck.
var NotFinishedWithYou = card.New(
	"Not Finished with You",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 63),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
