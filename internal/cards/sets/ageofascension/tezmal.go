package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Tezmal
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Reap: Choose a house. Your opponent cannot choose that house as their active house during their next turn.
var Tezmal = set.New(
	"Tezmal",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "66"),
	card.OneCopyPerDeck(),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.ChooseHouseThen{
			Then: card.CannotChooseHouse{
				Player:    card.Opponent,
				Reference: card.ChosenActiveHouse,
			},
		}),
)
