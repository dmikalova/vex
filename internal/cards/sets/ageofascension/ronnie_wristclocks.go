package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Ronnie Wristclocks
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: Steal 1 Æmber, or 2 if your opponent has 7 Æmber or more.
var RonnieWristclocks = card.New(
	"Ronnie Wristclocks",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 276),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{
			Amount: 1,
			Or: card.OrAmount{
				Amount: 2,
				When:   card.OpponentAember{Is: card.AtLeast, Amount: 7},
			},
		}),
)
