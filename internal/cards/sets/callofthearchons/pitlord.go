package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pitlord
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  9
//	Bonus:  Æmber Æmber
//	Traits: Demon
//
//	Taunt.
//	While Pitlord is in play you must choose Dis as your active house.
var Pitlord = set.New(
	"Pitlord",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "93"),
	card.WithBonus(card.Bonus.Aember, card.Bonus.Aember),
	card.WithPower(9),
	card.WithTraits(card.Traits.Demon),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithHouseLock(card.HouseLock{
		Player: card.Controller,
		House:  card.House.Self,
	}),
)
