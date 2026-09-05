//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// YxlyxStimrager
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Martian • Soldier
//
//	Fight: Deal 2D to a creature and move it to either flank of its controller's battleline.
var YxlyxStimrager = card.New(
	"Yxlyx Stimrager",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 209),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
