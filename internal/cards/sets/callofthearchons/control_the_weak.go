package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Control the Weak
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose a house - your opponent must choose that house as their active house during their next turn.
var ControlTheWeak = set.New(
	"Control the Weak",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "55"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.MustChooseHouse{
				Player:    card.Opponent,
				Reference: card.ChosenActiveHouse,
			},
		}),
)
