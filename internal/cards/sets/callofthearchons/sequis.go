package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sequis
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Reap: Sequis captures 1 Æmber from your opponent.
var Sequis = card.New(
	"Sequis",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 257),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Reap, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
