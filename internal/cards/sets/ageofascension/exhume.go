package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Exhume
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Play a creature from your discard pile.
var Exhume = card.New(
	"Exhume",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 59),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PlayFrom{
			From: card.Discard,
			Type: card.Type.Creature,
		}),
)
