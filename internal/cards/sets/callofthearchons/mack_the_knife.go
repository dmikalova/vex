package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mack the Knife
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive. Versatile.
//	Action: Deal 1 damage to a creature. If this damage destroys that creature, gain 1 Æmber.
var MackTheKnife = card.New(
	"Mack the Knife",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 302),
	card.WithPower(3),
	card.WithTraits("Elf", "Thief"),
	card.WithKeywords(card.Keyword.Elusive, card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.DamageIfDestroyed{
			Amount: 1,
			Target: card.Target.Creature,
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
