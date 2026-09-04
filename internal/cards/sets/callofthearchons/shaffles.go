package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Shaffles
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	At the end of your turn, your opponent loses 1 Æmber.
var Shaffles = card.New(
	"Shaffles",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 95),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithAbility(card.Trigger.EndOfTurn, card.LoseAember{
		Player: card.Opponent,
		Amount: 1,
	}),
)
