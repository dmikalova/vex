package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Charybdis
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Beast
//
//	Each enemy creature gains, "Before Fight: Lose 1 Æmber."
var Charybdis = set.New(
	"Charybdis",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.MM, "234"),
	card.InCluster(scyllaCluster),
	card.WithPower(7),
	card.WithTraits(card.Traits.Beast),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachEnemyCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.BeforeFight,
			Effect: card.LoseAember{
				Player: card.Controller,
				Amount: 1,
			},
		}},
	}),
)
