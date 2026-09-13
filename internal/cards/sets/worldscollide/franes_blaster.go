package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// franesBlasterCluster pulls a First Officer Frane into Frane's Blaster's pod — a
// Pull cluster, at least one averaging about one and a quarter (ADR 0036).
var franesBlasterCluster = card.Cluster{
	Name:     "Frane's Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Frane's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a Creature
//	- Attach Frane's Blaster to First Officer Frane -> move all Æmber from First Officer Frane to your pool."
var FranesBlaster = card.New(
	"Frane's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "346"),
	card.LeadsCluster(franesBlasterCluster),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First: card.AttachSelfTo{Host: FirstOfficerFrane.Name},
				Result: card.MoveAember{
					All:  true,
					From: card.Target.AttachedHost.Named(FirstOfficerFrane.Name),
					To:   card.Controller,
				},
			},
		}}),
	}),
)
