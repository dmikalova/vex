package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Bronze Key Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Players cannot forge their first key.
var BronzeKeyImp = card.New(
	"Bronze Key Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.AoA, 71),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithRestrictions(card.Restrictions{NoForgeKeyNumber: 1}),
)
