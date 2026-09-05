package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Grumpus Tamer
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	Reap: Search your deck and discard pile for a War Grumpus, reveal it, and put it into your hand.
var GrumpusTamer = card.New(
	"Grumpus Tamer",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 39),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Reap, card.SearchForName{Name: "War Grumpus"}),
)
