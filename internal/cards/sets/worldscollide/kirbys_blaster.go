package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// kirbysBlasterCluster pulls a Com Officer Kirby into Kirby's Blaster's pod — a
// Pull cluster, at least one averaging about one and a quarter (ADR 0036).
var kirbysBlasterCluster = card.Cluster{
	Name:     "Kirby's Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Kirby's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a Creature
//	- Attach Kirby's Blaster to Com. Officer Kirby -> draw 2 cards."
var KirbysBlaster = card.New(
	"Kirby's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "350"),
	card.LeadsCluster(kirbysBlasterCluster),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: ComOfficerKirby.Name},
				Result: card.Draw{Amount: 2},
			},
		}}),
	}),
)
