package massmutation

import "github.com/dmikalova/vex/internal/card"

// monumentToPrimusCluster pulls a Consul Primus into Monument to Primus's pod, so
// the Monument's discard-pile bonus has its namesake to feed it (ADR 0036). Consul
// Primus also drafts on its own, so the pull only guarantees the pairing.
var monumentToPrimusCluster = card.Cluster{
	Name:     "Monument to Primus",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Monument to Primus
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: If Consul Primus is in your discard pile, move 1 Æmber from a creature to another creature. Otherwise, move 1 Æmber from a friendly creature to another friendly creature.
var MonumentToPrimus = set.New(
	"Monument to Primus",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "238"),
	card.LeadsCluster(monumentToPrimusCluster),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.NamedCardInDiscard{Name: ConsulPrimus.Name},
			Then: card.MoveAember{
				Amount: 1,
				From:   card.Target.Creature,
				Onto:   card.Target.OtherCreature,
			},
			Else: card.MoveAember{
				Amount: 1,
				From:   card.Target.FriendlyCreature,
				Onto:   card.Target.OtherFriendlyCreature,
			},
		}),
)
