package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Brammo
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Giant • Knight
//
//	Play: Deal 2 damage to each enemy flank creature.
var Brammo = card.New(
	"Brammo",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 3),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Giant, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.EachEnemyCreature.OnFlank(),
		}),
)
