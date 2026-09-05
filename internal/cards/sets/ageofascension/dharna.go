package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Dharna
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Witch
//
//	Play: For each friendly damaged creature in play, gain 1 Æmber.
//	Reap: Heal 2 damage from a friendly creature.
var Dharna = card.New(
	"Dharna",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 319),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Witch),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per: card.InPlay{
				Player:  card.Controller,
				Type:    card.Type.Creature,
				Damaged: true,
			},
		}),
	card.WithAbility(
		card.Trigger.Reap, card.Heal{
			Amount: 2,
			Target: card.Target.FriendlyCreature,
		}),
)
