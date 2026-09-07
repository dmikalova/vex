package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Maruck the Marked
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Spirit • Knight
//
//	After Maruck the Marked prevents damage with its armor, for each damage just prevented, Maruck the Marked captures 1 Æmber from your opponent.
var MaruckTheMarked = card.New(
	"Maruck the Marked",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "220"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Spirit, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.AfterArmorPrevents, card.CaptureAember{
			Amount: 1,
			Per:    card.DamagePrevented{},
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
