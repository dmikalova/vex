package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mother
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Robot • Scientist
//
//	Your hand size is 1 more.
var Mother = set.New(
	"Mother",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, "145"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Robot, card.Traits.Scientist),
	card.WithDrawModifier(card.Controller, 1),
)
