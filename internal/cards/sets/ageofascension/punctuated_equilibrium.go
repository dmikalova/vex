//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// PunctuatedEquilibrium
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Each player discards their hand, then refills their hand as if it were the end of their turn.
var PunctuatedEquilibrium = card.New(
	"Punctuated Equilibrium",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 363),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
