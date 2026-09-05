package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Poke
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 1 damage to an enemy creature. If this damage destroys that creature, draw a card.
var Poke = card.New(
	"Poke",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 117),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DamageThenIfDestroyed{
			Amount: 1,
			Target: card.Target.EnemyCreature,
			Then:   card.Draw{Amount: 1},
		}),
)
