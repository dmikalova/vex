package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Little Niff
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Elf • Thief
//
//	Deploy, Elusive, Omega.
//	After a neighbor of Little Niff is used to fight, steal 1 Æmber.
var LittleNiff = card.New(
	"Little Niff",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 289),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(
		card.Keyword.Deploy,
		card.Keyword.Elusive,
		card.Keyword.Omega,
	),
	card.WithAbility(
		card.Trigger.AfterNeighborFights, card.StealAember{Amount: 1}),
)
