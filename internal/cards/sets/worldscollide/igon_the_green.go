package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// IgonTheGreenName lets Igon the Terrible reference the Green by name without a
// package var-init cycle (the Green already references IgonTheTerrible.Name).
const IgonTheGreenName = "Igon the Green"

// igonCluster pulls one Igon the Terrible per Igon the Green: the two are a
// PullExact pair, so the Terrible (Rarity.Connected) rides in with the Green
// (ADR 0036).
var igonCluster = card.Cluster{
	Name:     "Igon the Green",
	Strategy: card.ClusterStrategy.PullExact,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Igon the Green
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	Destroyed: Purge Igon the Green, and put an Igon the Terrible from your discard pile into your hand.
var IgonTheGreen = card.New(
	IgonTheGreenName,
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "39"),
	card.LeadsCluster(igonCluster),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Destroyed, card.Sequence{Effects: []card.Effect{
			card.PurgeCreature{Target: card.Target.This},
			card.PutFromDiscard{
				Match:       card.Match{Name: IgonTheTerrible.Name},
				Destination: card.To.Hand,
			},
		}}),
)
