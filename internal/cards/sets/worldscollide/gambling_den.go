package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Gambling Den
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	At the start of each player's turn, you may choose a house - reveal the top card of your deck. If it is of the chosen house, gain 2 Æmber. Otherwise, lose 2 Æmber.
var GamblingDen = card.New(
	"Gambling Den",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "268"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.AfterAnyPlayerStartOfTurn, card.May{
			Do: card.ChooseHouseThen{
				Then: card.Sentences{
					Effects: []card.Effect{
						card.RevealTopOfDeck{},
						card.Conditional{
							Cond: card.ItIsOfHouse{House: card.TheChosenHouse},
							Then: card.GainAember{Player: card.Controller, Amount: 2},
							Else: card.LoseAember{Player: card.Controller, Amount: 2},
						},
					},
				},
			},
		}),
)
