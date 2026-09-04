package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ember Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Your opponent cannot play more than 2 cards each turn.
var EmberImp = card.New(
	"Ember Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 85),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithRestrictions(card.Restrictions{
		PlayCardLimit: card.PlayCardLimit{
			Player: card.Opponent,
			Amount: 2,
		},
	}),
)
