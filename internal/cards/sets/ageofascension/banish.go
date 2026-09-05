package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Banish
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Archive an enemy creature from play.
var Banish = card.New(
	"Banish",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 54),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ArchiveFromPlay{Target: card.Target.EnemyCreature}),
)
