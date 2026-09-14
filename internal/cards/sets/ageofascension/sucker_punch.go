package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Sucker Punch
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Alpha.
//	Play: Deal 2 damage to an enemy Creature. If this damage destroys that Creature, archive Sucker Punch.
var SuckerPunch = set.New(
	"Sucker Punch",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "277"),
	card.WithAemberBonus(1),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.DamageThen{
			Amount: 2,
			After:  card.IfDestroyed,
			Target: card.Target.EnemyCreature,
			Then:   card.ArchiveSource{},
		}),
)
