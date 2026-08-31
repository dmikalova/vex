package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Deep Probe
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Choose a house - reveal your opponent's hand, and discard each creature of the chosen house from your opponent's hand.
var DeepProbe = card.New(
	"Deep Probe",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 162),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Sequence{Effects: []card.Effect{
				card.RevealHand{Player: card.Opponent},
				card.DiscardHand{
					Player:        card.Opponent,
					Types:         card.Types(card.Type.Creature),
					OfChosenHouse: true,
				},
			}},
		}),
)
