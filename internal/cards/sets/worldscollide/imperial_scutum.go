package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Imperial Scutum
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This Creature gains +2 armor.
//	This Creature gains, "Destroyed: Move each Æmber on this Creature to the common supply."
var ImperialScutum = card.New(
	"Imperial Scutum",
	card.House.Saurian,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, "185"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		ArmorBonus: 2,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect:  card.MoveAemberToSupply{All: true, Target: card.Target.This},
		}},
	}),
)
