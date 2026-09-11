package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// bearFluteCluster pulls a couple of Ancient Bears into Bear Flute's pod — a Pull
// cluster, at least two averaging about two and a half (ADR 0036).
var bearFluteCluster = card.Cluster{
	Name:     "Bear Flute",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Bear Flute
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Fully heal Ancient Bear. If there are no Ancient Bears in play, search your deck and discard pile and put each Ancient Bear from them into your hand -> shuffle your discard pile into your deck.
var BearFlute = card.New(
	"Bear Flute",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "340"),
	card.LeadsCluster(bearFluteCluster),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Heal{
				Fully:  true,
				Target: card.Target.Creature.Named(AncientBear.Name),
			},
			card.Conditional{
				Cond: card.InPlay{
					Player: card.EachPlayer,
					Type:   card.Type.Creature,
					Name:   AncientBear.Name,
					None:   true,
				},
				Then: card.Then{
					First:  card.SearchForName{Name: AncientBear.Name, All: true},
					Result: card.ShuffleIntoDeck{Zones: []card.Zone{card.Discard}},
				},
			},
		}}),
)
