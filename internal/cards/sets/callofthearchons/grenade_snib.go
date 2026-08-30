package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Grenade Snib
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Goblin
//
//	Destroyed: Your opponent loses 2 Æmber.
var GrenadeSnib = card.New(
	"Grenade Snib",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 34),
	card.WithPower(2),
	card.WithTraits("Goblin"),
	card.WithAbility(
		card.Trigger.Destroyed, card.LoseAember{
			Player: card.Opponent,
			Amount: 2,
		}),
)
