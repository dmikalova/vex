package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// 1-2 Punch
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose an enemy creature - if that creature was already stunned, destroy it. Otherwise, stun it.
var Card12Punch = card.New(
	"1-2 Punch",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 1),
	card.WithAemberBonus(1),
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
