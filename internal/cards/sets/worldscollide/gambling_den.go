package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Gambling Den
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	At the start of each player's turn, you may choose a house. If you do, reveal the top card of your deck. If it is of the chosen house, gain 2 Æmber. Otherwise, lose 2 Æmber.
var GamblingDen = set.New(
	"Gambling Den",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "268"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithEachPlayerAbility(
		card.Trigger.StartOfTurn, card.May{
			Do: card.ChooseHouseThen{
				Then: card.Sequence{
					Effects: []card.Effect{
						card.RevealTopOfDeck{Amount: 1},
						card.Conditional{
							Cond: card.ItIs{House: card.Houses.Chosen},
							Then: card.GainAember{
								Player: card.Controller,
								Amount: 2,
							},
							Else: card.LoseAember{
								Player: card.Controller,
								Amount: 2,
							},
						},
					},
				},
			},
		}),
)
