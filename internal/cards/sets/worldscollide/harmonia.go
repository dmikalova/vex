package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Harmonia
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Witch
//
//	Elusive.
//	After you play a Creature, if you are overwhelmed, gain 1 Æmber.
var Harmonia = card.New(
	"Harmonia",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "357"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(card.Trigger.AfterCardPlayed, card.Conditional{
		Cond: card.ItIs{Type: card.Type.Creature},
		Then: card.Conditional{
			Cond: card.Overwhelmed{},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		},
	}),
)
