package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Fyre-Breath
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains +3 power.
//	This Creature gains, "Before Fight: Deal 2 damage to each neighbor of the Creature this Creature fights."
var FyreBreath = card.New(
	"Fyre-Breath",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "20"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		PowerBonus: 3,
		Granted: []card.Ability{{
			Trigger: card.Trigger.BeforeFight,
			Effect: card.DealDamage{
				Amount: 2,
				Target: card.Target.CreatureFought.NeighborsOf(),
			},
		}},
	}),
)
