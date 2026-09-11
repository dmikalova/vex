package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// hydeCluster pulls one Velum per Hyde: the two are a PullExact pair, so Velum
// (Rarity.Connected) rides in whenever Hyde is placed (ADR 0036).
var hydeCluster = card.Cluster{
	Name:     "Hyde",
	Strategy: card.ClusterStrategy.PullExact,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Hyde
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Human • Scientist
//
//	Reap: Draw a card. If you control Velum, draw a card.
//	Destroyed: Archive Velum from your discard pile -> archive Hyde from play.
var Hyde = card.New(
	"Hyde",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "167"),
	card.LeadsCluster(hydeCluster),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(card.Trigger.Reap, card.Sentences{Effects: []card.Effect{
		card.Draw{Amount: 1},
		card.Conditional{
			Cond: card.ControlsNamed{Name: Velum.Name},
			Then: card.Draw{Amount: 1},
		},
	}}),
	card.WithAbility(card.Trigger.Destroyed, card.Then{
		First:  card.ArchiveCard{Zone: card.Discard, Selection: card.Named{Name: Velum.Name}},
		Result: card.ArchiveFromPlay{Target: card.Target.This},
	}),
)
