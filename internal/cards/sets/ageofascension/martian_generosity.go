package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Martian Generosity
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Lose all your Æmber, and for each Æmber you lost this way, draw 2 cards.
var MartianGenerosity = card.New(
	"Martian Generosity",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 202),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.LoseAember{
					Player: card.Controller,
					By:     card.AllAember,
				},
				card.Draw{
					Amount: 2,
					Per:    card.AemberLostThisWay{Player: card.Controller},
				},
			},
		}),
)
