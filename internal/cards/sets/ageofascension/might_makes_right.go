package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Might Makes Right
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: You may destroy any number of friendly Creatures with total power of 25 or more - forge a key at no cost -> purge Might Makes Right.
var MightMakesRight = card.New(
	"Might Makes Right",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "43"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DestroyFriendlyCreaturesToForge{
			Target:        card.Target.EachFriendlyCreature,
			MinTotalPower: 25,
			Then:          card.ForgeKey{FreeOfCost: true},
		}),
)
