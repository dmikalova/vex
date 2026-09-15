package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Scylla
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Beast
//
//	Each enemy creature gains, "Reap: Deal 4 damage to this creature."
var Scylla = set.New(
	"Scylla",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "230"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Beast),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachEnemyCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect: card.DealDamage{
				Amount: 4,
				Target: card.Target.This,
			},
		}},
	}),
)
