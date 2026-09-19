package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// grumpusTamerCluster pulls a couple of War Grumpuses into Grumpus Tamer's pod — a
// Pull cluster, at least two averaging three (ADR 0036). War Grumpus is Rare and
// rolls on its own too, so this tops the pod up to at least that many.
var grumpusTamerCluster = card.Cluster{
	Name:     "Grumpus Tamer",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Grumpus Tamer
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	Reap: Search your deck and discard pile for a War Grumpus, reveal it, and put it into your hand. Shuffle your deck.
var GrumpusTamer = set.New(
	"Grumpus Tamer",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "39"),
	card.LeadsCluster(grumpusTamerCluster),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Reap, card.Sequence{
			Effects: []card.Effect{
				card.Search{
					Sources: []card.Zone{card.Deck, card.Discard},
					Filter:  card.Filter{Name: WarGrumpus.Name},
					Reveal:  true,
					Dest:    card.To.Hand,
				},
				card.Shuffle{},
			},
		}),
)
