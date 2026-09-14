package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Lethal Distraction
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: For the remainder of the turn, whenever a Creature takes damage, it takes an additional 2 damage.
var LethalDistraction = set.New(
	"Lethal Distraction",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "305"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.TakesExtraDamage{
			Target: card.Target.Creature,
			Amount: 2,
		}),
)
