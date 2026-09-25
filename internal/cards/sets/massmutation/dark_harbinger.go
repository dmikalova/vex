package massmutation

import "github.com/dmikalova/vex/internal/card"

// darkHarbingerCluster binds the three Mutation tactics into one RandomCount
// cluster led by Dark Harbinger: whenever the Harbinger rolls into a pod, one to
// three of its Mutations ride in with it (ADR 0036), feeding the Harbinger's own
// "after you play an Untamed tactic, ready" payoff. The Mutations are
// Rarity.Connected, so the only way they reach a deck is on Dark Harbinger's
// coattails.
var darkHarbingerCluster = card.Cluster{
	Name:     "Dark Harbinger",
	Strategy: card.ClusterStrategy.RandomCount,
	Trigger:  card.ClusterTrigger.ByLead,
	Min:      1,
	Max:      3,
}

// Dark Harbinger
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant • Witch
//
//	After you play an Untamed tactic, ready Dark Harbinger.
var DarkHarbinger = set.New(
	"Dark Harbinger",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "381"),
	card.LeadsCluster(darkHarbingerCluster),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Witch),
	card.WithAbility(card.Trigger.AfterCardPlayed, card.Conditional{
		Cond: card.ItIs{
			House: card.Houses.Named(card.House.Self),
			Type:  card.Type.Tactic,
		},
		Then: card.Ready{Target: card.Target.This},
	}),
)
