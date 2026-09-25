package worldscollide

import "github.com/dmikalova/vex/internal/card"

// cowfynesBrewCluster pulls a Mega Cowfyne into Cowfyne's Brew's pod — a Pull
// cluster the brew leads, so the giant it is brewed for rides along.
var cowfynesBrewCluster = card.Cluster{
	Name:     "Cowfyne's Brew",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Cowfyne's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains +2 splash-attack.
var CowfynesBrew = set.New(
	"Cowfyne's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "63"),
	card.LeadsCluster(cowfynesBrewCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		SplashAttackBonus: 2,
	}),
)
