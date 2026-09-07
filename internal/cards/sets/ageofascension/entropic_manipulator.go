package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Entropic Manipulator
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Redistribute the damage among a player's creatures.
var EntropicManipulator = card.New(
	"Entropic Manipulator",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 195),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.RedistributeDamage{}),
)
