package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Whistling Darts
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to each enemy creature.
var WhistlingDarts = set.New(
	"Whistling Darts",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "281"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			Target: card.Target.EachEnemyCreature,
		}),
)
