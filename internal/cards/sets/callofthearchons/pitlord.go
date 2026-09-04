package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pitlord
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  9
//	Æmber:  2
//	Traits: Demon
//
//	Taunt.
//	While Pitlord is in play you must choose Dis as your active house.
var Pitlord = card.New(
	"Pitlord",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 93),
	card.WithPower(9),
	card.WithAemberBonus(2),
	card.WithTraits(card.Traits.Demon),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithHouseLock(card.HouseLock{
		Player: card.Controller,
		House:  card.House.Self,
	}),
)
