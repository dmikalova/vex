package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// They're Everywhere!
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to each enemy flank creature. Deal 1 damage to each enemy creature that is not on a flank.
var TheyreEverywhere = card.New(
	"They're Everywhere!",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 334),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.DealDamage{
				Amount: 2,
				Target: card.Target.EachEnemyCreature.OnFlank(),
			},
			card.DealDamage{
				Amount: 1,
				Target: card.Target.EachEnemyCreature.NotOnFlank(),
			},
		}}),
)
