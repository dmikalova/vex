package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Ortannu's Binding
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Connected
//	Æmber:  1
//
//	Play: Deal 2 damage to a friendly creature.
var OrtannusBinding = card.New(
	"Ortannu's Binding",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Connected,
	card.Provenance(card.AoA, 98),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.FriendlyCreature,
		}),
)
