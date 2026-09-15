package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Imperial Scutum
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains +2 armor.
//	This creature gains, "Destroyed: Move each Æmber on this creature to the common supply."
var ImperialScutum = set.New(
	"Imperial Scutum",
	card.House.Saurian,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, "185"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		ArmorBonus: 2,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect:  card.MoveAemberToSupply{All: true, Target: card.Target.This},
		}},
	}),
)
