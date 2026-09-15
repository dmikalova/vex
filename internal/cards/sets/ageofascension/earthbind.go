package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Earthbind
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature cannot be used unless you have discarded a card from your hand this turn.
var Earthbind = set.New(
	"Earthbind",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "352"),
	card.WithBonus(card.Bonus.Aember),
	card.WithRestrictions(card.Restrictions{
		UseCondition: card.CardsDiscarded{
			Player: card.Controller,
			Amount: 1,
		},
	}),
)
