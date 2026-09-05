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
//	Play: If you control creatures from 3 or more houses, gain 3 Æmber.
var PrinceDerricUnifier = card.New(
	"Prince Derric, Unifier",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 240),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.ControlsCreaturesOfHouses{Amount: 3},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 3,
			},
		}),
)
