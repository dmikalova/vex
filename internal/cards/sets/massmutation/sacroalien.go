//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SacroAlien
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Armor:  2
//	Traits: Mutant • Alien
//
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
var SacroAlien = set.New(
	"Sacro-Alien",
	card.House.Staralliance,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "359"),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Alien),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
