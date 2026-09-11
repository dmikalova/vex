package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// timetravellerCluster pulls one Help from Future Self per Timetraveller: the two
// are a PullExact pair, so N Timetravellers ride in with N copies of Help (ADR
// 0036). Help is Rarity.Connected, reachable only through its lead.
var timetravellerCluster = card.Cluster{
	Name:     "Timetraveller",
	Strategy: card.ClusterStrategy.PullExact,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Timetraveller
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Scientist
//
//	Play: Draw 2 cards.
//	Action: Shuffle Timetraveller into its owner's deck.
var Timetraveller = card.New(
	"Timetraveller",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "153"),
	card.LeadsCluster(timetravellerCluster),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Play, card.Draw{Amount: 2}),
	card.WithAbility(
		card.Trigger.Action, card.PutFromPlay{
			Target:      card.Target.This,
			Destination: card.To.DeckShuffled,
		}),
)
