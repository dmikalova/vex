package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Hexpion
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Robot
//
//	Destroyed: Archive Hexpion from play. Archive the top card of your deck.
var Hexpion = set.New(
	"Hexpion",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "113"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.Destroyed, card.Sequence{Effects: []card.Effect{
			card.ArchiveFromPlay{Target: card.Target.This},
			card.ArchiveCard{Zone: card.Deck, Selection: card.Top{}},
		}}),
)
