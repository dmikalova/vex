package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Zysysyx Shockworm
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Armor:  1
//	Traits: Martian • Soldier
//
//	After an enemy creature reaps, stun it.
var ZysysyxShockworm = card.New(
	"Zysysyx Shockworm",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 175),
	card.WithPower(3),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.AfterEnemyCreatureReaps, card.Stun{Target: card.Target.Triggering}),
)
