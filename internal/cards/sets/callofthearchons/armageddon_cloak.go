package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Armageddon Cloak
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains +2 hazardous and, "If this creature would be destroyed, instead fully heal it, and destroy Armageddon Cloak."
var ArmageddonCloak = set.New(
	"Armageddon Cloak",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "263"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		HazardousBonus: 2,
		Replaces: card.Replace{
			When: card.Event.Destroyed,
			With: card.Sequence{
				Effects: []card.Effect{
					card.Heal{
						Fully:  true,
						Target: card.Target.Triggering,
					},
					card.Destroy{Target: card.Target.This},
				},
			},
		},
	}),
)
