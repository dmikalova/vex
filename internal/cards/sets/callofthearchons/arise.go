package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Arise!
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a house - put each creature of the chosen house from your discard pile into your hand, and gain 1 chain.
var Arise = card.New(
	"Arise!",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 54),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Sequence{
				Effects: []card.Effect{
					card.PutFromDiscard{
						Type:          card.Type.Creature,
						Destination:   card.To.Hand,
						All:           true,
						OfChosenHouse: true,
					},
					card.GainChains{Amount: 1},
				},
			},
		}),
)
