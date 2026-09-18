package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Sucker Punch
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Alpha.
//	Play: Deal 2 damage to an enemy creature. If this damage destroys that creature, archive Sucker Punch.
var SuckerPunch = set.New(
	"Sucker Punch",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "277"),
	card.WithBonus(card.Bonus.Aember),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			After:  card.IfDestroyed,
			Target: card.Target.EnemyCreature,
			Then:   card.ArchiveSource{},
		}),
)
