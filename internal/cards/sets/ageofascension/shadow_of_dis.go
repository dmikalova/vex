package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Shadow of Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Until your next turn, enemy Creatures' text boxes are considered blank (except for traits).
var ShadowOfDis = card.New(
	"Shadow of Dis",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "103"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.BlankEnemyText{}),
)
