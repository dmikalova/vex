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
//	Reap: Choose a house - reveal the top card of your deck. If it is of the chosen house, archive it, and gain 1 Æmber. Otherwise, discard it.
var VespilonTheorist = set.New(
	"Vespilon Theorist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "155"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.ChooseHouseThen{
			Then: card.Sentences{
				Effects: []card.Effect{
					card.RevealTopOfDeck{Amount: 1},
					card.Conditional{
						Cond: card.ItIs{House: card.Houses.Chosen},
						Then: card.Sequence{
							Effects: []card.Effect{
								card.PutRevealedCard{To: card.Into.Archives},
								card.GainAember{
									Player: card.Controller,
									Amount: 1,
								},
							},
						},
						Else: card.PutRevealedCard{To: card.Into.Discard},
					},
				},
			},
		}),
)
