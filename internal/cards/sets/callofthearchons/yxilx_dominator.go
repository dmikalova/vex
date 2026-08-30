package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Yxilx Dominator
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  9
//	Traits: Robot
//
//	Taunt.
//	Yxilx Dominator enters play stunned.
var YxilxDominator = card.New(
	"Yxilx Dominator",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 205),
	card.WithPower(9),
	card.WithTraits("Robot"),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithEntersPlay(card.Stun{Target: card.Target.This}),
)
