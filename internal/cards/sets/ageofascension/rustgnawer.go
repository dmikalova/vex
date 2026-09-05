//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Rustgnawer
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast • Insect
//
//	Fight: Destroy an artifact. If that artifact had an Aember bonus, you gain that much A.
var Rustgnawer = card.New(
	"Rustgnawer",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 330),
	card.WithPower(4),
	card.WithTraits(card.Traits.Beast, card.Traits.Insect),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
