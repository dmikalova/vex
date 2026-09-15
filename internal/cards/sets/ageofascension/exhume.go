package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Exhume
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Play a creature from your discard pile.
var Exhume = set.New(
	"Exhume",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "59"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.PlayFrom{
			From:  card.Discard,
			Types: card.Types.Of(card.Type.Creature),
		}),
)
