package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Force Field
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Reap: Ward this creature."
var ForceField = set.New(
	"Force Field",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "310"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect:  card.Ward{Target: card.Target.This},
		}},
	}),
)
