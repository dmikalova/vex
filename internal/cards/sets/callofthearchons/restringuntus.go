package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Restringuntus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Demon
//
//	Play: Choose a house - your opponent cannot choose that house as their active house until Restringuntus leaves play.
var Restringuntus = card.New(
	"Restringuntus",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 94),
	card.WithPower(1),
	card.WithTraits(card.Traits.Demon),
	card.WithHouseLock(card.HouseLock{
		Player: card.Opponent,
		Bars:   true,
	}),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.NameHouse{Player: card.Opponent},
		}),
)
