package massmutation

import "github.com/dmikalova/vex/internal/card"

// monumentToFaustCluster pulls a Faust the Great into Monument to Faust's pod, so
// the Monument's discard-pile bonus has its namesake to feed it (ADR 0036). Faust
// the Great also drafts on its own, so the pull only guarantees the pairing.
var monumentToFaustCluster = card.Cluster{
	Name:     "Monument to Faust",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Monument to Faust
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: If Faust the Great is in your discard pile, keys cost +2 Æmber during your opponent's next turn. Otherwise, keys cost +1 Æmber during your opponent's next turn.
var MonumentToFaust = set.New(
	"Monument to Faust",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "236"),
	card.LeadsCluster(monumentToFaustCluster),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.NamedCardInDiscard{Name: FaustTheGreat.Name},
			Then: card.RaiseKeyCost{
				Player:   card.Opponent,
				Amount:   2,
				Duration: card.Duration.OpponentNextTurn,
			},
			Else: card.RaiseKeyCost{
				Player:   card.Opponent,
				Amount:   1,
				Duration: card.Duration.OpponentNextTurn,
			},
		}),
)
