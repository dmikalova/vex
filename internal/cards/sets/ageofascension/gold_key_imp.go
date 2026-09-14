package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Gold Key Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Players cannot forge their third key.
var GoldKeyImp = set.New(
	"Gold Key Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "72"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithRestrictions(card.Restrictions{NoForgeKeyNumber: 3}),
)
