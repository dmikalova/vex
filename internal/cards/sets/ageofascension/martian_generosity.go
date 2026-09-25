package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Martian Generosity
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Lose all your Æmber. For each Æmber you lost this way, draw 2 cards.
var MartianGenerosity = set.New(
	"Martian Generosity",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "202"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.LoseAember{
					Player: card.Controller,
					By:     card.AllAember,
				},
				card.Draw{
					Amount: 2,
					Per: card.ProducedThisWay{
						Tally:  card.Tally.AemberLost,
						Player: card.Controller,
					},
				},
			},
		}),
)
