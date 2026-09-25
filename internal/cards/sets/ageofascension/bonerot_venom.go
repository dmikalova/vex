package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Bonerot Venom
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "After this creature is used, deal 2 damage to it."
var BonerotVenom = set.New(
	"Bonerot Venom",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "283"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.UsedSelf,
			Effect: card.DealDamage{
				Amount: 2,
				Target: card.Target.This,
			},
		}},
	}),
)
