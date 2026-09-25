package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Cull the Weak
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Destroy the least powerful enemy creature.
var CullTheWeak = set.New(
	"Cull the Weak",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "57"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachEnemyCreature.Refine(card.LeastPowerful),
		}),
)
