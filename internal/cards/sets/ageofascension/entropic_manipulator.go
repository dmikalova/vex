package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Entropic Manipulator
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Redistribute the damage among a player's creatures.
var EntropicManipulator = set.New(
	"Entropic Manipulator",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "195"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.RedistributeDamage{}),
)
