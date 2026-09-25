package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Dusk Chronicles
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: If your opponent has more Æmber than you, draw a card. If you have more Æmber than your opponent, archive a card from your hand.
var DuskChronicles = set.New(
	"Dusk Chronicles",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "268"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Conditional{
					Cond: card.PoolAember{
						Player: card.Opponent,
						Is:     card.MoreThanYou,
					},
					Then: card.Draw{Amount: 1},
				},
				card.Conditional{
					Cond: card.PoolAember{
						Player: card.Controller,
						Is:     card.MoreThanOpponent,
					},
					Then: card.ArchiveCard{
						Zone:      card.Hand,
						Selection: card.Chosen{},
					},
				},
			},
		}),
)
