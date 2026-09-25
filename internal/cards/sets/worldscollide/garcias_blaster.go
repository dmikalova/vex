package worldscollide

import "github.com/dmikalova/vex/internal/card"

// garciasBlasterCluster pulls a Sensor Chief Garcia into Garcia's Blaster's pod —
// a Pull cluster, at least one averaging about one and a quarter (ADR 0036).
var garciasBlasterCluster = card.Cluster{
	Name:     "Garcia's Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Garcia's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Garcia's Blaster to Sensor Chief Garcia -> steal 1 Æmber."
var GarciasBlaster = set.New(
	"Garcia's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "347"),
	card.LeadsCluster(garciasBlasterCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{
				Amount: 2,
				Target: card.Target.Creature,
			},
			card.Then{
				First: card.AttachSelfTo{
					Target: card.Target.FriendlyCreature.Named(SensorChiefGarcia.Name),
				},
				Result: card.StealAember{Amount: 1},
			},
		}}),
	}),
)
