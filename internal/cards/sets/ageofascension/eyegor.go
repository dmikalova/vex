//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Eyegor
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg
//
//	Play: Look at the top 3 cards of your deck. Add 1 to your hand and discard the others.
var Eyegor = card.New(
	"Eyegor",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 111),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
