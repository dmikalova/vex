//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Jargogle
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Beast • Mutant
//
//	Elusive.
//	Play: Put a card from your hand facedown under Jargogle.
//	Destroyed: If it is your turn, play the card under Jargogle; otherwise, archive that card.
var Jargogle = card.New(
	"Jargogle",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 131),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast, card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
