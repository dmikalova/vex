package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Healing Blast
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Fully heal a creature. If you healed 4 or more damage, gain 2 Æmber.
var HealingBlast = card.New(
	"Healing Blast",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 219),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Heal{Fully: true, Target: card.Target.Creature},
			card.Conditional{
				Cond: card.CountIs{
					Count:  card.DamageHealed{},
					Is:     card.AtLeast,
					Amount: 4,
				},
				Then: card.GainAember{Player: card.Controller, Amount: 2},
			},
		}}),
)
