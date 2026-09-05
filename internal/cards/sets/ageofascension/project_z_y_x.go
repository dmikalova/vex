//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ProjectZYX
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Cyborg • Mutant
//
//	Fight/Reap: You may play one of your archived cards as if it were in your hand and in the active house.
var ProjectZYX = card.New(
	"Project Z.Y.X.",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 152),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
