package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// We Can ALL Win
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Each player's keys cost -2 Æmber until the end of your next turn.
var WeCanALLWin = card.New(
	"We Can ALL Win",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "344"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.LowerKeyCost{
			Player:   card.EachPlayer,
			Amount:   2,
			Duration: card.Duration.EndOfPlayerNextTurn,
		}),
)
