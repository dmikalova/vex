package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// alakasBrewCluster pulls a Mega Alaka into Alaka's Brew's pod — a Pull cluster
// the brew leads, so the giant it is brewed for rides along.
var alakasBrewCluster = card.Cluster{
	Name:     "Alaka's Brew",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Alaka's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight: Play a creature -> ready it."
var AlakasBrew = set.New(
	"Alaka's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "2"),
	card.LeadsCluster(alakasBrewCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Fight,
			Effect: card.Then{
				First:  card.PlayFrom{From: card.Hand, Types: card.Types.Of(card.Type.Creature)},
				Result: card.Ready{Target: card.Target.Triggering},
			},
		}},
	}),
)
