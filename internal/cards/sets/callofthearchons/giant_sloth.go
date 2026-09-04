package callofthearchons

import (
	"github.com/dmikalova/vactrol/internal/card"
)

// Giant Sloth
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Beast
//
//	You cannot use this card unless you have discarded an Untamed card from your hand this turn.
//	Action: Gain 3 Æmber.
var GiantSloth = card.New(
	"Giant Sloth",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 354),
	card.WithPower(6),
	card.WithTraits(card.Traits.Beast),
	card.WithRestrictions(card.Restrictions{
		UseCondition: card.CardsDiscarded{
			Player: card.Controller,
			House:  card.House.Self,
			Amount: 1,
		},
	}),
	card.WithAbility(
		card.Trigger.Action, card.GainAember{
			Player: card.Controller,
			Amount: 3,
		}),
)
