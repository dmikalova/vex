package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// They're Everywhere!
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to each enemy flank creature. Deal 1 damage to each enemy creature that is not on a flank.
var TheyreEverywhere = set.New(
	"They're Everywhere!",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "334"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
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
