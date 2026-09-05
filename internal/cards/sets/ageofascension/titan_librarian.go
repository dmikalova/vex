//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// TitanLibrarian
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	At the end of your turn, if Titan Librarian is not on a flank, archive a card.
var TitanLibrarian = card.New(
	"Titan Librarian",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 120),
	card.WithPower(4),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
