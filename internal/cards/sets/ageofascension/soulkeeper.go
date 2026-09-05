package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Soulkeeper
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Destroyed: Destroy the 1 most powerful enemy creatures."
var Soulkeeper = card.New(
	"Soulkeeper",
	card.House.Dis,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 83),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect: card.Destroy{
				Target: card.Target.EachEnemyCreature.Selector(card.MostPowerful(1)),
			},
		}},
	}),
)
