package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// fayginCluster pulls a couple of Urchins into Faygin's pod — a Pull cluster, at
// least two averaging about two and a half (ADR 0036).
var fayginCluster = card.Cluster{
	Name:     "Faygin",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Faygin
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Thief
//
//	Elusive.
//	Reap: Put an Urchin from play or from your discard pile into your hand.
var Faygin = set.New(
	"Faygin",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "300"),
	card.LeadsCluster(fayginCluster),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.ReturnNamedToHand{Name: Urchin.Name}),
)
