package massmutation

import "github.com/dmikalova/vex/internal/card"

// Cephaloist
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	While you have 4 Æmber or more, your Æmber cannot be stolen.
var Cephaloist = set.New(
	"Cephaloist",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "362"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithAemberCannotBeStolen(
		card.PoolAember{
			Player: card.Controller,
			Is:     card.AtLeast,
			Amount: 4,
		},
	),
)
