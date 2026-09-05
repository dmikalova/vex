package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Binding Irons
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent gains 3 chains.
var BindingIrons = card.New(
	"Binding Irons",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 55),
	card.WithAbility(
		card.Trigger.Play, card.GainChains{
			Player: card.Opponent,
			Amount: 3,
		}),
)
