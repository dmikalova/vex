//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Skoll
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Assault 3. (Before this creature attacks, deal 3D to the attacked enemy.)
//	After an enemy creature is destroyed by Skoll's assault damage, give a friendly creature a +1 power counter.
var Skoll = card.New(
	"Skoll",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "29"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
