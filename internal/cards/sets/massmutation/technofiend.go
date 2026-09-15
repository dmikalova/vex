package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Techno-Fiend
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Demon
//
//	Reap: Discard a card from your hand -> draw a card.
//	Destroyed: Steal 1 Æmber.
var TechnoFiend = set.New(
	"Techno-Fiend",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "016"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Reap, card.Then{
			First: card.DiscardCard{
				Player:    card.Controller,
				Zone:      card.Hand,
				Selection: card.Chosen{},
				Amount:    1,
			},
			Result: card.Draw{Amount: 1},
		}),
	card.WithAbility(
		card.Trigger.Destroyed, card.StealAember{Amount: 1}),
)
