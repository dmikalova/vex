package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Pound
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to a creature that is not on a flank and 1 damage to each of its neighbors.
var Pound = card.New(
	"Pound",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 15),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{Spread: card.CreatureAndNeighbors{
			Amount:     2,
			Splash:     1,
			NotOnFlank: true,
		}}),
)
