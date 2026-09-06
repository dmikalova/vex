package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Streke
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	While Streke is not on a flank, during their "draw cards" phase, your opponent refills their hand to 1 less card.
var Streke = card.New(
	"Streke",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 65),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithDrawModifierOffFlank(card.Opponent, -1),
)
