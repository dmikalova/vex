package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Prince Derric, Unifier
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  1
//	Traits: Human • Knight
//
//	Play: If 3 or more houses are represented among friendly creatures, gain 3 Æmber.
var PrinceDerricUnifier = set.New(
	"Prince Derric, Unifier",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "240"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.CountIs{
				Count: card.HousesAmong{
					Player: card.Controller,
					Type:   card.Type.Creature,
				},
				Is:     card.AtLeast,
				Amount: 3,
			},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 3,
			},
		}),
)
