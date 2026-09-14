package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// shortysBrewCluster pulls a Mega Shorty into Shorty's Brew's pod — a Pull
// cluster the brew leads, so the giant it is brewed for rides along.
var shortysBrewCluster = card.Cluster{
	Name:     "Shorty's Brew",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Shorty's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains +4 assault.
var ShortysBrew = set.New(
	"Shorty's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "68"),
	card.LeadsCluster(shortysBrewCluster),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		AssaultBonus: 4,
	}),
)
