package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Firespitter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Giant
//
//	Before Fight: Deal 1 damage to each enemy creature.
var Firespitter = card.New(
	"Firespitter",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 32),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.BeforeFight, card.DealDamage{
			Amount: 1,
			Target: card.Target.EachEnemyCreature,
		}),
)
