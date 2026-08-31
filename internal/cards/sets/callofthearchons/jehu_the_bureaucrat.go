package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Jehu the Bureaucrat
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	After you choose Sanctum as your active house, gain 2 Æmber.
var JehuTheBureaucrat = card.New(
	"Jehu the Bureaucrat",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 250),
	card.WithPower(3),
	card.WithTraits("Human"),
	card.WithAbility(
		card.Trigger.AfterChooseHouse, card.Conditional{
			Cond: card.ChoseHouse{House: card.House.Sanctum},
			Then: card.GainAember{Player: card.Controller, Amount: 2},
		}),
)
