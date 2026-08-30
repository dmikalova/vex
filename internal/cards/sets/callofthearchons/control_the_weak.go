package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Control the Weak
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a house - your opponent must choose that house as their active house during their next turn.
var ControlTheWeak = card.New(
	"Control the Weak",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 55),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.ForceOpponentActiveHouse{},
		}),
)
