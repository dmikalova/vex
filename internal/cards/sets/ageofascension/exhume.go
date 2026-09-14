package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Exhume
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Play a Creature from your discard pile.
var Exhume = set.New(
	"Exhume",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "59"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PlayFrom{
			From:  card.Discard,
			Types: card.Types.Of(card.Type.Creature),
		}),
)
