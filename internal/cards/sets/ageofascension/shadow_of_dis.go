package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Shadow of Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Until the start of your next turn, enemy creatures' text boxes are considered blank, except for traits.
var ShadowOfDis = set.New(
	"Shadow of Dis",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "103"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.BlankEnemyText{}),
)
