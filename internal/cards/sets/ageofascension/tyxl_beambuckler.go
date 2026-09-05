//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// TyxlBeambuckler
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Martian • Soldier
//
//	Play: Deal 2D to a creature and move it to either flank of its controller's battleline.
var TyxlBeambuckler = card.New(
	"Tyxl Beambuckler",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 171),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
