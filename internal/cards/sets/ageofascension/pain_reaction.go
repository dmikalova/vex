package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Pain Reaction
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to an enemy creature. If this damage destroys that creature, deal 2 damage to each of that creature's neighbors.
var PainReaction = set.New(
	"Pain Reaction",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "78"),
	card.WithBonus(card.Bonus.Aember),
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
