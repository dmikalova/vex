package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Stampede
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If you used 3 or more creatures this turn, steal 2 Æmber.
var Stampede = card.New(
	"Stampede",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 335),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.CountIs{
				Count:  card.CreaturesUsed{Player: card.Controller},
				Is:     card.AtLeast,
				Amount: 3,
			},
			Then: card.StealAember{Amount: 2},
		}),
)
