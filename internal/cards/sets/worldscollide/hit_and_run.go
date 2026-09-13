package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hit and Run
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 2 damage to a Creature. Put a friendly Creature into its owner's hand.
var HitAndRun = card.New(
	"Hit and Run",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "238"),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.DealDamage{
				Amount: 2,
				Target: card.Target.Creature,
			},
			card.PutFromPlay{
				Target:      card.Target.FriendlyCreature,
				Destination: card.To.Hand,
			},
		}}),
)
