package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Perplexing Sophistry
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If you have more Æmber than your opponent, your opponent discards a random card from their hand, and you draw a card.
var PerplexingSophistry = card.New(
	"Perplexing Sophistry",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 293),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.YourAember{Is: card.MoreThanOpponent},
			Then: card.Sequence{
				Effects: []card.Effect{
					card.DiscardRandomFromHand{Player: card.Opponent},
					card.Draw{
						Amount: 1,
						You:    true,
					},
				},
			},
		}),
)
