package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Banish
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Archive an enemy creature from play.
var Banish = set.New(
	"Banish",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "54"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ArchiveFromPlay{Target: card.Target.EnemyCreature}),
)
