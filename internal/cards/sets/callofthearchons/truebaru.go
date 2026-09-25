package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Truebaru
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	Taunt.
//	In order to play Truebaru, you must lose 3 Æmber.
//	Destroyed: Gain 5 Æmber.
var Truebaru = set.New("Truebaru",
	card.House.Dis, card.Type.Creature, card.Rarity.Rare,
	card.Provenance(card.CotA, "104"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Demon),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAemberCost(3),
	card.WithAbility(
		card.Trigger.Destroyed, card.GainAember{
			Player: card.Controller,
			Amount: 5,
		}),
)
