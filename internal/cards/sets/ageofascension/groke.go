package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Groke
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Fight: Your opponent loses 1 Æmber.
var Groke = card.New(
	"Groke",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 12),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Fight, card.LoseAember{
			Player: card.Opponent,
			Amount: 1,
		}),
)
