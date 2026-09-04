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
//	During your "draw cards" step, refill your hand to 1 additional card.
var Mother = card.New(
	"Mother",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 145),
	card.WithPower(5),
	card.WithTraits(card.Traits.Robot, card.Traits.Scientist),
	card.WithDrawModifier(card.Controller, 1),
)
