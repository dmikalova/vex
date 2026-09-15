package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lifeweb
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If your opponent played 3 or more creatures on their previous turn, steal 2 Æmber.
var Lifeweb = set.New(
	"Lifeweb",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "326"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.CountIs{
				Count: card.TurnCount{
					Player: card.Opponent,
					Of:     card.TurnStat.CreaturesPlayedLastTurn,
				},
				Is:     card.AtLeast,
				Amount: 3,
			},
			Then: card.StealAember{Amount: 2},
		}),
)
