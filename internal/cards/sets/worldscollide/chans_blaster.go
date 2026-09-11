package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// chansBlasterCluster pulls a Commander Chan into Chan's Blaster's pod — a Pull
// cluster, at least one averaging about one and a quarter (ADR 0036). The officer
// is a Common that rolls on its own too, so this tops the pod up.
var chansBlasterCluster = card.Cluster{
	Name:     "Chan's Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Chan's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Chan's Blaster to Commander Chan -> use another creature."
var ChansBlaster = card.New(
	"Chan's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "345"),
	card.LeadsCluster(chansBlasterCluster),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: CommanderChan.Name},
				Result: card.Use{Max: 1, Target: card.Target.OtherFriendlyCreature},
			},
		}}),
	}),
)
