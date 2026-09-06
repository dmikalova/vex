package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Baron Mengevin
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Human • Knight
//
//	After you discard a card from your hand, if it is a Sanctum card, Baron Mengevin captures 1 Æmber from your opponent.
var BaronMengevin = card.New(
	"Baron Mengevin",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 227),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithAbility(card.Trigger.AfterDiscardFromHand, card.Conditional{
		Cond: card.ItIs{House: card.House.Self},
		Then: card.CaptureAember{
			Target: card.Target.This,
			Amount: 1,
			Source: card.Opponent,
		},
	}),
)
