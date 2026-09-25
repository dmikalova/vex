package worldscollide

import "github.com/dmikalova/vex/internal/card"

// narpsBrewCluster pulls a Mega Narp into Narp's Brew's pod — a Pull cluster the
// brew leads, so the giant it is brewed for rides along.
var narpsBrewCluster = card.Cluster{
	Name:     "Narp's Brew",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Narp's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains +2 armor and taunt.
var NarpsBrew = set.New(
	"Narp's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "67"),
	card.LeadsCluster(narpsBrewCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		ArmorBonus: 2,
		Keywords:   card.Keywords(card.Keyword.Taunt),
	}),
)
