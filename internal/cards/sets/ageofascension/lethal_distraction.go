package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Lethal Distraction
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: For the remainder of the turn, whenever a creature takes damage, it takes an additional 2 damage.
var LethalDistraction = set.New(
	"Lethal Distraction",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "305"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.TakesExtraDamage{
			Target: card.Target.Creature,
			Amount: 2,
		}),
)
