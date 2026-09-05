package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Bonerot Venom
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "After this creature is used, deal 2 damage to this creature."
var BonerotVenom = card.New(
	"Bonerot Venom",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 283),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.UsedSelf,
			Effect:  card.DealDamage{Amount: 2, Target: card.Target.This},
		}},
	}),
)
