package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Vespilon Theorist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	Elusive.
//	Reap: Choose a house - reveal the top card of your deck. If it is of the chosen house, archive the top card of your deck, and gain 1 Æmber. Otherwise, discard the top card of your deck.
var VespilonTheorist = card.New(
	"Vespilon Theorist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 155),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.ChooseHouseThen{
			Then: card.Sentences{
				Effects: []card.Effect{
					card.RevealTopOfDeck{},
					card.Conditional{
						Cond: card.ItIsOfHouse{House: card.TheChosenHouse},
						Then: card.Sequence{
							Effects: []card.Effect{
								card.ArchiveTopOfDeck{Amount: 1},
								card.GainAember{
									Player: card.Controller,
									Amount: 1,
								},
							},
						},
						Else: card.DiscardTopOfDeck{Player: card.Controller},
					},
				},
			},
		}),
)
