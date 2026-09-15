package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Daemo-Bot
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Scientist
//
//	Reap: Discard a card from your hand -> draw a card.
//	Destroyed: Steal 1 Æmber.
var DaemoBot = set.New(
	"Daemo-Bot",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "068"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
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
