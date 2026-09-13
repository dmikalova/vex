package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// molinasBlasterCluster pulls an Armsmaster Molina into Molina's Blaster's pod — a
// Pull cluster, at least one averaging about one and a quarter (ADR 0036).
var molinasBlasterCluster = card.Cluster{
	Name:     "Molina's Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Molina's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a Creature
//	- Attach Molina's Blaster to Armsmaster Molina -> deal 3 damage to a Creature."
var MolinasBlaster = card.New(
	"Molina's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "302"),
	card.LeadsCluster(molinasBlasterCluster),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: ArmsmasterMolina.Name},
				Result: card.DealDamage{Amount: 3, Target: card.Target.Creature},
			},
		}}),
	}),
)
