package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Shooler
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Play: If your opponent has 4 Æmber or more, steal 1 Æmber.
var Shooler = card.New(
	"Shooler",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 96),
	card.WithPower(5),
	card.WithTraits("Demon"),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAember{Is: card.AtLeast, Amount: 4},
			Then: card.StealAember{Amount: 1},
		}),
)
