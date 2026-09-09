package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Plasma Nozzle
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Before Fight: Deal 2 damage to the creature this creature fought and 2 damage to each of its neighbors."
var PlasmaNozzle = card.New(
	"Plasma Nozzle",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "336"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.BeforeFight,
			Effect: card.DealDamage{Spread: card.CreatureAndNeighbors{
				Amount: 2,
				Splash: 2,
				Target: card.Target.CreatureFought,
			}},
		}},
	}),
)
