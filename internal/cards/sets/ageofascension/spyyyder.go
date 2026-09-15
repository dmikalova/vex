package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Spyyyder
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Demon
//
//	Skirmish.
//	Spyyyder gains poison while attacking an enemy flank creature.
var Spyyyder = set.New(
	"Spyyyder",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "84"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Demon),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAttackKeywords(card.AttackKeywords{
		Keywords:  []card.KeywordValue{card.Keyword.Poison},
		FlankOnly: true,
	}),
)
