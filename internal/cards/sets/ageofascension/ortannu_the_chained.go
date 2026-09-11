package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ortannuCluster pulls a couple of Ortannu's Bindings into Ortannu the Chained's
// pod — a Pull cluster, at least two averaging three (ADR 0036). The Binding is
// Rarity.Connected, reachable only through its lead.
var ortannuCluster = card.Cluster{
	Name:     "Ortannu the Chained",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Ortannu the Chained
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	Reap: Put each Ortannu's Binding from your discard pile into your hand. For each card put into your hand this way, deal 2 damage to a creature and 2 damage to each of its neighbors.
var OrtannuTheChained = card.New(
	"Ortannu the Chained",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "97"),
	card.LeadsCluster(ortannuCluster),
	card.WithPower(7),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Reap, card.Sentences{
			Effects: []card.Effect{
				card.PutFromDiscard{
					Name:        OrtannusBinding.Name,
					All:         true,
					Destination: card.To.Hand,
				},
				card.ForEach{
					Times: card.ProducedThisWay{Tally: card.Tally.CardsReturned},
					Do: card.DealDamage{Spread: card.CreatureAndNeighbors{
						Amount: 2,
						Splash: 2,
					}},
				},
			},
		}),
)
