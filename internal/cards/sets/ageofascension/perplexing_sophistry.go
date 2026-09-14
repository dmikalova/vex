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
var PerplexingSophistry = set.New(
	"Perplexing Sophistry",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "293"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.PoolAember{Player: card.Controller, Is: card.MoreThanOpponent},
			Then: card.Sequence{
				Effects: []card.Effect{
					card.DiscardCard{
						Player:    card.Opponent,
						Zone:      card.Hand,
						Selection: card.Random{},
					},
					card.Draw{
						Amount: 1,
						You:    true,
					},
				},
			},
		}),
)
