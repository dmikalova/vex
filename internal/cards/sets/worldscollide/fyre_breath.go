package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Fyre-Breath
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains +3 power.
//	This creature gains, "Before Fight: Deal 2 damage to each neighbor of the creature this creature fights."
var FyreBreath = set.New(
	"Fyre-Breath",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "20"),
	card.WithBonus(card.Bonus.Aember),
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
