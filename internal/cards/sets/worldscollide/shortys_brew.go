package worldscollide

import "github.com/dmikalova/vex/internal/card"

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
//	Bonus:  Æmber
//
//	This creature gains +4 assault.
var ShortysBrew = set.New(
	"Shorty's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "68"),
	card.LeadsCluster(shortysBrewCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		AssaultBonus: 4,
	}),
)
