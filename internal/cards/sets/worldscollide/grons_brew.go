package worldscollide

import "github.com/dmikalova/vex/internal/card"

// gronsBrewCluster pulls a Mega Gron Nine-Toes into Gron's Brew's pod — a Pull
// cluster the brew leads, so the giant it is brewed for rides along.
var gronsBrewCluster = card.Cluster{
	Name:     "Gron's Brew",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Gron's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains +4 power.
var GronsBrew = set.New(
	"Gron's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "65"),
	card.LeadsCluster(gronsBrewCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		PowerBonus: 4,
	}),
)
