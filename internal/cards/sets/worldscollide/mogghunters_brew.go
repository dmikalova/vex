package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// mogghuntersBrewCluster pulls a Mega Mogghunter into Mogghunter's Brew's pod — a
// Pull cluster the brew leads, so the giant it is brewed for rides along.
var mogghuntersBrewCluster = card.Cluster{
	Name:     "Mogghunter's Brew",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Mogghunter's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Fight: Deal 2 damage to a flank Creature."
var MogghuntersBrew = set.New(
	"Mogghunter's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "66"),
	card.LeadsCluster(mogghuntersBrewCluster),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Fight,
			Effect: card.DealDamage{
				Target: card.Target.Creature.OnFlank(),
				Amount: 2,
			},
		}},
	}),
)
