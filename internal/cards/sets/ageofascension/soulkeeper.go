package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Soulkeeper
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Destroyed: Destroy the most powerful enemy creature."
var Soulkeeper = set.New(
	"Soulkeeper",
	card.House.Dis,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "83"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect: card.Destroy{
				Target: card.Target.EachEnemyCreature.Refine(card.MostPowerful),
			},
		}},
	}),
)
