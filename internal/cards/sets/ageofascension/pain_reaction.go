package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Pain Reaction
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 2 damage to an enemy creature. If this damage destroys that creature, deal 2 damage to each of that creature's neighbors.
var PainReaction = card.New(
	"Pain Reaction",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "78"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DamageThenIfDestroyed{
			Amount: 2,
			Target: card.Target.EnemyCreature,
			Then: card.DealDamage{
				Amount: 2,
				Target: card.Target.FormerNeighbors,
			},
		}),
)
