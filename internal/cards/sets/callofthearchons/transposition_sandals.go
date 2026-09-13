package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Transposition Sandals
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Action: Swap this Creature with another friendly Creature in your battleline. Use the other Creature."
var TranspositionSandals = card.New(
	"Transposition Sandals",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "159"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.Action, Effect: card.Sentences{Effects: []card.Effect{
				card.Swap{With: card.Target.OtherFriendlyCreature},
				card.OnChooseCreature{
					Target: card.Target.TheOtherCreature,
					Verbs:  []card.CreatureVerb{card.UseVerb{}},
				},
			}}},
		},
	}),
)
