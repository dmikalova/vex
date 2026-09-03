package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Vigor
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Heal 3 damage from a creature. If you healed exactly 3 damage, gain 1 Æmber.
var Vigor = card.New(
	"Vigor",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 338),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Heal{
				Amount: 3,
				Target: card.Target.Creature,
			},
			card.Conditional{
				Cond: card.CountIs{
					Count:  card.DamageHealed{},
					Is:     card.Exactly,
					Amount: 3,
				},
				Then: card.GainAember{
					Player: card.Controller,
					Amount: 1,
				},
			},
		}}),
)
