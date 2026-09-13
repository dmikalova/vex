package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// qincansBlasterCluster pulls a Sci Officer Qincan into Qincan's Blaster's pod — a
// Pull cluster, at least one averaging about one and a quarter (ADR 0036).
var qincansBlasterCluster = card.Cluster{
	Name:     "Qincan's Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Qincan's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a Creature
//	- Attach Qincan's Blaster to Sci. Officer Qincan -> archive a Creature from play."
var QincansBlaster = card.New(
	"Qincan's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "351"),
	card.LeadsCluster(qincansBlasterCluster),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: SciOfficerQincan.Name},
				Result: card.ArchiveFromPlay{Target: card.Target.Creature},
			},
		}}),
	}),
)
