package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Krump
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	After a creature is destroyed fighting Krump, your opponent loses 1 Æmber.
var Krump = card.New(
	"Krump",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 39),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.LoseAember{
			Player: card.Opponent,
			Amount: 1,
		}),
)
