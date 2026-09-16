package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// monumentToShrixCluster pulls a Citizen Shrix into Monument to Shrix's pod, so
// the Monument's discard-pile bonus has its namesake to feed it (ADR 0036). Citizen
// Shrix also drafts on its own, so the pull only guarantees the pairing.
var monumentToShrixCluster = card.Cluster{
	Name:     "Monument to Shrix",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Monument to Shrix
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	You may spend Æmber on Monument to Shrix when forging keys.
//	Action: If Citizen Shrix is in your discard pile, move 1 Æmber from any player's pool to Monument to Shrix. Otherwise, move 1 Æmber from your pool to Monument to Shrix.
var MonumentToShrix = set.New(
	"Monument to Shrix",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "239"),
	card.LeadsCluster(monumentToShrixCluster),
	card.WithTraits(card.Traits.Location),
	card.WithSpendableAember(),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.NamedCardInDiscard{Name: CitizenShrix.Name},
			Then: card.MoveAemberFromPool{
				Amount: 1,
				Target: card.Target.This,
				Source: card.ChosenPlayer,
			},
			Else: card.MoveAemberFromPool{
				Amount: 1,
				Target: card.Target.This,
			},
		}),
)
