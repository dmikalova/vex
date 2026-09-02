package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Teliga
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Witch
//
//	After your opponent plays a card, if it is a creature, gain 1 Æmber.
var Teliga = card.New(
	"Teliga",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 366),
	card.WithPower(3),
	card.WithTraits("Human", "Witch"),
	card.WithAbility(card.Trigger.AfterEnemyCardPlayed, card.Conditional{
		Cond: card.ItIs{Type: card.Type.Creature},
		Then: card.GainAember{
			Player: card.Controller,
			Amount: 1,
		},
	}),
)
