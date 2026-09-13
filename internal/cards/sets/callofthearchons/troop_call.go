package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// troopCallCluster pulls Niffle creatures into Troop Call's pod — a Pull cluster
// with a per-partner rate: a couple of Niffle Apes (averaging three) and, much
// less often, a Niffle Queen (min zero, averaging under one) (ADR 0036).
var troopCallCluster = card.Cluster{
	Name:     "Troop Call",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}

// Troop Call
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put each Niffle Creature from your discard pile into your hand. Put each friendly Niffle Creature into its owner's hand.
var TroopCall = card.New(
	"Troop Call",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "337"),
	card.LeadsCluster(troopCallCluster),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{
			Effects: []card.Effect{
				card.PutFromDiscard{
					Match:       card.Match{Type: card.Type.Creature, Trait: card.Traits.Niffle},
					All:         true,
					Destination: card.To.Hand,
				},
				card.PutFromPlay{
					Target:      card.Target.EachFriendlyCreature.WithTrait(card.Traits.Niffle),
					Destination: card.To.Hand,
				},
			},
		}),
)
