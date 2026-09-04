package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Valdr
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	Valdr deals +2 Damage while attacking an enemy creature on the flank.
var Valdr = card.New(
	"Valdr",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 29),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithAttackDamage(card.AttackDamage{
		Amount:    2,
		FlankOnly: true,
	}),
)
