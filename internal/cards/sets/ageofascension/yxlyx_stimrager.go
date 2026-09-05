package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Yxlyx Stimrager
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Martian • Soldier
//
//	Fight: Deal 2 damage to a creature and move it to either flank of its controller's battleline.
var YxlyxStimrager = card.New(
	"Yxlyx Stimrager",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 209),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.Fight, card.DamageThen{
			Amount: 2,
			Target: card.Target.Creature,
			Then:   card.MoveToFlank{Target: card.Target.Triggering},
		}),
)
