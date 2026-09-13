package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Force Field
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains, "Reap: Ward this Creature."
var ForceField = card.New(
	"Force Field",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "310"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect:  card.Ward{Target: card.Target.This},
		}},
	}),
)
