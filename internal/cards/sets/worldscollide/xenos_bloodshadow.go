package worldscollide

import "github.com/dmikalova/vex/internal/card"

// xenosBloodshadowCluster pulls a Toad into Xenos Bloodshadow's pod, one for one —
// a Pull cluster Xenos Bloodshadow leads.
var xenosBloodshadowCluster = card.Cluster{
	Name:     "Xenos Bloodshadow",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Xenos Bloodshadow
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Human • Witch
//
//	Elusive, Poison, Skirmish, Hazardous 6.
var XenosBloodshadow = set.New(
	"Xenos Bloodshadow",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "404"),
	card.LeadsCluster(xenosBloodshadowCluster),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithKeywords(card.Keyword.Elusive, card.Keyword.Poison, card.Keyword.Skirmish),
	card.WithHazardous(6),
)
