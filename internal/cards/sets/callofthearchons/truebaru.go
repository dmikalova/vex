package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Truebaru
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	Taunt.
//	You must lose 3 Æmber in order to play Truebaru.
//	Destroyed: Gain 5 Æmber.
var Truebaru = card.New("Truebaru",
	card.House.Dis, card.Type.Creature, card.Rarity.Rare,
	card.Provenance(card.CotA, 104),
	card.WithPower(7),
	card.WithTraits("Demon"),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAemberCost(3),
	card.WithAbility(
		card.Trigger.Destroyed, card.GainAember{Player: card.Controller, Amount: 5}),
)
