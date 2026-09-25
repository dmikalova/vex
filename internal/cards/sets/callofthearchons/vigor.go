package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Vigor
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Heal 3 damage from a creature. If you healed 3 or more damage, gain 1 Æmber.
var Vigor = set.New(
	"Vigor",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "338"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Heal{
				Amount: 3,
				Target: card.Target.Creature,
			},
			card.Conditional{
				Cond: card.CountIs{
					Count:  card.DamageHealed{},
					Is:     card.AtLeast,
					Amount: 3,
				},
				Then: card.GainAember{
					Player: card.Controller,
					Amount: 1,
				},
			},
		}}),
)
