package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Pain Reaction
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 2 damage to an enemy Creature. If this damage destroys that Creature, deal 2 damage to each of that Creature's neighbors.
var PainReaction = card.New(
	"Pain Reaction",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "78"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DamageThen{
			Amount: 2,
			After:  card.IfDestroyed,
			Target: card.Target.EnemyCreature,
			Then: card.DealDamage{
				Amount: 2,
				Target: card.Target.FormerNeighbors,
			},
		}),
)
