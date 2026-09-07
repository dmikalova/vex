package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Gold Key Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Players cannot forge their third key.
var GoldKeyImp = card.New(
	"Gold Key Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.AoA, 72),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithRestrictions(card.Restrictions{NoForgeKeyNumber: 3}),
)

// TODO: should not be special
