package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Troll
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  8
//	Traits: Giant
//
//	Reap: Heal 3 damage from Troll.
var Troll = card.New(
	"Troll",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 48),
	card.WithPower(8),
	card.WithTraits("Giant"),
	card.WithAbility(
		card.Trigger.Reap, card.Heal{
			Amount: 3,
			Target: card.Target.This,
		}),
)
