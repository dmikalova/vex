package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Cull the Weak
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Destroy the least powerful enemy creature.
var CullTheWeak = card.New(
	"Cull the Weak",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 57),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachEnemyCreature.Selector(card.LeastPowerful),
		}),
)
