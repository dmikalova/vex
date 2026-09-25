package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Brothers in Battle
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Choose a house. For the remainder of the turn, each friendly creature of the chosen house may fight.
var BrothersInBattle = set.New(
	"Brothers in Battle",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "4"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.MayPlayOrUse{
				Houses: card.GrantHouses.Chosen,
				Grant:  card.GrantFight,
			},
		}),
)
