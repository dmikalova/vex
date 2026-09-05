package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Aubade the Grim
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Spirit • Knight
//
//	Play: Aubade the Grim captures 3 Æmber from your opponent.
//	Reap: Move 1 Æmber from Aubade the Grim to the common supply.
var AubadeTheGrim = card.New(
	"Aubade the Grim",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 213),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Spirit, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 3,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
	card.WithAbility(
		card.Trigger.Reap, card.MoveAemberToCommonSupply{
			Amount: 1,
			Target: card.Target.This,
		}),
)
