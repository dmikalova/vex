package ageofascension

import "github.com/dmikalova/vex/internal/card"

// 1-2 Punch
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose an enemy creature. If that creature was already stunned, destroy it. Otherwise, stun it.
var Card12Punch = set.New(
	"1-2 Punch",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "1"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.EnemyCreature,
			Then: card.Conditional{
				Cond: card.ItIsStunned{},
				Then: card.Destroy{Target: card.Target.Triggering},
				Else: card.Stun{Target: card.Target.Triggering},
			},
		}),
)
