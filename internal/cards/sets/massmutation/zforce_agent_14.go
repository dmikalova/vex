package massmutation

import "github.com/dmikalova/vex/internal/card"

// zForceCluster binds the three Z- upgrades into one WholePool cluster led by
// Z-Force Agent 14: whenever the agent rolls into a pod, all three upgrades ride
// in with it, one of each (ADR 0036). The upgrades are Rarity.Connected, so the
// only way they reach a deck is on Z-Force Agent 14's coattails.
var zForceCluster = card.Cluster{
	Name:     "Z-Force Agent 14",
	Strategy: card.ClusterStrategy.WholePool,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Z-Force Agent 14
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Cyborg
//
//	Fight: For each upgrade on Z-Force Agent 14, gain 1 Æmber.
var ZForceAgent14 = set.New(
	"Z-Force Agent 14",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "353"),
	card.LeadsCluster(zForceCluster),
	card.WithPower(5),
	card.WithTraits(card.Traits.Cyborg),
	card.WithAbility(
		card.Trigger.Fight, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per:    card.UpgradesOn{Target: card.Target.This},
		}),
)
