package worldscollide

import "github.com/dmikalova/vex/internal/card"

// ingramsBlasterCluster pulls a Medic Ingram into Ingram's Blaster's pod — a Pull
// cluster, at least one averaging about one and a quarter (ADR 0036).
var ingramsBlasterCluster = card.Cluster{
	Name:     "Ingram's Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Ingram's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Ingram's Blaster to Medic Ingram -> fully heal a creature."
var IngramsBlaster = set.New(
	"Ingram's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "348"),
	card.LeadsCluster(ingramsBlasterCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{
				Amount: 2,
				Target: card.Target.Creature,
			},
			card.Then{
				First: card.AttachSelfTo{
					Target: card.Target.FriendlyCreature.Named(MedicIngram.Name),
				},
				Result: card.Heal{
					Fully:  true,
					Target: card.Target.Creature,
				},
			},
		}}),
	}),
)
