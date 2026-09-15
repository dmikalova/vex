package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// wallsBlasterCluster pulls a Chief Engineer Walls into Walls' Blaster's pod — a
// Pull cluster, at least one averaging about one and a quarter (ADR 0036).
var wallsBlasterCluster = card.Cluster{
	Name:     "Walls' Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Walls' Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Walls' Blaster to Chief Engineer Walls -> for each upgrade on Chief Engineer Walls, stun a creature."
var WallsBlaster = set.New(
	"Walls' Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "352"),
	card.LeadsCluster(wallsBlasterCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First: card.AttachSelfTo{Host: ChiefEngineerWalls.Name},
				Result: card.ForEach{
					Times: card.UpgradesOn{
						Target: card.Target.AttachedHost.Named(ChiefEngineerWalls.Name),
					},
					Do: card.Stun{Target: card.Target.Creature},
				},
			},
		}}),
	}),
)
