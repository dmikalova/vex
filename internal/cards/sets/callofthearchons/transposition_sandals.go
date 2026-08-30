package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Transposition Sandals
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Action: Swap this creature with another friendly creature in your battleline. Use the other creature."
var TranspositionSandals = card.New(
	"Transposition Sandals",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 159),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.Action, Effect: card.Sequence{Effects: []card.Effect{
				card.Sentence{Effect: card.Swap{With: card.Target.OtherFriendlyCreature}},
				card.Sentence{Effect: card.OnChooseCreature{
					Target: card.Target.TheOtherCreature,
					Verbs:  []card.CreatureVerb{card.UseVerb{}},
				}},
			}}},
		},
	}),
)
