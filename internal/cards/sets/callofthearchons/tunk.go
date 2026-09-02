package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tunk
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Robot
//
//	After you play a Mars creature, fully heal Tunk.
var Tunk = card.New(
	"Tunk",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 199),
	card.WithPower(6),
	card.WithTraits("Robot"),
	card.WithAbility(card.Trigger.AfterCardPlayed, card.Conditional{
		Cond: card.ItIs{
			House: card.House.Self,
			Type:  card.Type.Creature,
		},
		Then: card.Heal{
			Fully:  true,
			Target: card.Target.This,
		},
	}),
)
