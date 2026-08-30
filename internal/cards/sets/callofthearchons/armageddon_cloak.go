package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Armageddon Cloak
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains +2 hazardous and, "If this creature would be destroyed, instead fully heal it, and destroy Armageddon Cloak."
var ArmageddonCloak = card.New(
	"Armageddon Cloak",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 263),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		HazardousBonus: 2,
		Replaces: card.Replace{
			When: card.Event.Destroyed,
			With: card.Sequence{
				Effects: []card.Effect{
					card.Heal{Fully: true, Target: card.Target.Triggering},
					card.Destroy{Target: card.Target.This},
				},
			},
		},
	}),
)
