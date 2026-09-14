package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// chieftainsBrewCluster pulls a Mega Ganger Chieftain into Chieftain's Brew's
// pod — a Pull cluster the brew leads, so the giant it is brewed for rides along.
var chieftainsBrewCluster = card.Cluster{
	Name:     "Chieftain's Brew",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Chieftain's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Fight: Ready and fight with a neighboring Creature."
var ChieftainsBrew = set.New(
	"Chieftain's Brew",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "62"),
	card.LeadsCluster(chieftainsBrewCluster),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Fight,
			Effect: card.OnChooseCreature{
				Target: card.Target.Creature.Neighboring(),
				Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
			},
		}},
	}),
)
