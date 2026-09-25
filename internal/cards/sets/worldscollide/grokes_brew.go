package worldscollide

import "github.com/dmikalova/vex/internal/card"

// grokesBrewCluster pulls a Mega Groke into Groke's Brew's pod — a Pull cluster
// the brew leads, so the giant it is brewed for rides along.
var grokesBrewCluster = card.Cluster{
	Name:     "Groke's Brew",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Groke's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight: Your opponent loses 1 Æmber."
var GrokesBrew = set.New(
	"Groke's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "64"),
	card.LeadsCluster(grokesBrewCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Fight,
			Effect: card.LoseAember{
				Player: card.Opponent,
				Amount: 1,
			},
		}},
	}),
)
