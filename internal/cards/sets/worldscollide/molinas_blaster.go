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
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Molina's Blaster to Armsmaster Molina -> deal 3 damage to a creature."
var MolinasBlaster = set.New(
	"Molina's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "302"),
	card.LeadsCluster(molinasBlasterCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First: card.AttachSelfTo{
					Target: card.Target.FriendlyCreature.Named(ArmsmasterMolina.Name),
				},
				Result: card.DealDamage{Amount: 3, Target: card.Target.Creature},
			},
		}}),
	}),
)
