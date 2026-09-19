package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// monumentToOctaviaCluster pulls a Cornicen Octavia into Monument to Octavia's
// pod, so the Monument's discard-pile bonus has its namesake to feed it (ADR
// 0036). Cornicen Octavia also drafts on its own, so the pull only guarantees the
// pairing.
var monumentToOctaviaCluster = card.Cluster{
	Name:     "Monument to Octavia",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Monument to Octavia
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: A friendly creature captures 1 Æmber from your opponent. If Cornicen Octavia is in your discard pile, it captures 1 Æmber from your opponent.
var MonumentToOctavia = set.New(
	"Monument to Octavia",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "237"),
	card.LeadsCluster(monumentToOctaviaCluster),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.CaptureAember{
				Amount: 1,
				Target: card.Target.FriendlyCreature,
				Source: card.Opponent,
			},
			card.Conditional{
				Cond: card.NamedCardInDiscard{Name: CornicenOctavia.Name},
				Then: card.CaptureAember{
					Amount: 1,
					Target: card.Target.Triggering,
					Source: card.Opponent,
				},
			},
		}}),
)
