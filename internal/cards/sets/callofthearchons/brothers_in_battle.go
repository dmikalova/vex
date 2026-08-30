package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Brothers in Battle
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Choose a house - for the remainder of the turn, each friendly creature of the chosen house may fight.
var BrothersInBattle = card.New(
	"Brothers in Battle",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 4),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.GrantFightForChosenHouse{},
		}),
)
