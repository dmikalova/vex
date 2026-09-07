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
//	Play: Deal 2 damage to an enemy creature. If this damage destroys that creature, archive Sucker Punch.
var SuckerPunch = card.New(
	"Sucker Punch",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "277"),
	card.WithAemberBonus(1),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.DamageThenIfDestroyed{
			Amount: 2,
			Target: card.Target.EnemyCreature,
			Then:   card.ArchiveSource{},
		}),
)
