package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Life for a Life
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Destroy a friendly creature -> deal 6 damage to a creature.
var LifeForALife = set.New(
	"Life for a Life",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "273"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.Destroy{Target: card.Target.FriendlyCreature},
			Result: card.DealDamage{
				Amount: 6,
				Target: card.Target.Creature,
			},
		}),
)
