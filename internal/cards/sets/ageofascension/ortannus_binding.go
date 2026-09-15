package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Ortannu's Binding
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Connected
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to a friendly creature.
var OrtannusBinding = set.New(
	"Ortannu's Binding",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Connected,
	card.Provenance(card.AoA, "98"),
	card.InCluster(card.Pulled(ortannuCluster, 2, 3)),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.FriendlyCreature,
		}),
)
