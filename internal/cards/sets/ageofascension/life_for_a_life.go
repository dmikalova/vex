package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Life for a Life
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Destroy a friendly creature -> deal 6 damage to a creature.
var LifeForALife = card.New(
	"Life for a Life",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 273),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First:  card.Destroy{Target: card.Target.FriendlyCreature},
			Result: card.DealDamage{Amount: 6, Target: card.Target.Creature},
		}),
)
