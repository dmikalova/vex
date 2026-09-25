package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Thero Centurion
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Play/Fight: Thero Centurion captures 1 Æmber from your opponent.
var TheroCenturion = set.New(
	"Thero Centurion",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "195"),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
	card.WithAbility(
		card.Trigger.Fight, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
