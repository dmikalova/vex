package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Snaglet
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Action: Choose a house. If your opponent chooses that house as their active house during their next turn, steal 2 Æmber.
var Snaglet = set.New(
	"Snaglet",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "118"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.ChooseHouseThen{
			Then: card.WagerOpponentChoosesChosenHouse{Amount: 2},
		}),
)
