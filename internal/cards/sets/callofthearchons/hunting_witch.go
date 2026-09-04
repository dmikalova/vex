package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hunting Witch
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Witch
//
//	After you play a creature, gain 1 Æmber.
var HuntingWitch = card.New(
	"Hunting Witch",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 367),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithAbility(card.Trigger.AfterCardPlayed, card.Conditional{
		Cond: card.ItIs{Type: card.Type.Creature},
		Then: card.GainAember{
			Player: card.Controller,
			Amount: 1,
		},
	}),
)
