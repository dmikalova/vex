package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Key Abduction
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put each Mars creature into its owner's hand. Forge a key at +9 Æmber current cost, reduced by 1 Æmber for each card in your hand.
var KeyAbduction = card.New(
	"Key Abduction",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 166),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Sentence{Effect: card.PutFromPlay{
				Target:      card.Target.EachCreature.OfHouse(card.House.Mars),
				Destination: card.To.Hand,
			}},
			card.Sentence{Effect: card.ForgeKey{
				Extra: 9,
				ReducedBy: card.CardsInHand{
					Player: card.Controller,
					House:  card.AnyHouse,
				},
			}},
		}}),
)
