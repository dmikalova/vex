//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// PerplexingSophistry
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If you have more A than your opponent, they discard a random card from their hand and you draw a card.
var PerplexingSophistry = card.New(
	"Perplexing Sophistry",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 293),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
