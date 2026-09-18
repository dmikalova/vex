package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Poke
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to an enemy creature. If this damage destroys that creature, draw a card.
var Poke = set.New(
	"Poke",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "117"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			After:  card.IfDestroyed,
			Target: card.Target.EnemyCreature,
			Then:   card.Draw{Amount: 1},
		}),
)
