package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Tyxl Beambuckler
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Martian • Soldier
//
//	Play: Deal 2 damage to a creature and move it to either flank of its controller's battleline.
var TyxlBeambuckler = card.New(
	"Tyxl Beambuckler",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 171),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.Play, card.DamageThen{
			Amount: 2,
			Target: card.Target.Creature,
			Then:   card.MoveToFlank{Target: card.Target.Triggering},
		}),
)
