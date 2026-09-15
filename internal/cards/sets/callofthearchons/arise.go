package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Arise!
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a house - put each creature of the chosen house from your discard pile into your hand. Gain 1 chain.
var Arise = set.New(
	"Arise!",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "54"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Sentences{
				Effects: []card.Effect{
					card.PutFromDiscard{
						Selection: card.Each{
							Type:  card.Type.Creature,
							House: card.Houses.Chosen,
						},
						Destination: card.To.Hand,
					},
					card.GainChains{Amount: 1},
				},
			},
		}),
)
