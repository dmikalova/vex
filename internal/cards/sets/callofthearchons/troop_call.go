package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Troop Call
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put each Niffle trait creature from your discard pile into your hand. Put each friendly Niffle trait creature into its owner's hand.
var TroopCall = card.New(
	"Troop Call",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 337),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Sentence{Effect: card.PutFromDiscard{
					Type:        card.Type.Creature,
					Trait:       "Niffle",
					All:         true,
					Destination: card.To.Hand,
				}},
				card.Sentence{Effect: card.PutFromPlay{
					Target:      card.Target.EachFriendlyCreature.WithTrait("Niffle"),
					Destination: card.To.Hand,
				}},
			},
		}),
)
