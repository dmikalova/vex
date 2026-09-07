package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Earthbind
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature cannot be used unless you have discarded a card from your hand this turn.
var Earthbind = card.New(
	"Earthbind",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "352"),
	card.WithAemberBonus(1),
	card.WithRestrictions(card.Restrictions{
		UseCondition: card.CardsDiscarded{
			Player: card.Controller,
			Amount: 1,
		},
	}),
)
