package worldscollide

import "github.com/dmikalova/vex/internal/card"

// khrkharsBlasterCluster pulls a Lieutenant Khrkhar into Khrkhar's Blaster's pod —
// a Pull cluster, at least one averaging about one and a quarter (ADR 0036).
var khrkharsBlasterCluster = card.Cluster{
	Name:     "Khrkhar's Blaster",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Khrkhar's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Khrkhar's Blaster to Lieutenant Khrkhar -> ward Lieutenant Khrkhar."
var KhrkharsBlaster = set.New(
	"Khrkhar's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "349"),
	card.LeadsCluster(khrkharsBlasterCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{
				Amount: 2,
				Target: card.Target.Creature,
			},
			card.Then{
				First: card.AttachSelfTo{
					Target: card.Target.FriendlyCreature.Named(LieutenantKhrkhar.Name),
				},
				Result: card.Ward{Target: card.Target.AttachedHost.Named(LieutenantKhrkhar.Name)},
			},
		}}),
	}),
)
