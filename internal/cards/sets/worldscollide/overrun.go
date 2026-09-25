package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Overrun
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If 3 or more enemy creatures have been destroyed this turn, your opponent loses 2 Æmber.
var Overrun = set.New(
	"Overrun",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "25"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.CountIs{
				Count: card.TurnCount{
					Player: card.Controller,
					Of:     card.TurnStat.EnemyCreaturesDestroyed,
				},
				Is:     card.AtLeast,
				Amount: 3,
			},
			Then: card.LoseAember{
				Player: card.Opponent,
				Amount: 2,
			},
		}),
)
